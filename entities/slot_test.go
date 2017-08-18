package entities

import "testing"

func TestSlot_PutItem_SingleItem(t *testing.T) {
	var testData = []struct {
		slot             *Slot
		itemToPut        Item
		expectedContains int
		errIsNil         bool
		errMsg           string
	}{
		{
			slot:             NewSlot("", 100, true),
			itemToPut:        Item{Id: 100, Size: 50},
			expectedContains: 50,
			errIsNil:         true,
		},
		{
			slot:             NewSlot("", 10, true),
			itemToPut:        Item{Id: 100, Size: 50},
			expectedContains: 0,
			errIsNil:         false,
			errMsg:           ItemTooBig,
		},
		{
			slot:             NewSlot("", 100, false),
			itemToPut:        Item{Id: 100, Size: 50},
			expectedContains: 50,
			errIsNil:         true,
		},
	}

	for i, testItem := range testData {
		var err = testItem.slot.PutItem(testItem.itemToPut)
		if err == nil != testItem.errIsNil {
			t.Errorf("Expected err == nil %v, got %v (%d)", testItem.errIsNil, err == nil, i)
		}

		if err != nil && err.Error() != testItem.errMsg {
			t.Errorf("Expected \"%s\", got %s (%d)", testItem.errMsg, err.Error(), i)
		}

		if testItem.slot.contains != testItem.expectedContains {
			t.Errorf("Expected %d, got %d (%d)", testItem.expectedContains, testItem.slot.contains, i)
		}
	}
}

func TestSlot_PutItem_MultipleItem(t *testing.T) {
	var testData = []struct {
		slot             *Slot
		itemsToPut       []Item
		expectedContains int
	}{
		{
			slot: NewSlot("", 100, true),
			itemsToPut: []Item{
				{Id: 0, Size: 10},
				{Id: 1, Size: 20},
			},
			expectedContains: 30,
		},
		{
			slot: NewSlot("", 100, true),
			itemsToPut: []Item{
				{Id: 0, Size: 10},
				{Id: 1, Size: 20},
				{Id: 1, Size: 200},
			},
			expectedContains: 30,
		},
		{
			slot: NewSlot("", 100, true),
			itemsToPut: []Item{
				{Id: 0, Size: 10},
				{Id: 1, Size: 200},
				{Id: 1, Size: 20},
			},
			expectedContains: 30,
		},
	}

	for i, testItem := range testData {
		for _, item := range testItem.itemsToPut {
			testItem.slot.PutItem(item)
		}

		if testItem.slot.contains != testItem.expectedContains {
			t.Errorf("Expected %d, got %d (%d)", testItem.expectedContains, testItem.slot.contains, i)
		}
	}
}

func TestSlot_GetItem(t *testing.T) {
	var testData = []struct {
		slot             *Slot
		itemsToPut       []Item
		idToGet          int
		expectedContains int
	}{
		{
			slot: NewSlot("", 100, true),
			itemsToPut: []Item{
				{Id: 0, Size: 10},
				{Id: 1, Size: 20},
			},
			idToGet:          0,
			expectedContains: 20,
		},
		{
			slot: NewSlot("", 100, true),
			itemsToPut: []Item{
				{Id: 0, Size: 10},
				{Id: 1, Size: 20},
				{Id: 1, Size: 200},
			},
			idToGet:          0,
			expectedContains: 20,
		},
		{
			slot: NewSlot("", 100, true),
			itemsToPut: []Item{
				{Id: 0, Size: 10},
				{Id: 1, Size: 200},
				{Id: 1, Size: 20},
			},
			idToGet:          100,
			expectedContains: 30,
		},
	}

	for i, testItem := range testData {
		for _, item := range testItem.itemsToPut {
			testItem.slot.PutItem(item)
		}
		testItem.slot.GetItem(testItem.idToGet)

		if testItem.slot.contains != testItem.expectedContains {
			t.Errorf("Expected %d, got %d (%d)", testItem.expectedContains, testItem.slot.contains, i)
		}
	}
}

func TestSlot_MoveItem(t *testing.T) {
	var testData = []struct {
		slot1      *Slot
		slot2      *Slot
		itemsToPut []Item
		idToMove   int
		contains1  int
		contains2  int
	}{
		{
			slot1: NewSlot("", 100, true),
			slot2: NewSlot("", 0, true),
			itemsToPut: []Item{
				{Id: 0, Size: 10},
				{Id: 1, Size: 20},
			},
			idToMove:  0,
			contains1: 30,
			contains2: 0,
		},
		{
			slot1: NewSlot("", 100, true),
			slot2: NewSlot("", 100, true),
			itemsToPut: []Item{
				{Id: 0, Size: 10},
				{Id: 1, Size: 20},
				{Id: 1, Size: 200},
			},
			idToMove:  0,
			contains1: 20,
			contains2: 10,
		},
		{
			slot1: NewSlot("", 100, true),
			slot2: NewSlot("", 100, true),
			itemsToPut: []Item{
				{Id: 0, Size: 10},
				{Id: 1, Size: 200},
				{Id: 1, Size: 20},
			},
			idToMove:  100,
			contains1: 30,
			contains2: 0,
		},
	}

	for i, testItem := range testData {
		for _, item := range testItem.itemsToPut {
			testItem.slot1.PutItem(item)
		}
		testItem.slot1.MoveItem(testItem.idToMove, testItem.slot2)

		if testItem.slot1.contains != testItem.contains1 {
			t.Errorf("Expected %d, got %d (%d_1)", testItem.contains1, testItem.slot1.contains, i)
		}

		if testItem.slot2.contains != testItem.contains2 {
			t.Errorf("Expected %d, got %d (%d_2)", testItem.contains2, testItem.slot2.contains, i)
		}
	}
}
