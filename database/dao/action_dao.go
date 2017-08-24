package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/rentgen94/QuestGoMail/entities"
)

const (
	roomDoesNotExistsTemplate        = "Room id = %d does not exists in labyrinth"
	slotDoesNotExistsTemplate        = "Slot id = %d does not exists in room %d"
	interactiveDoesNotExistsTemplate = "Interactive id = %d does not exists in room %d"
	doorDoesNotExistsTemplate        = "Door id = %d does not exists in room %d"

	actionRelatedDoorsQuery = `
		SELECT d.id, d.room1, sw.newState FROM Door d
			JOIN ActionDoorSwitch sw ON d.id = sw.door
		WHERE sw.action = $1
		ORDER BY d.id
	`
	actionRelatedInteractivesQuery = `
		SELECT i.id, i.room, sw.newState FROM Interactive i
			JOIN ActionInteractiveSwitch sw ON i.id = sw.interactive
		WHERE sw.action = $1
		ORDER BY i.id
	`
	actionRelatedSlotsQuery = `
		SELECT s.id, s.room, sw.newState FROM Slot s
			JOIN ActionSlotSwitch sw ON s.id = sw.slot
		WHERE sw.action = $1
		ORDER BY s.id
	`
	getActionByIdQuery = `
		SELECT id, name, resultCode, resultMsg FROM Action WHERE id = $1
	`
)

type ActionDAO struct {
	db *sql.DB
}

type entityData struct {
	id     int
	roomId int
	state  bool
}

func NewActionDAO(db *sql.DB) *ActionDAO {
	var result = new(ActionDAO)
	result.db = db
	return result
}

func (dao *ActionDAO) GetById(id int) (*entities.Action, error) {
	return dao.getAction(id)
}

func (dao *ActionDAO) getAction(id int) (*entities.Action, error) {
	var acceptor = struct {
		id         int
		name       string
		resultCode int
		resultMsg  string
	}{}

	var err = dao.db.
		QueryRow(getActionByIdQuery, id).
		Scan(&acceptor.id, &acceptor.name, &acceptor.resultCode, &acceptor.resultMsg)

	if err != nil {
		return nil, err
	}

	var actFunc, funcErr = dao.getActFunc(acceptor.id, acceptor.resultMsg, acceptor.resultCode)
	if funcErr != nil {
		return nil, funcErr
	}

	return entities.NewAction(id, acceptor.name, actFunc), nil
}

func (dao *ActionDAO) getActFunc(actionId int, resultMsg string, resultCode int) (entities.ActType, error) {
	var slotsData, slotsErr = dao.getRelatedSlotsInfo(actionId)
	if slotsErr != nil {
		return nil, slotsErr
	}

	var interactivesData, interactivesErr = dao.getRelatedInteractivesInfo(actionId)
	if interactivesErr != nil {
		return nil, interactivesErr
	}

	var doorsData, doorsErr = dao.getRelatedDoorsInfo(actionId)
	if doorsErr != nil {
		return nil, doorsErr
	}

	return func(labyrinth *entities.Labyrinth) (entities.InteractionResult, error) {
		var processAccessible = func(
			accessibleData entityData,
			extractor func(room *entities.Room) (entities.Accessible, error),
		) error {
			var room, okRoom = labyrinth.Rooms()[accessibleData.roomId]
			if !okRoom {
				return errors.New(fmt.Sprintf(roomDoesNotExistsTemplate, accessibleData.roomId))
			}

			var accessible, err = extractor(room)
			if err != nil {
				return err
			}

			accessible.SetAccessible(accessibleData.state)
			return nil
		}

		var err error

		for _, slotData := range slotsData {
			err = processAccessible(
				slotData,
				func(room *entities.Room) (entities.Accessible, error) {
					var slot, ok = room.Slots()[slotData.id]
					if !ok {
						return nil, errors.New(fmt.Sprintf(slotDoesNotExistsTemplate, slotData.id, slotData.roomId))
					}
					return slot, nil
				},
			)
		}
		if err != nil {
			return entities.InteractionResult{}, err
		}

		for _, interData := range interactivesData {
			err = processAccessible(
				interData,
				func(room *entities.Room) (entities.Accessible, error) {
					var inter, ok = room.Interactives()[interData.id]
					if !ok {
						return nil, errors.New(fmt.Sprintf(interactiveDoesNotExistsTemplate, interData.id, interData.roomId))
					}
					return inter, nil
				},
			)
		}
		if err != nil {
			return entities.InteractionResult{}, err
		}

		// TODO maybe doors must be keeped inside the labyrinth
		for _, doorData := range doorsData {
			err = processAccessible(
				doorData,
				func(room *entities.Room) (entities.Accessible, error) {
					var inter, ok = room.Doors()[doorData.id]
					if !ok {
						return nil, errors.New(fmt.Sprintf(doorDoesNotExistsTemplate, doorData.id, doorData.roomId))
					}
					return inter, nil
				},
			)
		}
		if err != nil {
			return entities.InteractionResult{}, err
		}

		return entities.InteractionResult{
			Msg:  resultMsg,
			Code: resultCode,
		}, nil
	}, nil
}

func (dao *ActionDAO) getRelatedSlotsInfo(actionId int) ([]entityData, error) {
	return dao.getRelatedEntityInfo(actionId, actionRelatedSlotsQuery)
}

func (dao *ActionDAO) getRelatedInteractivesInfo(actionId int) ([]entityData, error) {
	return dao.getRelatedEntityInfo(actionId, actionRelatedInteractivesQuery)
}

func (dao *ActionDAO) getRelatedDoorsInfo(actionId int) ([]entityData, error) {
	return dao.getRelatedEntityInfo(actionId, actionRelatedDoorsQuery)
}

func (dao *ActionDAO) getRelatedEntityInfo(actionId int, query string) ([]entityData, error) {
	var rows, err = dao.db.Query(query, actionId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result = make([]entityData, 0)
	for rows.Next() {
		var item entityData

		err = rows.Scan(&item.id, &item.roomId, &item.state)
		if err != nil {
			break
		}

		result = append(result, item)
	}

	if err != nil {
		return nil, err
	} else {
		err = rows.Err()
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}
