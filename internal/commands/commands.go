package commands

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/postsa/strut/internal/messages"
	"github.com/postsa/strut/internal/models"
	"log"
	"os"
	"os/exec"
	"time"
)

func ViewPortResizeCmd(width int) tea.Cmd {
	return func() tea.Msg {
		return messages.ViewPortResizeMessage{Width: width}
	}
}

func HistoryResizedCmd(newWidth int, totalWidth int) tea.Cmd {
	return func() tea.Msg {
		return messages.HistoryResizedMessage{NewWidth: newWidth, TotalWidth: totalWidth}
	}
}

func NewAnswerCmd(answer string, prompt string) tea.Cmd {
	return func() tea.Msg {
		return messages.NewAnswerMessage{Answer: answer, Prompt: prompt}
	}
}

func NewRenderCmd(content string) tea.Cmd {
	return func() tea.Msg {
		return messages.NewRenderMessage{Content: content}
	}
}

func SetAnswerCmd(answer string, answerRendered string) tea.Cmd {
	return func() tea.Msg {
		return messages.SetAnswerMessage{Answer: answer, AnswerRendered: answerRendered}
	}
}

func ExecuteVim(content string) tea.Cmd {
	file, err := os.CreateTemp("", "editor_*.md")
	if err != nil {
		log.Fatal(err)
	}
	_, err = file.WriteString(content)
	cmd := exec.Command("vim", file.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return messages.EditorFinishedMsg{Err: err, File: file}
	})
}

func TickCmd() tea.Cmd {
	return tea.Tick(time.Second/4, func(t time.Time) tea.Msg {
		return messages.TickMsg(t)
	})
}

func FetchResponseCmd(client models.ChatClient, prompt string) tea.Cmd {
	return func() tea.Msg {
		resp, err := client.GenerateContent(context.Background(), prompt)
		if err != nil {
			log.Printf("Error generating content: %v", err)
			return messages.ErrMsg{Err: err}
		}
		return messages.ResponseMsg{Response: resp, Prompt: prompt}
	}
}

func ExecutePromptCmd(prompt string) tea.Cmd {
	return func() tea.Msg {
		return messages.ExecutePromptMessage{Prompt: prompt}
	}
}
