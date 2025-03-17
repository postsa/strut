package history

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func changeAnswerDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		position := m.Index()
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				return changeAnswerCommand(position)
			}
		}
		return nil
	}
	return d
}

type ChangeAnswerMessage struct{ position int }

func changeAnswerCommand(position int) tea.Cmd {
	return func() tea.Msg {
		return ChangeAnswerMessage{position}
	}
}
