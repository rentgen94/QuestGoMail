package dao

import (
	"errors"
	"fmt"
	"github.com/rentgen94/QuestGoMail/entities"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

func TestActionDAO_getRelatedDoors_Success(t *testing.T) {
	var db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var rows = sqlmock.NewRows([]string{"id", "roomId", "state"}).
		AddRow(0, 0, true).
		AddRow(1, 1, false)

	mock.
		ExpectQuery("SELECT *").
		WithArgs(1).
		WillReturnRows(rows)

	var actionDao = NewActionDAO(db)
	var doorsInfo, doorErr = actionDao.getRelatedDoorsInfo(1)

	assert.Nil(t, doorErr)

	var correctData = []entityData{
		{0, 0, true},
		{1, 1, false},
	}

	for i := range correctData {
		assert.Equal(t, correctData[i], doorsInfo[i])
	}
}

func TestActionDAO_getRelatedDoors_Empty(t *testing.T) {
	var db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var rows = sqlmock.NewRows([]string{"id", "roomId", "state"})

	mock.
		ExpectQuery("SELECT *").
		WithArgs(1).
		WillReturnRows(rows)

	var actionDao = NewActionDAO(db)
	var doorInfo, doorErr = actionDao.getRelatedDoorsInfo(1)

	assert.Nil(t, doorErr)
	assert.Empty(t, doorInfo)
}

func TestActionDAO_getRelatedDoors_DBErr(t *testing.T) {
	var db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.
		ExpectQuery("SELECT *").
		WithArgs(1).
		WillReturnError(errors.New("err"))

	var actionDao = NewActionDAO(db)
	var _, doorErr = actionDao.getRelatedDoorsInfo(1)

	assert.Error(t, doorErr, "err")
}

func TestActionDAO_getActFunc_Success(t *testing.T) {
	var db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var slotsRows = sqlmock.NewRows([]string{"id", "roomId", "state"}).
		AddRow(0, 0, true).
		AddRow(1, 0, false)

	var interactivesRows = sqlmock.NewRows([]string{"id", "roomId", "state"}).
		AddRow(2, 0, true).
		AddRow(3, 0, false)

	var doorsRows = sqlmock.NewRows([]string{"id", "roomId", "state"}).
		AddRow(4, 0, true).
		AddRow(5, 0, false)

	mock.
		ExpectQuery("SELECT *").
		WithArgs(1).
		WillReturnRows(slotsRows)

	mock.
		ExpectQuery("SELECT *").
		WithArgs(1).
		WillReturnRows(interactivesRows)

	mock.
		ExpectQuery("SELECT *").
		WithArgs(1).
		WillReturnRows(doorsRows)

	var actionDao = NewActionDAO(db)

	var resultMsg = "Result message"
	var resultCode = 10

	var room = entities.NewRoom(0, "Some", "Descr")
	room.Slots()[0] = entities.NewSlot(0, "Slot", 100, true)
	room.Slots()[1] = entities.NewSlot(1, "Slot", 100, true)
	room.Interactives()[2] = entities.NewInteractiveObject(
		2,
		"Int",
		"Descr",
		true,
		func(args []string, items []entities.Item) error {
			return nil
		},
		nil,
	)
	room.Interactives()[3] = entities.NewInteractiveObject(
		3,
		"Int",
		"Descr",
		true,
		func(args []string, items []entities.Item) error {
			return nil
		},
		nil,
	)
	room.Doors()[4] = entities.NewDoor(4, "Some", false)
	room.Doors()[5] = entities.NewDoor(5, "Any", false)

	var lab = entities.NewLabyrinth(0, "lab")
	lab.SetRooms(entities.RoomsType{0: room})
	lab.SetStartRoom(room)

	var actFunc, funcErr = actionDao.getActFunc(1, resultMsg, resultCode)
	assert.Nil(t, funcErr)

	var res, actErr = actFunc(lab)
	assert.Nil(t, actErr)

	assert.Equal(t, entities.InteractionResult{Code: resultCode, Msg: resultMsg}, res)
}

func TestActionDAO_getActFunc_NoRoom(t *testing.T) {
	var db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var slotsRows = sqlmock.NewRows([]string{"id", "roomId", "state"}).
		AddRow(0, 0, true).
		AddRow(1, 0, false)

	var interactivesRows = sqlmock.NewRows([]string{"id", "roomId", "state"}).
		AddRow(2, 0, true).
		AddRow(3, 0, false)

	var doorsRows = sqlmock.NewRows([]string{"id", "roomId", "state"}).
		AddRow(4, 0, true).
		AddRow(5, 0, false)

	mock.
		ExpectQuery("SELECT *").
		WithArgs(1).
		WillReturnRows(slotsRows)

	mock.
		ExpectQuery("SELECT *").
		WithArgs(1).
		WillReturnRows(interactivesRows)

	mock.
		ExpectQuery("SELECT *").
		WithArgs(1).
		WillReturnRows(doorsRows)

	var actionDao = NewActionDAO(db)

	var resultMsg = "Result message"
	var resultCode = 10

	var lab = entities.NewLabyrinth(0, "lab")

	var actFunc, funcErr = actionDao.getActFunc(1, resultMsg, resultCode)
	assert.Nil(t, funcErr)

	var _, actErr = actFunc(lab)
	assert.Error(t, actErr, fmt.Sprintf(roomDoesNotExistsTemplate, 0))
}

func TestActionDAO_getActFunc_NoSlot(t *testing.T) {
	var db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var slotsRows = sqlmock.NewRows([]string{"id", "roomId", "state"}).
		AddRow(0, 0, true).
		AddRow(1, 0, false)

	var interactivesRows = sqlmock.NewRows([]string{"id", "roomId", "state"}).
		AddRow(2, 0, true).
		AddRow(3, 0, false)

	var doorsRows = sqlmock.NewRows([]string{"id", "roomId", "state"}).
		AddRow(4, 0, true).
		AddRow(5, 0, false)

	mock.
		ExpectQuery("SELECT *").
		WithArgs(1).
		WillReturnRows(slotsRows)

	mock.
		ExpectQuery("SELECT *").
		WithArgs(1).
		WillReturnRows(interactivesRows)

	mock.
		ExpectQuery("SELECT *").
		WithArgs(1).
		WillReturnRows(doorsRows)

	var actionDao = NewActionDAO(db)

	var resultMsg = "Result message"
	var resultCode = 10

	var room = entities.NewRoom(0, "Some", "Descr")
	room.Slots()[0] = entities.NewSlot(0, "Slot", 100, true)
	room.Interactives()[2] = entities.NewInteractiveObject(
		2,
		"Int",
		"Descr",
		true,
		func(args []string, items []entities.Item) error {
			return nil
		},
		nil,
	)
	room.Interactives()[3] = entities.NewInteractiveObject(
		3,
		"Int",
		"Descr",
		true,
		func(args []string, items []entities.Item) error {
			return nil
		},
		nil,
	)
	room.Doors()[4] = entities.NewDoor(4, "Some", false)
	room.Doors()[5] = entities.NewDoor(5, "Any", false)

	var lab = entities.NewLabyrinth(0, "lab")
	lab.SetRooms(entities.RoomsType{1: room})
	lab.SetStartRoom(room)

	var actFunc, funcErr = actionDao.getActFunc(1, resultMsg, resultCode)
	assert.Nil(t, funcErr)

	var _, actErr = actFunc(lab)
	assert.Error(t, actErr, fmt.Sprintf(slotDoesNotExistsTemplate, 1, 0))
}
