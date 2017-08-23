package dao

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

func TestSlotDAO_getSlot_Success(t *testing.T) {
	var db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()
	var rows = sqlmock.NewRows([]string{"id", "name", "capacity", "isAccessible"}).
		AddRow(1, "box", 100, true)

	mock.
		ExpectQuery("SELECT ").
		WithArgs(1).
		WillReturnRows(rows)

	var slotDao = NewSlotDAO(db)
	var item, itemErr = slotDao.getSlot(1)

	assert.Nil(t, itemErr)

	assert.Equal(t, item.Id(), 1)
	assert.Equal(t, item.Name(), "box")
	assert.Equal(t, item.Capacity(), 100)
	assert.Equal(t, item.IsAccessible(), true)
}

func TestSlotDAO_getSlot_Empty(t *testing.T) {
	var db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()
	var rows = sqlmock.NewRows([]string{"id", "name", "capacity", "isAccessible"})

	mock.
		ExpectQuery("SELECT ").
		WithArgs(1).
		WillReturnRows(rows)

	var slotDao = NewSlotDAO(db)
	var _, itemErr = slotDao.getSlot(1)

	assert.NotNil(t, itemErr)
}

func TestSlotDAO_getSlotItems_Success(t *testing.T) {
	var db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()
	var itemRows = sqlmock.NewRows([]string{"id", "name", "description", "size"}).
		AddRow(0, "sword", "big", 100).
		AddRow(1, "bow", "mid", 50).
		AddRow(2, "axe", "small", 10)
	var slotRows = sqlmock.NewRows([]string{"id", "name", "capacity", "isAccessible"}).
		AddRow(1, "box", 1000, true)

	mock.
		ExpectQuery("SELECT id").
		WithArgs(1).
		WillReturnRows(slotRows)

	mock.
		ExpectQuery("SELECT i.id").
		WithArgs(1).
		WillReturnRows(itemRows)

	var slotDAO = NewSlotDAO(db)
	var slot, itemErr = slotDAO.GetById(1)

	assert.Nil(t, itemErr)

	var cnt = 0
	for range slot.Items() {
		cnt++
	}

	assert.Equal(t, cnt, 3)
}

func TestSlotDAO_GetById_Overfilled(t *testing.T) {
	var db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()
	var itemRows = sqlmock.NewRows([]string{"id", "name", "description", "size"}).
		AddRow(0, "sword", "big", 100).
		AddRow(1, "bow", "mid", 50).
		AddRow(2, "axe", "small", 10)
	var slotRows = sqlmock.NewRows([]string{"id", "name", "capacity", "isAccessible"}).
		AddRow(1, "box", 100, true)

	mock.
		ExpectQuery("SELECT id").
		WithArgs(1).
		WillReturnRows(slotRows)

	mock.
		ExpectQuery("SELECT i.id").
		WithArgs(1).
		WillReturnRows(itemRows)

	var slotDAO = NewSlotDAO(db)
	var _, itemErr = slotDAO.GetById(1)

	assert.NotNil(t, itemErr)
}
