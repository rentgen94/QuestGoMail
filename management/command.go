package management

import "github.com/rentgen94/QuestGoMail/entities"

const (
	getDoorsCode = iota
	getSlotsCode
	getItemsCode
	getInteractivesCode

	enterCode
	interactCode
	takeCode
	putCode
)

type Command struct {
	typeCode int
	itemKey  string
	args     []string
	items    []entities.Item
}

func NewCommand(typeCode int, itemKey string, args []string, items []entities.Item) Command {
	return Command{
		typeCode: typeCode,
		itemKey:  itemKey,
		args:     args,
		items:    items,
	}
}
