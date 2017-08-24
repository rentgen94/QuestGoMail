package dao

import (
	"testing"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"github.com/stretchr/testify/assert"
	"errors"
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
