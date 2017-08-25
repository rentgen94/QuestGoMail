package management

import "github.com/rentgen94/QuestGoMail/entities"

const (
	GetRoomCode = iota
	GetSlotsCode
	GetBagCode
	GetDoorsCode
	GetItemsCode
	GetInteractivesCode

	enterCode
	interactCode
	takeCode
	putCode
)

type Command struct {
	TypeCode int `json:"code"`
	ItemKey  int `json:"item_key"`
	Args     []string `json:"args"`
	Items    []entities.Item `json:"items"`
}

func NewCommand(typeCode int, itemKey int, args []string, items []entities.Item) Command {
	return Command{
		TypeCode: typeCode,
		ItemKey:  itemKey,
		Args:     args,
		Items:    items,
	}
}
