package tui

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/generative-ai-go/genai"
)

// Model represents the TUI's state.
type Model struct {
	textarea       textarea.Model
	viewport       viewport.Model
	renderer       glamour.TermRenderer
	response       string
	geminiResponse *genai.GenerateContentResponse
	err            error
	quitting       bool
	viewing        bool
}

// NewModel creates a new TUI model.
func NewModel() Model {

	ta := textarea.New()
	ta.Placeholder = "Enter your prompt here..."
	ta.Focus()

	width := 78
	vp := viewport.New(width, 30)

	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("228")).
		BorderBackground(lipgloss.Color("63")).
		BorderTop(true).
		BorderLeft(true).
		BorderBottom(true).
		BorderRight(true).
		PaddingTop(2).
		PaddingLeft(4).
		PaddingRight(4).
		PaddingBottom(4).
		Width(78)

	const glamourGutter = 2
	glamourRenderWidth := width - vp.Style.GetHorizontalFrameSize() - glamourGutter

	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(glamourRenderWidth),
	)

	return Model{
		textarea: ta,
		viewport: vp,
		renderer: *r,
	}
}
