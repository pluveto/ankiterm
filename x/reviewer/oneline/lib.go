package oneline

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pluveto/ankiterm/x/automata"
	"github.com/pluveto/ankiterm/x/reviewer"
	"github.com/pluveto/ankiterm/x/xmisc"
)

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
		m, _ := tea.NewProgram(initialModel(card)).Run()
		if err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
		out := m.(model)

		action := out.action
		if action == nil {
			// quit
			fmt.Println("Bye!")
			return
		}

		switch code := action.GetCode(); code {
		case reviewer.ActionAbort:
			return
		case reviewer.ActionSkip:
			continue
		case reviewer.ActionAnswer:
			am.AnswerCard(action.(reviewer.AnswerAction).CardEase)
		default:
			panic("unknown action code " + code)
		}
	}
}

func format(text string) string {
	text = xmisc.PurgeStyle(text)
	text = xmisc.TtyColor(text)
	return text
}
