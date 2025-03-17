package tui

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/google/generative-ai-go/genai"
	"github.com/postsa/strut-cli/internal/gemini"
	"github.com/postsa/strut-cli/internal/history"
	"github.com/postsa/strut-cli/internal/input"
	"github.com/postsa/strut-cli/internal/viewer"
)

// Model represents the TUI's state.
type Model struct {
	inputModel     input.Model
	viewerModel    viewer.Model
	historyModel   history.Model
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

	i := input.NewModel("gemini-2.0-flash")
	h := history.NewModel()
	v := viewer.NewModel()

	p := progress.New(progress.WithDefaultGradient())

	return Model{
		inputModel:   i,
		historyModel: h,
		viewerModel:  v,
		viewing:      true,
		loading:      false,
		listFocus:    false,
		progress:     p,
		modelName:    "gemini-2.0-flash",
		client:       client,
	}
}
