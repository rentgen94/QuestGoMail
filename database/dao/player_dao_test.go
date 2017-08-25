package dao

import (
	"github.com/rentgen94/QuestGoMail/entities"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

func TestDbPlayerDAO_FindPlayer_Successful(t *testing.T) {
	var db, mock, err = sqlmock.New()

	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var rows = sqlmock.NewRows([]string{"id", "name", "password"}).
		AddRow(2, "qqq", "111")

	mock.
		ExpectQuery("SELECT *").
		WithArgs("qqq", "111").
		WillReturnRows(rows)

	var player = entities.Player{Login: "qqq", Password: "111"}

	var playerDAO = NewDBPlayerDAO(db)
	var founded, findErr = playerDAO.FindPlayer(&player)

	assert.Nil(t, findErr)

	assert.Equal(t, founded.Id, 2)
	assert.Equal(t, founded.Login, "qqq")
	assert.Equal(t, founded.Password, "111")
}

func TestDbPlayerDAO_CreatePlayer_Successful(t *testing.T) {
	var db, mock, err = sqlmock.New()

	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var rows = sqlmock.NewRows([]string{"cnt"}).AddRow(0)

	mock.ExpectQuery("SELECT *").
		WithArgs("qqq").
		WillReturnRows(rows)

	mock.ExpectExec("INSERT INTO").
		WithArgs("qqq", "111").
		WillReturnResult(sqlmock.NewResult(1, 1))

	var player = entities.Player{Login: "qqq", Password: "111"}
	var playerDAO = NewDBPlayerDAO(db)
	var createErr = playerDAO.CreatePlayer(&player)

	assert.Nil(t, createErr)
}

func TestDbPlayerDAO_CreatePlayer_AlreadyExist(t *testing.T) {
	var db, mock, err = sqlmock.New()

	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var rows = sqlmock.NewRows([]string{"cnt"}).
		AddRow(1)

	mock.ExpectQuery("SELECT *").
		WithArgs("qqq").
		WillReturnRows(rows)

	mock.ExpectExec("INSERT INTO").
		WithArgs("qqq", "111").
		WillReturnResult(sqlmock.NewResult(1, 1))

	var player = entities.Player{Login: "qqq", Password: "111"}
	var playerDAO = NewDBPlayerDAO(db)
	var createErr = playerDAO.CreatePlayer(&player)

	assert.NotNil(t, createErr)
}

func TestDbPlayerDAO_SelectAllPlayers(t *testing.T) {
	var db, mock, err = sqlmock.New()

	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var testData = []entities.Player{
		{Id: 1, Login: "qqq", Password: "111"},
		{Id: 2, Login: "test", Password: "123"},
		{Id: 3, Login: "hello", Password: "world"},
		{Id: 4, Login: "check", Password: "something"},
	}
	var rows = sqlmock.NewRows([]string{"id", "name", "password"}).
		AddRow(1, "qqq", "111").
		AddRow(2, "test", "123").
		AddRow(3, "hello", "world").
		AddRow(4, "check", "something")

	mock.ExpectQuery("SELECT *").WillReturnRows(rows)

	var playerDAO = NewDBPlayerDAO(db)
	var players, findErr = playerDAO.SelectAllPlayers()

	assert.Nil(t, findErr)

	for i, player := range players {
		assert.Equal(t, player.Id, testData[i].Id)
		assert.Equal(t, player.Login, testData[i].Login)
		assert.Equal(t, player.Password, testData[i].Password)
	}
}
