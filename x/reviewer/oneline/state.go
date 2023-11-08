package oneline

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/pluveto/ankiterm/x/ankicc"
	"github.com/pluveto/ankiterm/x/reviewer"
)

// model
// [R] <content> 1. <choice1> 2. <choice2> 3. <choice3> 4. <choice4>
// /                                            ^cursor     |||||selected
type model struct {
	question    string
	answer      string
	answerShown bool
	choices     []string
	chosen      int
	action      reviewer.Action
}

func initialModel(card *ankicc.CurrentCard) model {
	buttons := []string{}
	lookup := []string{"Again", "Hard", "Good", "Easy"}
	for i, button := range card.Buttons {
		button := fmt.Sprintf("%d. %s (%s)", button, lookup[i], card.NextReviews[i])
		buttons = append(buttons, button)
	}
	return model{
		question:    strings.ReplaceAll(format(card.Question), "\n", " "),
		answer:      strings.ReplaceAll(format(card.Answer), "\n", " "),
		choices:     buttons,
		answerShown: false,
		chosen:      0,
	}
}

func (m model) View() string {
	if !m.answerShown {
		return fmt.Sprintf("[R] %s %s", m.question, color.HiBlueString("Show Answer"))
	}
	choicesStr := ""
	for i, choice := range m.choices {
		choiceStr := ""
		if i == m.chosen {
			choiceStr = color.HiBlueString(choice)
		} else {
			choiceStr = choice
		}

		choicesStr += choiceStr + " "
	}
	s := fmt.Sprintf("[R] %s -- %s", m.answer, choicesStr)
	return s
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.answerShown {
		return updateWhenAnswerHidden(m, msg)
	}
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "left":
			m.chosen--
			if m.chosen < 0 {
				m.chosen = 0
			}
		case "right":
			m.chosen++
			if m.chosen >= len(m.choices) {
				m.chosen = len(m.choices) - 1
			}
		case "1", "2", "3", "4":
			m.chosen = int(msg.String()[0] - '1')
			m.action = reviewer.AnswerAction{
				CardEase: m.chosen + 1,
			}
			return m, tea.Quit
		case "enter":
			m.action = reviewer.AnswerAction{
				CardEase: m.chosen + 1,
			}
			return m, tea.Quit
		case "q", "ctrl+c":
			m.action = reviewer.AbortAction{}
			return m, tea.Quit
		case "s":
			m.action = reviewer.SkipAction{}
			return m, tea.Quit
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func updateWhenAnswerHidden(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	// check what msg it is
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.answerShown = true
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}
