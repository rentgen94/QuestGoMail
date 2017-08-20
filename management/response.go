package management

import "github.com/rentgen94/QuestGoMail/entities"

type Response struct {
	Result entities.InteractionResult
	ErrMsg string
	Data   interface{}
}

func NewResponse() Response {
	return Response{
		Result: entities.ContinueResult(""),
		ErrMsg: "",
		Data:   nil,
	}
}
