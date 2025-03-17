package input

import (
	"github.com/charmbracelet/bubbles/textinput"
)

type Model struct {
	textinput textinput.Model
	inFocus   bool
}

func (m Model) Focus() Model {
	m.inFocus = true
	return m
}

func (m Model) Blur() Model {
	m.inFocus = false
	return m
}

func NewModel(modelName string) Model {
	ti := textinput.New()
	ti.Prompt = "(" + modelName + ")" + " > "
	ti.Placeholder = "ask a question ..."
	ti.Focus()

	return Model{
		textinput: ti,
		inFocus:   true,
	}
}
