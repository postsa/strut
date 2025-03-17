package viewer

import (
	"fmt"
	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/postsa/strut-cli/internal/commands"
	"github.com/postsa/strut-cli/internal/messages"
	"strings"
)

func (m ViewerModel) Init() tea.Cmd {
	return nil
}

func (m ViewerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		windowUpdateCmd tea.Cmd
		cmds            []tea.Cmd
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
		m.pane, windowUpdateCmd = m.pane.Update(msg)
		return m, windowUpdateCmd

	case messages.ViewPortResizeMessage:
		m.pane.Style.MaxWidth(msg.Width)
		m.pane.Width = msg.Width
		m.pane, windowUpdateCmd = m.pane.Update(msg)
		return m, windowUpdateCmd

	case messages.ErrMsg:
		m.pane.SetContent(fmt.Sprintf("Error: %s", msg.Err))
	}

	if m.inFocus {
		m.pane, windowUpdateCmd = m.pane.Update(msg)
		cmds = append(cmds, windowUpdateCmd)
	}

	return m, tea.Batch(cmds...)
}
