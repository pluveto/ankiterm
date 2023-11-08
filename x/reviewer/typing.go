package reviewer

import (
	"strconv"
)

const (
	ActionAnswer = "answer"
	ActionSkip   = "skip"
	ActionAbort  = "abort"
)

type Action interface {
	GetCode() string
}

type AnswerAction struct {
	CardEase int
}

func (a AnswerAction) GetCode() string {
	return ActionAnswer
}

type SkipAction struct {
}

func (a SkipAction) GetCode() string {
	return ActionSkip
}

type AbortAction struct {
}

func (a AbortAction) GetCode() string {
	return ActionAbort
}

func ActionFromString(input string) Action {
	i, err := strconv.Atoi(input)
	if err == nil {
		return AnswerAction{CardEase: i}
	}

	switch input {
	case "s":
		return SkipAction{}
	case "a":
		return AbortAction{}
	default:
		return nil
	}
}
