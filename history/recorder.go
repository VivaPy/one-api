package history

import "github.com/songquanpeng/one-api/relay/model"

type Recorder interface {
	Push(tokenID int, messages []model.Message, usage *Usage) error
	Pull(userID int) ([]model.Message, error)
}

func NewRecorder(typ string) Recorder {
	switch typ {
	case "rbmq":
		return NewMQRecorder()
	default:
		panic("no expected recorder type!")
	}
}
