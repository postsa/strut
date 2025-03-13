package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/postsa/strut-cli/internal/tui"
	"os"
)

func main() {
	program := tea.NewProgram(tui.NewModel(), tea.WithAltScreen())
	_, err := program.Run()
	if err != nil {
		fmt.Println("Bummer, there's been an error:", err)
		os.Exit(1)
	}
}
