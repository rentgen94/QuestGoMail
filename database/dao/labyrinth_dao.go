package dao

import (
	"database/sql"
	"github.com/rentgen94/QuestGoMail/entities"
)

const (
	getAllLabyrinths = `
		SELECT id, name, description, startRoom FROM Labyrinth ORDER BY id
	`

	getLabyrinthByIdQuery = `
		SELECT id, name, description, startRoom FROM Labyrinth WHERE id = $1
	`
	getLabyrinthRelatedRoomIdsQuery = `
		SELECT room FROM LabyrinthRoomLink WHERE labyrinth = $1 ORDER BY room
	`
	getLabyrinthRelatedDoorIdsQuery = `
		SELECT d.id FROM Labyrinth lab
			JOIN LabyrinthRoomLink link ON lab.id = link.labyrinth
			JOIN Room r1 ON r1.id = link.room
			JOIN Room r2 ON r2.id = link.room
			JOIN Door d ON d.room1 = r1.id AND d.room2 = r2.id
		WHERE lab.id = $1
		ORDER BY d.id
	`
	getLabyrinthRelatedActionIdsQuery = `
		SELECT action FROM LabyrinthActionLink WHERE labyrinth = $1 ORDER BY action
	`
	getActionBindsQuery = `
		SELECT r.id room_id, i.id interactive_id, ai_link.action action_id
		FROM LabyrinthRoomLink lr_link
			JOIN Room r ON r.id = lr_link.room
			JOIN Interactive i ON i.room = r.id
			JOIN ActionInteractiveLink ai_link ON i.id = ai_link.interactive
		WHERE lr_link.labyrinth = $1

	`
)

type LabyrinthDAO struct {
	db        *sql.DB
	roomDao   *RoomDAO
	doorDao   *DoorDAO
	actionDao *ActionDAO
}

type LabyrinthDescription struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	startRoomId int
}

func NewLabyrinthDAO(db *sql.DB) *LabyrinthDAO {
	return &LabyrinthDAO{
		db:        db,
		roomDao:   NewRoomDAO(db),
		doorDao:   NewDoorDAO(db),
		actionDao: NewActionDAO(db),
	}
}

func (dao *LabyrinthDAO) GetById(id int) (*entities.Labyrinth, error) {
	var labInfo, infoErr = dao.getEmptyLabyrinthInfo(id)
	if infoErr != nil {
		return nil, infoErr
	}

	var actions, actErr = dao.getRelatedActions(id)
	if actErr != nil {
		return nil, actErr
	}

	var rooms, roomErr = dao.getRelatedRooms(id)
	if roomErr != nil {
		return nil, roomErr
	}

	var doors, doorErr = dao.getRelatedDoors(id, rooms)
	if doorErr != nil {
		return nil, doorErr
	}

	var lab = entities.NewLabyrinth(labInfo.Id, labInfo.Name, labInfo.Description)
	lab.SetRooms(rooms)
	lab.SetStartRoom(rooms[labInfo.startRoomId])
	lab.SetDoors(doors)
	lab.SetActions(actions)

	dao.bindActions(lab)

	return lab, nil
}

func (dao *LabyrinthDAO) GetAll() ([]LabyrinthDescription, error) {
	var rows, err = dao.db.Query(getAllLabyrinths)
	if err != nil {
		return nil, err
	}

	var result = make([]LabyrinthDescription, 0)
	for rows.Next() {
		var description = LabyrinthDescription{}
		err = rows.Scan(&description.Id, &description.Name, &description.Description, &description.startRoomId)
		if err != nil {
			return nil, err
		}

		result = append(result, description)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, nil
}

// this function has to be called at the end of labyrinth loading
func (dao *LabyrinthDAO) bindActions(labyrinth *entities.Labyrinth) error {
	for _, action := range labyrinth.Actions() {
		action.SetLabyrinth(labyrinth)
	}

	type idInfoType struct {
		roomId        int
		interactiveId int
		actionId      int
	}
	var idInfoSlice = make([]idInfoType, 0)

	var rows, err = dao.db.Query(getActionBindsQuery, labyrinth.Id())
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var idInfo idInfoType
		err = rows.Scan(&idInfo.roomId, &idInfo.interactiveId, &idInfo.actionId)
		if err != nil {
			return err
		}

		idInfoSlice = append(idInfoSlice, idInfo)
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	for _, item := range idInfoSlice {
		labyrinth.
			Rooms()[item.roomId].
			Interactives()[item.interactiveId].
			SetAction(
				labyrinth.Actions()[item.actionId],
			)
	}

	return nil
}

func (dao *LabyrinthDAO) getRelatedActions(labId int) (entities.ActionsType, error) {
	var ids, idsErr = getRelatedEntityIds(dao.db, getLabyrinthRelatedActionIdsQuery, labId)
	if idsErr != nil {
		return nil, idsErr
	}

	var actions = make(entities.ActionsType)
	for _, id := range ids {
		var action, err = dao.actionDao.GetById(id)
		if err != nil {
			return nil, err
		}

		actions[action.Id()] = action
	}

	return actions, nil
}

func (dao *LabyrinthDAO) getRelatedDoors(labId int, rooms entities.RoomsType) (entities.DoorsType, error) {
	var ids, idsErr = getRelatedEntityIds(dao.db, getLabyrinthRelatedDoorIdsQuery, labId)
	if idsErr != nil {
		return nil, idsErr
	}

	var doors = make(entities.DoorsType)
	for _, id := range ids {
		var doorInfo, err = dao.doorDao.getDoorInfo(id)
		if err != nil {
			return nil, err
		}

		var door = doorInfo.door
		door.SetRoom1(rooms[doorInfo.room1Id]) // assumed, room always found
		door.SetRoom2(rooms[doorInfo.room2Id])

		rooms[doorInfo.room1Id].Doors()[door.Id()] = door
		rooms[doorInfo.room2Id].Doors()[door.Id()] = door

		doors[door.Id()] = door
	}

	return doors, nil
}

func (dao *LabyrinthDAO) getRelatedRooms(labId int) (entities.RoomsType, error) {
	var ids, idsErr = getRelatedEntityIds(dao.db, getLabyrinthRelatedRoomIdsQuery, labId)
	if idsErr != nil {
		return nil, idsErr
	}

	var rooms = make(entities.RoomsType)
	for _, id := range ids {
		var room, roomErr = dao.roomDao.GetById(id)
		if roomErr != nil {
			return nil, roomErr
		}

		rooms[room.Id()] = room
	}

	return rooms, nil
}

func (dao *LabyrinthDAO) getEmptyLabyrinthInfo(id int) (*LabyrinthDescription, error) {
	var acceptor = new(LabyrinthDescription)

	var err = dao.db.
		QueryRow(getLabyrinthByIdQuery, id).
		Scan(&acceptor.Id, &acceptor.Name, &acceptor.Description, &acceptor.startRoomId)

	if err != nil {
		return nil, err
	}

	return acceptor, nil
}
