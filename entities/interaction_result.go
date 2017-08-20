package entities

const (
	GameContinue = iota
	GameFinished
)

type InteractionResult struct {
	code int
	Msg  string
}

func ContinueResult(msg string) InteractionResult {
	return InteractionResult{
		code: GameContinue,
		Msg:  msg,
	}
}

func FinishResult(msg string) InteractionResult {
	return InteractionResult{
		code: GameFinished,
		Msg:  msg,
	}
}

func (res InteractionResult) IsFinish() bool {
	return res.code == GameFinished
}
