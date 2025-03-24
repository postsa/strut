package tui

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/google/generative-ai-go/genai"
	"github.com/postsa/strut/internal/history"
	"github.com/postsa/strut/internal/input"
	"github.com/postsa/strut/internal/models"
	"github.com/postsa/strut/internal/viewer"
)

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
	client         *models.ChatClient
	Next           string
}

func NewModel(client models.ChatClient, model string) Model {

	i := input.NewModel(model)
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
		client:       &client,
		Next:         "quit",
	}
}
