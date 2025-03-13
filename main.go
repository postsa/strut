package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/postsa/strut-cli/internal/tui"
	"os"
)

func main() {
	var dump *os.File
	var err error
	dump, err = os.OpenFile("messages.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		os.Exit(1)
	}
	program := tea.NewProgram(tui.NewModel(dump), tea.WithAltScreen())
	_, err = program.Run()
	if err != nil {
		fmt.Println("Bummer, there's been an error:", err)
		os.Exit(1)
	}
}
