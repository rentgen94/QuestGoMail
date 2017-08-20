package entities

const (
	GameContinue = iota
	GameFinished
)

type InteractionResult struct {
	Code int
	Msg  string
}

func ContinueResult(msg string) InteractionResult {
	return InteractionResult{
		Code: GameContinue,
		Msg:  msg,
	}
}

func FinishResult(msg string) InteractionResult {
	return InteractionResult{
		Code: GameFinished,
		Msg:  msg,
	}
}

func (res InteractionResult) IsFinish() bool {
	return res.Code == GameFinished
}
