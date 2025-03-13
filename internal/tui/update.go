package tui

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/google/generative-ai-go/genai"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/postsa/strut-cli/internal/gemini"
)

// Update handles TUI events.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			m.quitting = true
			return m, tea.Quit
		case tea.KeyEsc:
			m.viewing = false
			return m, tea.WindowSize()
		case tea.KeyEnter:
			prompt := m.textarea.Value()
			m.textarea.Reset()

			return m, func() tea.Msg {
				geminiClient, err := gemini.NewClient(context.Background())
				if err != nil {
					log.Printf("Error creating Gemini client: %v", err)
					return errMsg{err}
				}
				defer geminiClient.Close()

				resp, err := geminiClient.GenerateContent(context.Background(), prompt)
				if err != nil {
					log.Printf("Error generating content: %v", err)
					return errMsg{err}
				}
				return geminiResponseMsg{resp}
			}
		}

	case geminiResponseMsg:
		m.viewing = true
		m.geminiResponse = msg.response
		m.response = fmt.Sprintf("%v", msg.response.Candidates[0].Content.Parts[0])
		output, _ := m.renderer.Render(m.response)
		m.viewport.SetContent(output)

	case errMsg:
		m.err = msg.err
		m.viewport.SetContent(fmt.Sprintf("Error: %s", msg.err))
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

type geminiResponseMsg struct {
	response *genai.GenerateContentResponse
}

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}
