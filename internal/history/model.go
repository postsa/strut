package history

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

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
