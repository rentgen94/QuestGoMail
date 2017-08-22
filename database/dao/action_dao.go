package dao

import (
	"database/sql"
	"github.com/rentgen94/QuestGoMail/entities"
)

type ActionDAO struct {
	db *sql.DB
}

type actionModel struct {
	roomIds []int
}

func (dao *ActionDAO) GetById(id int) *entities.Action {

}

func (dao *ActionDAO) getAction(id int) *entities.Action {
	var acceptor = struct {
		id int
		name string
		isAccessible bool
		capacity int
	}{}

	var err = dao.db.
		QueryRow(`
			SELECT id, name, capacity, isAccessible FROM Slot WHERE id = $1
		`, id).
		Scan(&acceptor.id, &acceptor.name, &acceptor.capacity, &acceptor.isAccessible)

	if err != nil {
		return nil, err
	}

	var result = entities.NewSlot(acceptor.id, acceptor.name, acceptor.capacity, acceptor.isAccessible)
	return result, err
}


//CREATE TABLE Action (
//id         SERIAL PRIMARY KEY,
//name       VARCHAR(100),
//resultCode INT,
//resultMsg  VARCHAR(500)
//);
