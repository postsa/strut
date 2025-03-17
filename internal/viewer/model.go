package viewer

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
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
