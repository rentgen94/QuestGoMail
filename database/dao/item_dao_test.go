package dao

import (
	"testing"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"github.com/stretchr/testify/assert"
	"errors"
)

func TestItemDAO_GetByName_Success(t *testing.T) {
	var db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()
	var rows = sqlmock.NewRows([]string{"id", "name", "description", "size"}).
		AddRow(0, "sword", "big", 100)

	mock.
	ExpectQuery("SELECT *").
	WithArgs("sword").
	WillReturnRows(rows)

	var itemDao = NewItemDAO(db)
	var item, itemErr = itemDao.GetByName("sword")

	assert.Nil(t, itemErr)

	assert.Equal(t, item.Id, 0)
	assert.Equal(t, item.Name, "sword")
	assert.Equal(t, item.Description, "big")
	assert.Equal(t, item.Size, 100)
}

func TestItemDAO_GetByName_Empty(t *testing.T) {
	var db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()
	var rows = sqlmock.NewRows([]string{"id", "name", "description", "size"})

	mock.
	ExpectQuery("SELECT *").
		WithArgs("sword").
		WillReturnRows(rows)

	var itemDao = NewItemDAO(db)
	var _, itemErr = itemDao.GetByName("sword")

	assert.NotNil(t, itemErr)
}

func TestItemDAO_GetByName_DBFail(t *testing.T) {
	var db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	mock.
	ExpectQuery("SELECT *").
		WithArgs("sword").
		WillReturnError(errors.New("DBError"))

	var itemDao = NewItemDAO(db)
	var _, itemErr = itemDao.GetByName("sword")

	assert.NotNil(t, itemErr)
	assert.Equal(t, itemErr.Error(), "DBError")
}
