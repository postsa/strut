package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/generative-ai-go/genai"
	"io"
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

// Model represents the TUI's state.
type Model struct {
	textinput                  textinput.Model
	resultsViewport            viewport.Model
	previousQuestionsListModel list.Model
	mdRenderer                 glamour.TermRenderer
	response                   string
	geminiResponse             *genai.GenerateContentResponse
	err                        error
	quitting                   bool
	viewing                    bool
	previousQuestionsList      []list.Item
	dump                       io.Writer
	loading                    bool
	progress                   progress.Model
	listFocus                  bool
	previousAnswersRendered    []string
	currentContentRendered     string
	previousAnswers            []string
	currentContent             string
}

// NewModel creates a new TUI model.
func NewModel(dump io.Writer) Model {

	ti := textinput.New()
	ti.Placeholder = "Enter your prompt here..."

	ti.Focus()

	width := 50
	rvp := viewport.New(width, 30)

	var pql []list.Item
	pqlm := list.New(pql, getAnswerDelegate(), 20, 20)

	pqlm.Title = "History"

	pqlm.DisableQuitKeybindings()

	var pa []string

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

	rvp.Style = viewportStyle

	pb := progress.New(progress.WithDefaultGradient())

	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)

	rvp.KeyMap = viewport.KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "down"),
		),
	}

	return Model{
		textinput:                  ti,
		resultsViewport:            rvp,
		mdRenderer:                 *r,
		previousQuestionsList:      pql,
		previousQuestionsListModel: pqlm,
		dump:                       dump,
		viewing:                    true,
		loading:                    false,
		progress:                   pb,
		listFocus:                  false,
		previousAnswersRendered:    pa,
	}
}
