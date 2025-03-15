package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func getAnswerDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		position := m.Index()
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				return getAnswerCmd(position)
			}
		}
		return nil
	}
	return d
}

type GetAnswerMsg struct{ position int }

func getAnswerMsg(position int) tea.Msg {
	return GetAnswerMsg{position}
}

func getAnswerCmd(position int) tea.Cmd {
	return func() tea.Msg {
		return getAnswerMsg(position)
	}
}
