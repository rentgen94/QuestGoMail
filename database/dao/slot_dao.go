package dao

import (
	"database/sql"
	"github.com/rentgen94/QuestGoMail/entities"
)

type SlotDAO struct {
	db *sql.DB
}

func NewSlotDAO(db *sql.DB) *SlotDAO {
	var result = new(SlotDAO)
	result.db = db
	return result
}

func (dao *SlotDAO) GetById(id int) (*entities.Slot, error) {
	var slot, slotErr = dao.getSlot(id)
	if slotErr != nil {
		return nil, slotErr
	}

	var items, itemErr = dao.getSlotItems(id)
	if itemErr != nil {
		return nil, itemErr
	}

	var err error
	for _, item := range items {
		err = slot.PutItem(item)
		if err != nil {
			break
		}
	}

	return slot, err
}

func (dao *SlotDAO) getSlotItems(id int) ([]entities.Item, error) {
	var rows, err = dao.db.Query(`
			SELECT i.id, i.name, i.description, i.size FROM Item i
			JOIN SlotItemLink l ON i.id = l.item
		    WHERE l.slot = $1 ORDER BY i.id
		`, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items = make([]entities.Item, 0)
	for rows.Next() {
		var item = entities.Item{}

		err = rows.Scan(
			&item.Id,
			&item.Name,
			&item.Description,
			&item.Size,
		)
		if err != nil {
			break
		}
		items = append(items, item)
	}

	if err != nil {
		return nil, err
	} else {
		err = rows.Err()
	}

	if err != nil {
		return nil, err
	}

	return items, nil
}

func (dao *SlotDAO) getSlot(id int) (*entities.Slot, error) {
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
