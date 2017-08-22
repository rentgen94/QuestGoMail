package dao

import (
	"database/sql"
	"github.com/rentgen94/QuestGoMail/entities"
)

type ItemDAO struct {
	db *sql.DB
}

func NewItemDAO(db *sql.DB) *ItemDAO {
	var result = new(ItemDAO)
	result.db = db
	return result
}

func (dao *ItemDAO) GetByName(name string) (entities.Item, error) {
	return dao.getByParameter(
		"SELECT id, name, description, size FROM Item WHERE name = $1",
		name,
	)
}

func (dao *ItemDAO) GetById(id int) (entities.Item, error) {
	return dao.getByParameter(
		"SELECT id, name, description, size FROM Item WHERE id = $1",
		id,
	)
}

func (dao *ItemDAO) getByParameter(query string, parameter interface{}) (entities.Item, error) {
	var result entities.Item
	var err = dao.db.
		QueryRow(query, parameter).
		Scan(&result.Id, &result.Name, &result.Description, &result.Size)

	return result, err
}
