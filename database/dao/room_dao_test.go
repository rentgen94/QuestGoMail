package dao

import (
	"errors"
	"github.com/rentgen94/QuestGoMail/entities"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

func TestRoomDAO_getRelatedEntityIds_Success(t *testing.T) {
	var db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var rows = sqlmock.NewRows([]string{"id"}).
		AddRow(0).
		AddRow(1).
		AddRow(2)

	mock.
		ExpectQuery("SELECT ").
		WithArgs(1).
		WillReturnRows(rows)

	var roomDao = NewRoomDAO(db)
	var ids, idsErr = roomDao.getRelatedEntityIds(getRoomRelatedSlotIdsQuery, 1)

	assert.Nil(t, idsErr)
	assert.Equal(t, len(ids), 3)

	for i, id := range ids {
		assert.Equal(t, i, id)
	}
}

func TestRoomDAO_getRelatedEntityIds_Empty(t *testing.T) {
	var db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var rows = sqlmock.NewRows([]string{"id"})

	mock.
		ExpectQuery("SELECT ").
		WithArgs(1).
		WillReturnRows(rows)

	var roomDao = NewRoomDAO(db)
	var ids, idsErr = roomDao.getRelatedEntityIds(getRoomRelatedSlotIdsQuery, 1)

	assert.Nil(t, idsErr)
	assert.Equal(t, len(ids), 0)
}

func TestRoomDAO_getRelatedEntityIds_DBErr(t *testing.T) {
	var db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.
		ExpectQuery("SELECT ").
		WithArgs(1).
		WillReturnError(errors.New("err"))

	var roomDao = NewRoomDAO(db)
	var _, idsErr = roomDao.getRelatedEntityIds(getRoomRelatedSlotIdsQuery, 1)

	assert.Error(t, idsErr, "err")
}

func TestRoomDAO_getEmptyRoom_Success(t *testing.T) {
	var db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var rows = sqlmock.NewRows([]string{"id", "name", "description"}).
		AddRow(1, "room", "big")

	mock.
		ExpectQuery("SELECT ").
		WithArgs(1).
		WillReturnRows(rows)

	var roomDao = NewRoomDAO(db)
	var room, roomErr = roomDao.getEmptyRoom(1)

	assert.Nil(t, roomErr)

	var correctRoom = entities.NewRoom(1, "room", "big")
	assert.Equal(t, correctRoom, room)
}

func TestRoomDAO_getEmptyRoom_NotFound(t *testing.T) {
	var db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var rows = sqlmock.NewRows([]string{"id", "name", "description"})

	mock.
		ExpectQuery("SELECT ").
		WithArgs(1).
		WillReturnRows(rows)

	var roomDao = NewRoomDAO(db)
	var _, roomErr = roomDao.getEmptyRoom(1)

	assert.NotNil(t, roomErr)
}

func TestRoomDAO_getEmptyRoom_DBErr(t *testing.T) {
	var db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.
		ExpectQuery("SELECT ").
		WithArgs(1).
		WillReturnError(errors.New("err"))

	var roomDao = NewRoomDAO(db)
	var _, roomErr = roomDao.getEmptyRoom(1)

	assert.Error(t, roomErr, "err")
}
