package tui

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

func (m ViewerModel) Focus() ViewerModel {
	m.inFocus = true
	m.pane.Style = m.pane.Style.BorderForeground(lipgloss.Color("228"))
	return m
}

func (m ViewerModel) Blur() ViewerModel {
	m.inFocus = false
	m.pane.Style = m.pane.Style.BorderForeground(lipgloss.Color("238"))
	return m
}

type ViewerModel struct {
	pane                   viewport.Model
	renderer               *glamour.TermRenderer
	currentContentRendered string
	currentContent         string
	inFocus                bool
}

func NewViewerModel() ViewerModel {
	vp := viewport.New(50, 30)
	viewportStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("228")).
		BorderTop(true).
		BorderLeft(true).
		BorderBottom(true).
		BorderRight(true).
		PaddingTop(2).
		PaddingLeft(2).
		PaddingRight(2).
		PaddingBottom(4)
	vp.Style = viewportStyle
	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	vp.KeyMap = viewport.KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "down"),
		),
	}
	return ViewerModel{
		pane:     vp,
		renderer: r,
		inFocus:  true,
	}
}

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
			return m, executeVim(m.currentContent)
		}
	case NewAnswerMessage:
		m.currentContent = msg.answer
		output, _ := m.renderer.Render(msg.answer)
		cmds = append(cmds, NewRenderCmd(output))

	case NewRenderMessage:
		m.currentContentRendered = msg.content
		m.pane.SetContent(m.currentContentRendered)
		m.pane.GotoTop()

	case SetAnswerMessage:
		m.currentContentRendered = msg.answerRendered
		m.currentContent = msg.answer
		m.pane.SetContent(m.currentContentRendered)
		m.pane.GotoTop()

	case tea.WindowSizeMsg:
		m.pane.Height = msg.Height - 9
		m.pane, windowUpdateCmd = m.pane.Update(msg)
		return m, windowUpdateCmd

	case ViewPortResizeMessage:
		m.pane.Style.MaxWidth(msg.width)
		m.pane.Width = msg.width
		m.pane, windowUpdateCmd = m.pane.Update(msg)
		return m, windowUpdateCmd

	case errMsg:
		m.pane.SetContent(fmt.Sprintf("Error: %s", msg.err))
	}

	if m.inFocus {
		m.pane, windowUpdateCmd = m.pane.Update(msg)
		cmds = append(cmds, windowUpdateCmd)
	}

	return m, tea.Batch(cmds...)
}

func (m ViewerModel) View() string {
	return m.pane.View()
}
