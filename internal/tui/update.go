package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/postsa/strut-cli/internal/commands"
	"github.com/postsa/strut-cli/internal/history"
	"github.com/postsa/strut-cli/internal/messages"
	"github.com/postsa/strut-cli/internal/viewer"
	"os"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd      tea.Cmd
		vpCmd      tea.Cmd
		historyCmd tea.Cmd
		viewerCmd  tea.Cmd
	)
	var (
		historyModel tea.Model
		viewerModel  tea.Model
	)

	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.QuitMsg:
		return m, nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			if m.viewing {
				m.listFocus = true
				m.viewing = false
			} else {
				m.viewing = true
				m.listFocus = false
			}
		case tea.KeyCtrlC:
			m.quitting = true
			return m, tea.Quit
		case tea.KeyEsc:
			return m, nil
		case tea.KeyEnter:
			if m.textinput.Focused() {
				prompt := m.textinput.Value()
				currentProgressWidth := m.progress.Width
				m.progress = progress.New(progress.WithDefaultGradient())
				m.progress.Width = currentProgressWidth
				m.progress.SetPercent(0)
				cmd := m.progress.IncrPercent(.1)
				m.loading = true
				m.textinput.Reset()
				return m, tea.Batch(cmd, commands.TickCmd(), commands.FetchResponseCmd(m.client, prompt))
			}
		}
	case tea.WindowSizeMsg:
		m.textinput.Blur()
		m.progress.Width = msg.Width - 6
		m.textinput.Focus()

	case messages.HistoryResizedMessage:
		return m, commands.ViewPortResizeCmd(msg.TotalWidth - msg.NewWidth)

	case messages.GeminiResponseMsg:
		m.geminiResponse = msg.Response
		m.response = fmt.Sprintf("%v", msg.Response.Candidates[0].Content.Parts[0])
		cmds = append(cmds, commands.NewAnswerCmd(m.response, msg.Prompt))

	case messages.NewRenderMessage:
		m.loading = false

	case messages.SetAnswerMessage:
		m.viewing = true
		m.listFocus = false

	case messages.ErrMsg:
		m.err = msg.Err

	case messages.TickMsg:
		cmd := m.progress.IncrPercent(((1 - m.progress.Percent()) / 3) * ((1 - m.progress.Percent()) / 1.2))
		return m, tea.Batch(commands.TickCmd(), cmd)

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	case messages.EditorFinishedMsg:
		defer os.Remove(msg.File.Name())
		if msg.Err != nil {
			m.err = msg.Err
			return m, tea.Quit
		}
	}

	if m.listFocus {
		m.textinput.Blur()
		m.historyModel = m.historyModel.Focus()
		m.viewerModel = m.viewerModel.Blur()
		m.textinput.TextStyle = lipgloss.NewStyle().Background(lipgloss.Color("238"))
		m.textinput.PromptStyle = lipgloss.NewStyle().Background(lipgloss.Color("238"))
		m.textinput.PlaceholderStyle = lipgloss.NewStyle().Background(lipgloss.Color("238")).Foreground(lipgloss.Color("238"))
	}
	if m.viewing {
		m.historyModel = m.historyModel.Blur()
		m.viewerModel = m.viewerModel.Focus()
		m.textinput.TextStyle = lipgloss.NewStyle().Background(lipgloss.Color("89"))
		m.textinput.PromptStyle = lipgloss.NewStyle().Background(lipgloss.Color("89"))
		m.textinput.PlaceholderStyle = lipgloss.NewStyle().Background(lipgloss.Color("89")).Foreground(lipgloss.Color("228"))
		m.textinput.Focus()
	}

	m.textinput, tiCmd = m.textinput.Update(msg)

	historyModel, historyCmd = m.historyModel.Update(msg)
	m.historyModel = historyModel.(history.HistoryModel)

	viewerModel, viewerCmd = m.viewerModel.Update(msg)
	m.viewerModel = viewerModel.(viewer.ViewerModel)

	cmds = append(cmds, tiCmd, vpCmd, historyCmd, viewerCmd)
	return m, tea.Batch(cmds...)
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}
