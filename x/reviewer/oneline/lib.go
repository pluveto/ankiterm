package oneline

import (
	"fmt"
	"strconv"
	"strings"

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
		if err != nil {
			fmt.Println("No more cards.")
			return
		}

		fmt.Printf("\n[REVIEW MODE]\n")
		fmt.Println(format(card.Question))
		fmt.Println("\n[ENTER] Show Answer")

		awaitEnter()
		fmt.Print("\n---\n")
		fmt.Println(format(card.Answer))

		lookup := []string{"Again", "Hard", "Good", "Easy"}
		for i, button := range card.Buttons {
			fmt.Printf("[%d] %s (%s)\n", button, lookup[i], card.NextReviews[i])
		}

		action := awaitAction(am.CurrentCard().Buttons)
		switch code := action.getCode(); code {
		case ActionAbort:
			return
		case ActionSkip:
			continue
		case ActionAnswer:
			am.AnswerCard(action.(AnswerAction).CardEase)
		default:
			panic("unknown action code")
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
