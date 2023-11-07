package automata

import (
	"errors"

	"github.com/pluveto/ankiterm/x/ankicc"
	"github.com/pluveto/ankiterm/x/xslices"
)


type Automata struct {
	client ankicc.Client
	deck *ankicc.DeckStat
	card *ankicc.CurrentCard
	reviewing bool
	needAnswer bool
}

func NewAutomata(client ankicc.Client) *Automata {
	return &Automata{
		client: client,
	}
}

func (m* Automata) CurrentCard() *ankicc.CurrentCard {
	return m.card
}

func (m* Automata) AllowReview() bool {
	return m.deck == nil
}

func (m* Automata) StartReview(deck string) error {
	if !m.AllowReview() {
		return errors.New("already reviewing")
	}
	deckStat, err := m.client.GetDeckStat(deck)
	if err != nil {
		return err
	}
	m.deck = &deckStat

	err = m.client.GuiDeckReview(deck)
	if err != nil {
		return err
	}
	m.reviewing = true
	return nil
}

func (m* Automata) NextCard() (card *ankicc.CurrentCard, err error) {
	if !m.reviewing {
		return nil, errors.New("not reviewing")
	}

	if m.needAnswer {
		return nil, errors.New("need answer first")
	}

	card, err = m.client.GuiCurrentCard()
	if err != nil {
		return nil, err
	}

	m.card = card
	m.needAnswer = true
	return card, nil
}

func (m* Automata) AnswerCard(ease int) (err error) {
	if !m.reviewing {
		return errors.New("not reviewing")
	}

	if !m.needAnswer {
		return errors.New("no need answer")
	}

	err = m.client.GuiShowAnswer()
	if err != nil {
		return err
	}

	if !xslices.Contains(m.card.Buttons, ease) {
		return errors.New("ease out of range")
	}

	err = m.client.GuiAnswerCard(ease)
	if err != nil {
		return err
	}

	m.needAnswer = false
	return nil
}

func (m* Automata) StopReview() error {
	if !m.reviewing {
		return errors.New("not reviewing")
	}

	m.reviewing = false
	m.needAnswer = false
	m.deck = nil
	m.card = nil
	return nil
}
