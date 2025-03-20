package main

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/postsa/strut/internal/gemini"
	"github.com/postsa/strut/internal/tui"
	"log"
	"os"
)

func main() {
	var err error
	client, err := gemini.NewClient(context.Background())
	if err != nil {
		log.Printf("Error creating Gemini client: %v", err)
		os.Exit(1)
	}
	defer client.Close()
	program := tea.NewProgram(tui.NewModel(client), tea.WithAltScreen())
	_, err = program.Run()
	if err != nil {
		log.Printf("Bummer, there's been an error:", err)
		os.Exit(1)
	}
}
