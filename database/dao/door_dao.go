package dao

import (
	"database/sql"
	"github.com/rentgen94/QuestGoMail/entities"
)

const (
	getDoorByIdQuery = `
		SELECT id, name, isAccessible FROM Door WHERE id = $1
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

func (dao *DoorDAO) getOrphanDoor(doorId int) (*entities.Door, error) {
	var acceptor = struct {
		id          int
		name        string
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
