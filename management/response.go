package management

import "github.com/rentgen94/QuestGoMail/entities"

type Response struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	ErrMsg string      `json:"errMsg"`
	Data   interface{} `json:"data"`
}

func NewResponse() Response {
	return Response{
		Code: entities.GameContinue,
	}
}

func (r Response) IsFinish() bool {
	return r.Code == entities.GameFinished
}
