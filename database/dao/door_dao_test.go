package dao

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

func TestDoorDAO_GetById_Success(t *testing.T) {
	var db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var rows = sqlmock.NewRows([]string{"id", "name", "isAccessible"}).
		AddRow(1, "door", true)

	mock.
		ExpectQuery("SELECT ").
		WithArgs(1).
		WillReturnRows(rows)

	var doorDao = NewDoorDAO(db)
	var door, doorErr = doorDao.GetById(1)

	assert.Nil(t, doorErr)
	assert.Equal(t, door.Id(), 1)
	assert.Equal(t, door.Name(), "door")
	assert.Equal(t, door.IsAccessible(), true)
}

func TestDoorDAO_GetById_NotFound(t *testing.T) {
	var db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var rows = sqlmock.NewRows([]string{"id", "name", "isAccessible"})

	mock.
		ExpectQuery("SELECT ").
		WithArgs(1).
		WillReturnRows(rows)

	var doorDao = NewDoorDAO(db)
	var _, doorErr = doorDao.GetById(1)

	assert.NotNil(t, doorErr)
}

func TestDoorDAO_GetById_DBErr(t *testing.T) {
	var db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.
		ExpectQuery("SELECT ").
		WithArgs(1).
		WillReturnError(errors.New("err"))

	var doorDao = NewDoorDAO(db)
	var _, doorErr = doorDao.GetById(1)

	assert.Error(t, doorErr, "err")
}

func TestDoorDAO_getDoorInfo_Success(t *testing.T) {
	var db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var rows = sqlmock.NewRows([]string{"id", "name", "isAccessible", "id1", "id2"}).
		AddRow(1, "door", true, 1, 2)

	mock.
		ExpectQuery("SELECT ").
		WithArgs(1).
		WillReturnRows(rows)

	var doorDao = NewDoorDAO(db)
	var doorInfo, infoErr = doorDao.getDoorInfo(1)

	assert.Nil(t, infoErr)
	assert.Equal(t, doorInfo.room1Id, 1)
	assert.Equal(t, doorInfo.room2Id, 2)
	assert.Equal(t, doorInfo.door.Id(), 1)
	assert.Equal(t, doorInfo.door.Name(), "door")
	assert.Equal(t, doorInfo.door.IsAccessible(), true)
}

func TestDoorDAO_getDoorInfo_NotFound(t *testing.T) {
	var db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var rows = sqlmock.NewRows([]string{"id", "name", "isAccessible", "id1", "id2"})

	mock.
		ExpectQuery("SELECT ").
		WithArgs(1).
		WillReturnRows(rows)

	var doorDao = NewDoorDAO(db)
	var _, infoErr = doorDao.getDoorInfo(1)

	assert.NotNil(t, infoErr)
}
