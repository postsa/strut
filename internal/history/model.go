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

type Model struct {
	inFocus         bool
	prompts         []list.Item
	listModel       list.Model
	answers         []string
	answersRendered []string
	width           int
	height          int
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Width() int {
	return m.listModel.Width()
}

func (m Model) Height() int {
	return m.listModel.Height()
}

func (m Model) Focus() Model {
	m.inFocus = true
	return m
}

func (m Model) Blur() Model {
	m.inFocus = false
	return m
}

func NewModel() Model {
	var l []list.Item
	lm := list.New(l, changeAnswerDelegate(), 20, 20)
	lm.Title = "History"
	lm.DisableQuitKeybindings()
	lm.Styles.TitleBar = lm.Styles.TitleBar.PaddingTop(1).AlignHorizontal(lipgloss.Center)

	return Model{
		inFocus:   false,
		prompts:   l,
		listModel: lm,
	}
}
