package oneline

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pluveto/ankiterm/x/automata"
	"github.com/pluveto/ankiterm/x/xmisc"
	"github.com/pluveto/ankiterm/x/xslices"
)

const (
	ActionAnswer = "answer"
	ActionSkip   = "skip"
	ActionAbort  = "abort"
)

type Action interface {
	getCode() string
}

type AnswerAction struct {
	CardEase int
}

func (a AnswerAction) getCode() string {
	return ActionAnswer
}

type SkipAction struct {
}

func (a SkipAction) getCode() string {
	return ActionSkip
}

type AbortAction struct {
}

func (a AbortAction) getCode() string {
	return ActionAbort
}

func Execute(am *automata.Automata, deck string) {
	if am == nil {
		panic("am (automata.Automata) is nil")
	}
	if deck == "" {
		panic("deck is empty")
	}

	err := am.StartReview(deck)
	if err != nil {
		panic(err)
	}
	defer am.StopReview()

	for {
		card, err := am.NextCard()
		if err != nil {
			if strings.Contains(err.Error(), "Gui review is not currently active") {
				fmt.Println("Congratulations! You have finished all cards.")
				return
			}
			panic(err)
		}

		p := tea.NewProgram(initialModel(card))
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	}
}

func awaitEnter() {
	fmt.Scanln()
}

func awaitAction(validRange []int) Action {
	print("awaitAction")
	var input string
	fmt.Scanln(&input)
	// try parse int
	i, err := strconv.Atoi(input)
	if err == nil {
		if xslices.Contains(validRange, i) {
			return AnswerAction{CardEase: i}
		} else {
			fmt.Printf("invalid input \"%s\" out of range, try again: \n", input)
			return awaitAction(validRange)
		}
	}

	switch input {
	case "s":
		return SkipAction{}
	case "a":
		return AbortAction{}
	default:
		fmt.Printf("invalid input \"%s\", try again: \n", input)
		return awaitAction(validRange)
	}
}

func format(text string) string {
	text = xmisc.PurgeStyle(text)
	text = xmisc.TtyColor(text)
	return text
}
