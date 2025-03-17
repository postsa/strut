package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/postsa/strut-cli/internal/commands"
	"github.com/postsa/strut-cli/internal/history"
	"github.com/postsa/strut-cli/internal/input"
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
		inputCmd   tea.Cmd
	)
	var (
		historyModel tea.Model
		viewerModel  tea.Model
		inputModel   tea.Model
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
		}

	case messages.ExecutePromptMessage:
		currentProgressWidth := m.progress.Width
		m.progress = progress.New(progress.WithDefaultGradient())
		m.progress.Width = currentProgressWidth
		m.progress.SetPercent(0)
		cmd := m.progress.IncrPercent(.1)
		m.loading = true
		cmds = append(cmds, cmd, commands.TickCmd(), commands.FetchResponseCmd(m.client, msg.Prompt))

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - 6
	
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
		m.historyModel = m.historyModel.Focus()
		m.viewerModel = m.viewerModel.Blur()
		m.inputModel = m.inputModel.Blur()

	}
	if m.viewing {
		m.historyModel = m.historyModel.Blur()
		m.viewerModel = m.viewerModel.Focus()
		m.inputModel = m.inputModel.Focus()
	}

	inputModel, inputCmd = m.inputModel.Update(msg)
	m.inputModel = inputModel.(input.Model)

	historyModel, historyCmd = m.historyModel.Update(msg)
	m.historyModel = historyModel.(history.Model)

	viewerModel, viewerCmd = m.viewerModel.Update(msg)
	m.viewerModel = viewerModel.(viewer.Model)

	cmds = append(cmds, tiCmd, vpCmd, historyCmd, viewerCmd, inputCmd)
	return m, tea.Batch(cmds...)
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}
