package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/generative-ai-go/genai"
	"github.com/postsa/strut-cli/internal/gemini"
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

// Model represents the TUI's state.
type Model struct {
	textinput              textinput.Model
	resultsViewport        viewport.Model
	mdRenderer             glamour.TermRenderer
	historyModel           HistoryModel
	response               string
	geminiResponse         *genai.GenerateContentResponse
	err                    error
	quitting               bool
	viewing                bool
	loading                bool
	progress               progress.Model
	listFocus              bool
	currentContentRendered string
	currentContent         string
	modelName              string
	client                 *gemini.Client
}

// NewModel creates a new TUI model.
func NewModel(client *gemini.Client) Model {
	modelName := "gemini-2.0-flash"

	ti := textinput.New()
	ti.Prompt = "(" + modelName + ")" + " > "
	ti.Placeholder = "ask a question ..."
	ti.Focus()

	h := NewHistoryModel()

	width := 50
	rvp := viewport.New(width, 30)

	p := progress.New(progress.WithDefaultGradient())

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
		textinput:       ti,
		resultsViewport: rvp,
		mdRenderer:      *r,
		viewing:         true,
		loading:         false,
		listFocus:       false,
		progress:        p,
		modelName:       "gemini-2.0-flash",
		currentContent:  "",
		client:          client,
		historyModel:    h,
	}
}
