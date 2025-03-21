package viewer

import (
	"fmt"
	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/postsa/strut/internal/commands"
	"github.com/postsa/strut/internal/messages"
	"strings"
)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlA:
			clipboard.WriteAll(m.currentContent)
		case tea.KeyCtrlS:
			clipboard.WriteAll(strings.Trim(m.currentContent, "`"))
		case tea.KeyCtrlE:
			return m, commands.ExecuteVim(m.currentContent)
		}
	case messages.NewAnswerMessage:
		m.currentContent = msg.Answer
		output, _ := m.renderer.Render(msg.Answer)
		cmds = append(cmds, commands.NewRenderCmd(output))

	case messages.NewRenderMessage:
		m.currentContentRendered = msg.Content
		m.pane.SetContent(m.currentContentRendered)
		m.pane.GotoTop()

	case messages.SetAnswerMessage:
		m.currentContentRendered = msg.AnswerRendered
		m.currentContent = msg.Answer
		m.pane.SetContent(m.currentContentRendered)
		m.pane.GotoTop()

	case tea.WindowSizeMsg:
		m.pane.Height = msg.Height - 9
		m.pane, cmd = m.pane.Update(msg)
		return m, cmd

	case messages.ViewPortResizeMessage:
		m.pane.Style.MaxWidth(msg.Width)
		m.pane.Width = msg.Width
		m.pane, cmd = m.pane.Update(msg)
		return m, cmd

	case messages.ErrMsg:
		m.pane.SetContent(fmt.Sprintf("Error: %s", msg.Err))
	}

	if m.inFocus {
		m.pane, cmd = m.pane.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}
