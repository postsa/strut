package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"time"
)

type HistoryModel struct {
	inFocus         bool
	prompts         []list.Item
	listModel       list.Model
	answers         []string
	answersRendered []string
	width           int
	height          int
}

func (m HistoryModel) Init() tea.Cmd {
	return nil
}

func (m HistoryModel) Width() int {
	return m.listModel.Width()
}

func (m HistoryModel) Height() int {
	return m.listModel.Height()
}

func (m HistoryModel) Focus() HistoryModel {
	m.inFocus = true
	return m
}

func (m HistoryModel) Blur() HistoryModel {
	m.inFocus = false
	return m
}

func NewHistoryModel() HistoryModel {
	var l []list.Item
	lm := list.New(l, changeAnswerDelegate(), 20, 20)
	lm.Title = "History"
	lm.DisableQuitKeybindings()
	lm.Styles.TitleBar = lm.Styles.TitleBar.PaddingTop(1).AlignHorizontal(lipgloss.Center)

	return HistoryModel{
		inFocus:   false,
		prompts:   l,
		listModel: lm,
	}
}

func (m HistoryModel) View() string {
	return m.listModel.View()
}

func (m HistoryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case ChangeAnswerMessage:
		if len(m.prompts) > 0 {
			return m, SetAnswerCmd(m.answers[msg.position], m.answersRendered[msg.position])
		}

	case NewAnswerMessage:
		m.prompts = append(m.prompts, item{title: msg.prompt, desc: time.Now().Format("01/02/06 03:04 PM")})
		m.listModel.SetItems(m.prompts)
		m.listModel.Select(len(m.prompts) - 1)
		m.answers = append(m.answers, msg.answer)
		m.listModel, cmd = m.listModel.Update(msg)

	case NewRenderMessage:
		m.answersRendered = append(m.answersRendered, msg.content)

	case tea.WindowSizeMsg:
		m.listModel.SetWidth(msg.Width / 3)
		m.listModel.SetHeight(msg.Height - 9)
		resizeCmd := HistoryResizedCmd(m.listModel.Width(), msg.Width)
		m.listModel, cmd = m.listModel.Update(msg)
		return m, tea.Batch(resizeCmd, cmd)
	}

	if m.inFocus {
		m.listModel, cmd = m.listModel.Update(msg)
	}
	return m, cmd
}
