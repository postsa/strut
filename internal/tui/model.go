package tui

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
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
	textinput      textinput.Model
	viewerModel    ViewerModel
	historyModel   HistoryModel
	response       string
	geminiResponse *genai.GenerateContentResponse
	err            error
	quitting       bool
	viewing        bool
	loading        bool
	progress       progress.Model
	listFocus      bool
	modelName      string
	client         *gemini.Client
}

// NewModel creates a new TUI model.
func NewModel(client *gemini.Client) Model {
	modelName := "gemini-2.0-flash"

	ti := textinput.New()
	ti.Prompt = "(" + modelName + ")" + " > "
	ti.Placeholder = "ask a question ..."
	ti.Focus()

	h := NewHistoryModel()
	v := NewViewerModel()

	p := progress.New(progress.WithDefaultGradient())

	return Model{
		textinput:    ti,
		viewing:      true,
		loading:      false,
		listFocus:    false,
		progress:     p,
		modelName:    "gemini-2.0-flash",
		client:       client,
		historyModel: h,
		viewerModel:  v,
	}
}
