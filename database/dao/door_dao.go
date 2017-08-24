package dao

import (
	"database/sql"
	"github.com/rentgen94/QuestGoMail/entities"
)

const (
	getDoorByIdQuery = `
		SELECT id, name, isAccessible FROM Door WHERE id = $1
	`
	getDoorInfoQuery = `
		SELECT id, name, isAccessible, room1, room2 FROM Door WHERE id = $1
	`
)

type DoorDAO struct {
	db *sql.DB
}

func NewDoorDAO(db *sql.DB) *DoorDAO {
	return &DoorDAO{
		db: db,
	}
}

func (dao *DoorDAO) GetById(id int) (*entities.Door, error) {
	return dao.getOrphanDoor(id)
}

type doorInfo struct {
	door    *entities.Door
	room1Id int
	room2Id int
}

func (dao *DoorDAO) getDoorInfo(doorId int) (*doorInfo, error) {
	var acceptor = struct {
		id           int
		name         string
		isAccessible bool
		room1Id      int
		room2Id      int
	}{}

	var err = dao.db.
		QueryRow(getDoorInfoQuery, doorId).
		Scan(&acceptor.id, &acceptor.name, &acceptor.isAccessible, &acceptor.room1Id, &acceptor.room2Id)

	if err != nil {
		return nil, err
	}

	var result = new(doorInfo)
	result.door = entities.NewDoor(acceptor.id, acceptor.name, acceptor.isAccessible)
	result.room1Id = acceptor.room1Id
	result.room2Id = acceptor.room2Id

	return result, nil
}

func (dao *DoorDAO) getOrphanDoor(doorId int) (*entities.Door, error) {
	var acceptor = struct {
		id           int
		name         string
		isAccessible bool
	}{}

	var err = dao.db.
		QueryRow(getDoorByIdQuery, doorId).
		Scan(&acceptor.id, &acceptor.name, &acceptor.isAccessible)

	if err != nil {
		return nil, err
	}

	return entities.NewDoor(acceptor.id, acceptor.name, acceptor.isAccessible), nil
}
