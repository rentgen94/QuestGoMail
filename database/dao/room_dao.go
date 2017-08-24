package dao

import (
	"database/sql"
	"github.com/rentgen94/QuestGoMail/entities"
)

const (
	getRoomByIdQuery = `
		SELECT id, name, description FROM Room WHERE id = $1
	`
	getRoomRelatedSlotIdsQuery = `
		SELECT id FROM Slot WHERE room = $1
	`
	getRoomRelatedInteractiveIdsQuery = `
		SELECT id FROM Interactive WHERE room = $1
	`
)

type RoomDAO struct {
	db             *sql.DB
	slotDao        *SlotDAO
	interactiveDao *InteractiveDAO
}

func NewRoomDAO(db *sql.DB) *RoomDAO {
	return &RoomDAO{
		db:             db,
		slotDao:        NewSlotDAO(db),
		interactiveDao: NewInteractiveDAO(db),
	}
}

func (dao *RoomDAO) GetById(id int) (*entities.Room, error) {
	var room, roomErr = dao.getEmptyRoom(id)
	if roomErr != nil {
		return nil, roomErr
	}

	var slots, slotsErr = dao.getRelatedSlots(id)
	if slotsErr != nil {
		return nil, slotsErr
	}

	var interactives, interactivesErr = dao.getRelatedInteractives(id)
	if interactivesErr != nil {
		return nil, interactivesErr
	}

	room.SetSlots(slots)
	room.SetInteractives(interactives)

	return room, nil
}

func (dao *RoomDAO) getEmptyRoom(roomId int) (*entities.Room, error) {
	var acceptor = struct {
		id          int
		name        string
		description string
	}{}

	var err = dao.db.
		QueryRow(getRoomByIdQuery, roomId).
		Scan(&acceptor.id, &acceptor.name, &acceptor.description)

	if err != nil {
		return nil, err
	}

	return entities.NewRoom(acceptor.id, acceptor.name, acceptor.description), nil
}

func (dao *RoomDAO) getRelatedInteractives(roomId int) (entities.InteractivesType, error) {
	var ids, idsErr = getRelatedEntityIds(dao.db, getRoomRelatedInteractiveIdsQuery, roomId)
	if idsErr != nil {
		return nil, idsErr
	}

	var interactives = make(entities.InteractivesType)
	for _, id := range ids {
		var interactive, interactiveErr = dao.interactiveDao.GetById(id)
		if interactiveErr != nil {
			return nil, interactiveErr
		}

		interactives[interactive.Id()] = interactive
	}

	return interactives, nil
}

func (dao *RoomDAO) getRelatedSlots(roomId int) (entities.SlotsType, error) {
	var ids, idsErr = getRelatedEntityIds(dao.db, getRoomRelatedSlotIdsQuery, roomId)
	if idsErr != nil {
		return nil, idsErr
	}

	var slots = make(entities.SlotsType)
	for _, id := range ids {
		var slot, slotErr = dao.slotDao.GetById(id)
		if slotErr != nil {
			return nil, slotErr
		}

		slots[slot.Id()] = slot
	}

	return slots, nil
}

func getRelatedEntityIds(db *sql.DB, query string, roomId int) ([]int, error) {
	var rows, err = db.Query(query, roomId)
	if err != nil {
	return nil, err
	}
	defer rows.Close()

	var ids = make([]int, 0)
	for rows.Next() {
	var id int
	err = rows.Scan(&id)
	if err != nil {
	return nil, err
	}

	ids = append(ids, id)
	}

	err = rows.Err()
	if err != nil {
	return nil, err
	}

	return ids, nil
}
