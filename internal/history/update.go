package history

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/postsa/strut/internal/commands"
	"github.com/postsa/strut/internal/messages"
	"time"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case ChangeAnswerMessage:
		if len(m.prompts) > 0 {
			return m, commands.SetAnswerCmd(m.answers[msg.position], m.answersRendered[msg.position])
		}

	case messages.NewAnswerMessage:
		m.prompts = append(m.prompts, item{title: msg.Prompt, desc: time.Now().Format("01/02/06 03:04 PM")})
		m.listModel.SetItems(m.prompts)
		m.listModel.Select(len(m.prompts) - 1)
		m.answers = append(m.answers, msg.Answer)
		m.listModel, cmd = m.listModel.Update(msg)

	case messages.NewRenderMessage:
		m.answersRendered = append(m.answersRendered, msg.Content)

	case tea.WindowSizeMsg:
		m.listModel.SetWidth(msg.Width / 3)
		m.listModel.SetHeight(msg.Height - 9)
		resizeCmd := commands.HistoryResizedCmd(m.listModel.Width(), msg.Width)
		m.listModel, cmd = m.listModel.Update(msg)
		return m, tea.Batch(resizeCmd, cmd)
	}

	if m.inFocus {
		m.listModel, cmd = m.listModel.Update(msg)
	}
	return m, cmd
}
