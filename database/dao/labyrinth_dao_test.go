package dao

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

func TestLabyrinthDAO_getEmptyLabyrinthInfo_Success(t *testing.T) {
	var db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var rows = sqlmock.NewRows([]string{"id", "name", "start_room_id"}).
		AddRow(0, "a", 100)

	mock.
		ExpectQuery("SELECT ").
		WithArgs(1).
		WillReturnRows(rows)

	var labyrinthDAO = NewLabyrinthDAO(db)
	var labInfo, infoErr = labyrinthDAO.getEmptyLabyrinthInfo(1)

	assert.Nil(t, infoErr)
	assert.Equal(t, labInfo.id, 0)
	assert.Equal(t, labInfo.name, "a")
	assert.Equal(t, labInfo.startRoomId, 100)
}

func TestLabyrinthDAO_getEmptyLabyrinthInfo_Empty(t *testing.T) {
	var db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var rows = sqlmock.NewRows([]string{"id", "name", "start_room_id"})

	mock.
		ExpectQuery("SELECT ").
		WithArgs(1).
		WillReturnRows(rows)

	var labyrinthDAO = NewLabyrinthDAO(db)
	var _, infoErr = labyrinthDAO.getEmptyLabyrinthInfo(1)

	assert.NotNil(t, infoErr)
}

func TestLabyrinthDAO_getEmptyLabyrinthInfo_DBErr(t *testing.T) {
	var db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.
		ExpectQuery("SELECT ").
		WithArgs(1).
		WillReturnError(errors.New("err"))

	var labyrinthDAO = NewLabyrinthDAO(db)
	var _, infoErr = labyrinthDAO.getEmptyLabyrinthInfo(1)

	assert.Error(t, infoErr, "err")
}
