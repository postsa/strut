package tui

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/generative-ai-go/genai"
	"log"
	"os"
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/postsa/strut-cli/internal/gemini"
)

type tickMsg time.Time

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd      tea.Cmd
		vpCmd      tea.Cmd
		historyCmd tea.Cmd
		viewerCmd  tea.Cmd
	)
	var (
		historyModel tea.Model
		viewerModel  tea.Model
	)

	var cmds []tea.Cmd

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
			if m.textinput.Focused() {
				prompt := m.textinput.Value()
				currentProgressWidth := m.progress.Width
				m.progress = progress.New(progress.WithDefaultGradient())
				m.progress.Width = currentProgressWidth
				m.progress.SetPercent(0)
				cmd := m.progress.IncrPercent(.1)
				m.loading = true
				m.textinput.Reset()
				return m, tea.Batch(cmd, tickCmd(), fetchResponseCmd(m.client, prompt))
			}
		}
	case tea.WindowSizeMsg:
		m.textinput.Blur()
		m.progress.Width = msg.Width - 6
		m.textinput.Focus()

	case HistoryResizedMessage:
		return m, ViewPortResizeCmd(msg.totalWidth - msg.newWidth)

	case geminiResponseMsg:
		m.geminiResponse = msg.response
		m.response = fmt.Sprintf("%v", msg.response.Candidates[0].Content.Parts[0])
		cmds = append(cmds, NewAnswerCmd(m.response, msg.prompt))

	case NewRenderMessage:
		m.loading = false

	case SetAnswerMessage:
		m.viewing = true
		m.listFocus = false

	case errMsg:
		m.err = msg.err

	case tickMsg:
		cmd := m.progress.IncrPercent(((1 - m.progress.Percent()) / 3) * ((1 - m.progress.Percent()) / 1.2))
		return m, tea.Batch(tickCmd(), cmd)

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	case editorFinishedMsg:
		defer os.Remove(msg.file.Name())
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}
	}

	if m.listFocus {
		m.textinput.Blur()
		m.historyModel = m.historyModel.Focus()
		m.viewerModel = m.viewerModel.Blur()
		m.textinput.TextStyle = lipgloss.NewStyle().Background(lipgloss.Color("238"))
		m.textinput.PromptStyle = lipgloss.NewStyle().Background(lipgloss.Color("238"))
		m.textinput.PlaceholderStyle = lipgloss.NewStyle().Background(lipgloss.Color("238")).Foreground(lipgloss.Color("238"))
	}
	if m.viewing {
		m.historyModel = m.historyModel.Blur()
		m.viewerModel = m.viewerModel.Focus()
		m.textinput.TextStyle = lipgloss.NewStyle().Background(lipgloss.Color("89"))
		m.textinput.PromptStyle = lipgloss.NewStyle().Background(lipgloss.Color("89"))
		m.textinput.PlaceholderStyle = lipgloss.NewStyle().Background(lipgloss.Color("89")).Foreground(lipgloss.Color("228"))
		m.textinput.Focus()
	}

	m.textinput, tiCmd = m.textinput.Update(msg)

	historyModel, historyCmd = m.historyModel.Update(msg)
	m.historyModel = historyModel.(HistoryModel)

	viewerModel, viewerCmd = m.viewerModel.Update(msg)
	m.viewerModel = viewerModel.(ViewerModel)

	cmds = append(cmds, tiCmd, vpCmd, historyCmd, viewerCmd)
	return m, tea.Batch(cmds...)
}

type geminiResponseMsg struct {
	response *genai.GenerateContentResponse
	prompt   string
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second/4, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func fetchResponseCmd(client *gemini.Client, prompt string) tea.Cmd {
	return func() tea.Msg {
		resp, err := client.GenerateContent(context.Background(), prompt)
		if err != nil {
			log.Printf("Error generating content: %v", err)
			return errMsg{err}
		}
		return geminiResponseMsg{response: resp, prompt: prompt}
	}
}

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

type editorFinishedMsg struct {
	err  error
	file *os.File
}

func executeVim(content string) tea.Cmd {
	file, err := os.CreateTemp("", "editor_*.md")
	if err != nil {
		log.Fatal(err)
	}
	_, err = file.WriteString(content)
	cmd := exec.Command("vim", file.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return editorFinishedMsg{err, file}
	})
}

type SetAnswerMessage struct {
	answer         string
	answerRendered string
}
type NewAnswerMessage struct {
	answer string
	prompt string
}
type NewRenderMessage struct{ content string }

func NewAnswerCmd(answer string, prompt string) tea.Cmd {
	return func() tea.Msg {
		return NewAnswerMessage{answer, prompt}
	}
}

func NewRenderCmd(content string) tea.Cmd {
	return func() tea.Msg {
		return NewRenderMessage{content}
	}
}

func SetAnswerCmd(answer string, answerRendered string) tea.Cmd {
	return func() tea.Msg {
		return SetAnswerMessage{answer, answerRendered}
	}
}

type HistoryResizedMessage struct {
	newWidth   int
	totalWidth int
}

func HistoryResizedCmd(newWidth int, totalWidth int) tea.Cmd {
	return func() tea.Msg {
		return HistoryResizedMessage{newWidth, totalWidth}
	}
}

type ViewPortResizeMessage struct {
	width int
}

func ViewPortResizeCmd(width int) tea.Cmd {
	return func() tea.Msg {
		return ViewPortResizeMessage{width}
	}
}
