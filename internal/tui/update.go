package tui

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/google/generative-ai-go/genai"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/postsa/strut-cli/internal/gemini"
)

type tickMsg time.Time

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd   tea.Cmd
		vpCmd   tea.Cmd
		listCmd tea.Cmd
		//prgsCmd tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.QuitMsg:
		return m, nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			if m.viewing {
				m.listFocus = true
				m.viewing = false
			} else {
				m.viewing = true
				m.listFocus = false
			}
		case tea.KeyCtrlC:
			m.quitting = true
			return m, tea.Quit
		case tea.KeyEsc:
			return m, nil
		case tea.KeyEnter:
			prompt := m.textinput.Value()
			m.loading = true
			m.textinput.Reset()
			newList := append(m.previousQuestionsList, item{title: prompt, desc: "some description"})
			m.previousQuestionsList = newList
			m.previousQuestionsListModel.SetItems(m.previousQuestionsList)
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
	case tea.WindowSizeMsg:
		m.textinput.Blur()
		m.progress.Width = msg.Width - 6
		m.previousQuestionsListModel.SetWidth(msg.Width / 3)
		m.previousQuestionsListModel.SetHeight(msg.Height - 4)
		m.resultsViewport.Height = msg.Height - 4
		m.resultsViewport.Style.MaxWidth(msg.Width - m.previousQuestionsListModel.Width())
		m.resultsViewport.Width = msg.Width - m.previousQuestionsListModel.Width()
		m.textinput.Focus()

	case geminiResponseMsg:
		m.geminiResponse = msg.response
		m.response = fmt.Sprintf("%v", msg.response.Candidates[0].Content.Parts[0])
		output, _ := m.mdRenderer.Render(m.response)
		m.resultsViewport.SetContent(output)
		m.loading = false

	case errMsg:
		m.err = msg.err
		m.resultsViewport.SetContent(fmt.Sprintf("Error: %s", msg.err))

	case tickMsg:
		var cmd tea.Cmd
		if !m.loading {
			cmd = m.progress.SetPercent(0)
		}

		cmd = m.progress.IncrPercent((1 - m.progress.Percent()) / 2 * (1 - m.progress.Percent()) / 2)
		return m, tea.Batch(tickCmd(), cmd)

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}

	if m.listFocus {
		m.previousQuestionsListModel, listCmd = m.previousQuestionsListModel.Update(msg)
	}
	m.textinput, tiCmd = m.textinput.Update(msg)
	if m.viewing {
		m.resultsViewport, vpCmd = m.resultsViewport.Update(msg)
	}
	cmds := []tea.Cmd{tiCmd, vpCmd, listCmd}

	if m.loading {
		cmds = append(cmds, tickCmd())
	}

	return m, tea.Batch(cmds...)
}

type geminiResponseMsg struct {
	response *genai.GenerateContentResponse
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}
