package main

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"io"
	"strings"

	//"context"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/postsa/strut/internal/models"
	"github.com/postsa/strut/internal/tui"
	"log"
	"os"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type Picker struct {
	choice   string
	list     list.Model
	quitting bool
}

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}
func NewPicker() Picker {
	items := []list.Item{item("gemini-2.0-flash"), item("gpt-4o")}
	l := list.New(items, itemDelegate{}, 20, 6)
	l.Title = "Choose a model"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	return Picker{
		quitting: false,
		list:     l,
	}
}

func (p Picker) Init() tea.Cmd {
	return nil
}

func (p Picker) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return p, tea.Quit
		case tea.KeyEnter:
			i, ok := p.list.SelectedItem().(item)
			if ok {
				p.choice = string(i)
			}
			return p, tea.Quit
		}
	case tea.WindowSizeMsg:
		p.list.SetWidth(msg.Width)
		return p, nil
	case tea.QuitMsg:
		return p, nil
	}
	var cmd tea.Cmd
	p.list, cmd = p.list.Update(msg)
	return p, cmd
}

func (p Picker) View() string {
	if p.choice != "" {
		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", p.choice))
	}
	if p.quitting {
		return quitTextStyle.Render("Exiting")
	}
	return "\n" + p.list.View()
}

func main() {
	err, client, model := chooseModel()
	if err != nil {
		log.Fatal(err)
	}
	if client != nil {
		defer client.Close()
		next := ""
		for next != "quit" {
			program := tea.NewProgram(tui.NewModel(client, model), tea.WithAltScreen())
			output, err := program.Run()
			next = output.(tui.Model).Next
			if err != nil {
				log.Fatalf("There's been an error:", err)
			}
			if next == "choose" {
				err, client, model = chooseModel()
			}
		}
	} else {
		os.Exit(0)
	}
}

func chooseModel() (error, models.ChatClient, string) {
	var err error
	var client models.ChatClient

	pickerProgram := tea.NewProgram(NewPicker())
	output, err := pickerProgram.Run()
	picker := output.(Picker)
	model := picker.choice

	if model == "gemini-2.0-flash" {
		client, err = models.NewGemini(context.Background())
	} else if model == "gpt-4o" {
		client, err = models.NewOpenAi()
	}
	return err, client, model
}
