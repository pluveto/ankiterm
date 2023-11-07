package oneline

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pluveto/ankiterm/x/ankicc"
)

// model
// Q: <question> A: <answer> 1. <choice1> 2. <choice2> 3. <choice3> 4. <choice4>
// /                                            ^cursor     |||||selected
type model struct {
	question string
	answer   string
	choices  []string
	cursor   int
}

func initialModel(card *ankicc.CurrentCard) model {
	buttons := []string{}
	lookup := []string{"Again", "Hard", "Good", "Easy"}
	for i, button := range card.Buttons {
		button := fmt.Sprintf("%d. %s (%s)", button, lookup[i], card.NextReviews[i])
		buttons = append(buttons, button)
	}
	return model{
		question: format(card.Question),
		answer:   format(card.Answer),
		choices:  buttons,
		cursor:   0,
	}
}

func (m model) View() string {
	choicesStr := ""
	for i, choice := range m.choices {
		choiceStr := ""
		if i == m.cursor {
			choiceStr = fmt.Sprintf("<%s>", choice)
		} else {
			choiceStr = choice
		}

		choicesStr += choiceStr + " "
	}
	s := fmt.Sprintf("Q: %s -- A: %s -- %s", m.question, m.answer, choicesStr)
	return s
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		case "left":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = 0
			}
		case "right":
			m.cursor++
			if m.cursor >= len(m.choices) {
				m.cursor = len(m.choices) - 1
			}
		case "enter":
			println("enter")
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}
