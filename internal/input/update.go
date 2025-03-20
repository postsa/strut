package input

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/postsa/strut/internal/commands"
)

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.textinput.Focused() {
				prompt := m.textinput.Value()
				m.textinput.Reset()
				cmds = append(cmds, commands.ExecutePromptCmd(prompt))
			}
		}
	}
	if m.inFocus {
		m.textinput.Focus()
		m.textinput.TextStyle = lipgloss.NewStyle().Background(lipgloss.Color("89"))
		m.textinput.PromptStyle = lipgloss.NewStyle().Background(lipgloss.Color("89"))
		m.textinput.PlaceholderStyle = lipgloss.NewStyle().Background(lipgloss.Color("89")).Foreground(lipgloss.Color("228"))

	} else {
		m.textinput.Blur()
		m.textinput.TextStyle = lipgloss.NewStyle().Background(lipgloss.Color("238"))
		m.textinput.PromptStyle = lipgloss.NewStyle().Background(lipgloss.Color("238"))
		m.textinput.PlaceholderStyle = lipgloss.NewStyle().Background(lipgloss.Color("238")).Foreground(lipgloss.Color("238"))
	}
	if m.inFocus {
		m.textinput, cmd = m.textinput.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}
