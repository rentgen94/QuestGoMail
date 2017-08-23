package dao

import (
	"database/sql"
	"errors"
	"github.com/rentgen94/QuestGoMail/entities"
)

const (
	getInteractiveQuery = `
		SELECT id, room, name, description, isAccessible, args FROM Interactive WHERE id = $1
	`
	getNecessaryItemsQuery = `
		SELECT id FROM Item i JOIN InteractiveObjectNeed need ON i.id = need.item WHERE need.interactive = $1
	`

	argumentLengthDoesNotMatch = "Argument length does not match"
	gotWrongArguments          = "Got wrong arguments"
	itemListLengthDoesNotMatch = "Item list length does not match"
	notFoundNecessaryItem      = "Not found necessary item"
)

type InteractiveDAO struct {
	db      *sql.DB
	itemDao *ItemDAO
}

func NewInteractiveDAO(db *sql.DB) *InteractiveDAO {
	return &InteractiveDAO{db: db, itemDao: NewItemDAO(db)}
}

func (dao *InteractiveDAO) GetById(id int) (entities.InteractiveObject, error) {
	return dao.getInteractiveById(id)
}

func (dao *InteractiveDAO) getInteractiveById(id int) (entities.InteractiveObject, error) {
	var acceptor = struct {
		id           int
		room         int
		name         string
		description  string
		isAccessible bool
		args         []string
	}{}

	var err = dao.db.
		QueryRow(getInteractiveQuery, id).
		Scan(&acceptor.id, &acceptor.room, &acceptor.name, &acceptor.description, &acceptor.isAccessible, &acceptor.args)

	if err != nil {
		return nil, err
	}

	var necessaryItems, itemsErr = dao.getNecessaryItems(id)
	if itemsErr != nil {
		return nil, itemsErr
	}

	var inputChecker = dao.getInputChecker(acceptor.args, necessaryItems)

	var inter = entities.NewInteractiveObject(
		acceptor.id,
		acceptor.name,
		acceptor.description,
		acceptor.isAccessible,
		inputChecker,
		nil,
	)

	return inter, nil
}

func (*InteractiveDAO) getInputChecker(args []string, items []entities.Item) entities.InputChecker {
	var argValidator = func(expectedArgs []string, gotArgs []string) error {
		if len(expectedArgs) == 0 {
			return nil
		}

		if len(expectedArgs) != len(gotArgs) {
			return errors.New(argumentLengthDoesNotMatch)
		}

		for i := range args {
			if expectedArgs[i] != gotArgs[i] {
				return errors.New(gotWrongArguments)
			}
		}

		return nil
	}

	var itemValidator = func(expectedItems []entities.Item, gotItems []entities.Item) error {
		if len(expectedItems) == 0 {
			return nil
		}

		if len(expectedItems) != len(gotItems) {
			return errors.New(itemListLengthDoesNotMatch)
		}

		var idMap = make(map[int]bool)
		for _, item := range expectedItems {
			idMap[item.Id] = true
		}

		for _, item := range gotItems {
			var _, ok = idMap[item.Id]
			if !ok {
				return errors.New(notFoundNecessaryItem)
			}
		}

		return nil
	}

	return func(innerArgs []string, innerItems []entities.Item) error {
		var argErr = argValidator(args, innerArgs)
		if argErr != nil {
			return argErr
		}

		var itemErr = itemValidator(items, innerItems)
		if itemErr != nil {
			return itemErr
		}

		return nil
	}
}

func (dao *InteractiveDAO) getNecessaryItems(interactiveId int) ([]entities.Item, error) {
	var rows, err = dao.db.Query(getNecessaryItemsQuery, interactiveId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var itemId int
	var items = make([]entities.Item, 0)
	for rows.Next() {
		var rowErr = rows.Scan(&itemId)
		if rowErr != nil {
			return nil, rowErr
		}

		var item, itemErr = dao.itemDao.GetById(itemId)
		if itemErr != nil {
			return nil, itemErr
		}

		items = append(items, item)
	}

	var lastErr = rows.Err()

	return items, lastErr
}
