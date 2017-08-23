package dao

import (
	"errors"
	"fmt"
	"github.com/rentgen94/QuestGoMail/entities"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

func TestInteractiveDAO_getNecessaryItems_Success(t *testing.T) {
	var db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var itemIds = sqlmock.NewRows([]string{"id"}).
		AddRow(0).
		AddRow(1)

	var firstItemRows = sqlmock.NewRows([]string{"id", "name", "description", "size"}).
		AddRow(0, "sword", "big", 100)

	var secondItemRows = sqlmock.NewRows([]string{"id", "name", "description", "size"}).
		AddRow(1, "axe", "small", 100)

	mock.
		ExpectQuery("SELECT *").
		WithArgs(0).
		WillReturnRows(itemIds)

	mock.
		ExpectQuery("SELECT *").
		WithArgs(0).
		WillReturnRows(firstItemRows)

	mock.
		ExpectQuery("SELECT *").
		WithArgs(1).
		WillReturnRows(secondItemRows)

	var interactiveDao = NewInteractiveDAO(db)

	var items, itemErr = interactiveDao.getNecessaryItems(0)

	assert.Nil(t, itemErr)
	assert.Equal(t, 2, len(items))

	var expectedItems = []entities.Item{
		{Id: 0, Size: 100, Description: "big", Name: "sword"},
		{Id: 1, Size: 100, Description: "small", Name: "axe"},
	}

	for i := range expectedItems {
		assert.Equal(t, expectedItems[i], items[i])
	}
}

func TestInteractiveDAO_getNecessaryItems_DBItemsFail(t *testing.T) {
	var db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.
		ExpectQuery("SELECT *").
		WithArgs(0).
		WillReturnError(errors.New("err"))

	var interactiveDao = NewInteractiveDAO(db)

	var _, itemErr = interactiveDao.getNecessaryItems(0)

	assert.Error(t, itemErr, "err")
}

func TestInteractiveDAO_getNecessaryItems_DBItemFail(t *testing.T) {
	var db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var itemIds = sqlmock.NewRows([]string{"id"}).
		AddRow(0).
		AddRow(1)

	mock.
		ExpectQuery("SELECT *").
		WithArgs(0).
		WillReturnRows(itemIds)

	mock.
		ExpectQuery("SELECT *").
		WithArgs(0).
		WillReturnError(errors.New("err"))

	var interactiveDao = NewInteractiveDAO(db)

	var _, itemErr = interactiveDao.getNecessaryItems(0)

	assert.Error(t, itemErr, "err")
}

func TestInteractiveDAO_getInputChecker(t *testing.T) {
	var testData = []struct {
		expectedArgs  []string
		expectedItems []entities.Item
		args          []string
		items         []entities.Item
		err           error
	}{
		{
			[]string{},
			[]entities.Item{},
			[]string{"kj;kj"},
			[]entities.Item{},
			nil,
		},
		{
			[]string{"some"},
			[]entities.Item{},
			[]string{"some"},
			[]entities.Item{},
			nil,
		},
		{
			[]string{"some"},
			[]entities.Item{},
			[]string{"any"},
			[]entities.Item{},
			errors.New(gotWrongArguments),
		},
		{
			[]string{},
			[]entities.Item{
				{Id: 100},
				{Id: 200},
			},
			[]string{},
			[]entities.Item{
				{Id: 200},
				{Id: 100},
			},
			nil,
		},
		{
			[]string{},
			[]entities.Item{},
			[]string{},
			[]entities.Item{
				{Id: 200},
				{Id: 100},
			},
			nil,
		},
		{
			[]string{},
			[]entities.Item{
				{Id: 100},
				{Id: 200},
			},
			[]string{},
			[]entities.Item{
				{Id: 200},
			},
			errors.New(itemListLengthDoesNotMatch),
		},
		{
			[]string{},
			[]entities.Item{
				{Id: 100},
				{Id: 200},
			},
			[]string{},
			[]entities.Item{
				{Id: 200},
				{Id: 300},
			},
			errors.New(notFoundNecessaryItem),
		},
	}

	for i, item := range testData {
		var interactiveDAO = NewInteractiveDAO(nil)
		var checker = interactiveDAO.getInputChecker(item.expectedArgs, item.expectedItems)
		var res = checker(item.args, item.items)
		assert.Equal(t, item.err, res, fmt.Sprintf("i = %d", i))
	}
}

func TestInteractiveDAO_GetById_Success(t *testing.T) {
	var db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var interRows = sqlmock.NewRows([]string{"id", "room", "name", "description", "isAccessible", "args"}).
		AddRow(0, 0, "Some", "Description", true, []string{"alpha", "beta", "gamma"})

	var itemIdRows = sqlmock.NewRows([]string{"id"})

	mock.
		ExpectQuery("SELECT").
		WithArgs(0).
		WillReturnRows(interRows)

	mock.
		ExpectQuery("SELECT").
		WithArgs(0).
		WillReturnRows(itemIdRows)

	var interactiveDao = NewInteractiveDAO(db)
	var inter, iErr = interactiveDao.GetById(0)
	assert.Nil(t, iErr)

	var action = entities.NewAction(
		"",
		func(labyrinth *entities.Labyrinth) (result entities.InteractionResult, err error) {
			return entities.InteractionResult{Msg: "ui"}, nil
		},
	)

	inter.SetAction(action)
	var testCases = []struct {
		args []string
		err  error
	}{
		{
			[]string{},
			errors.New(argumentLengthDoesNotMatch),
		},
		{
			[]string{"alpha", "beta", "gamma"},
			nil,
		},
	}

	for i, testCase := range testCases {
		var _, e = inter.Interact(testCase.args, []entities.Item{})
		assert.Equal(t, testCase.err, e, fmt.Sprintf("i = %d", i))
	}
}
