package main

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/postsa/strut/internal/models"
	"github.com/postsa/strut/internal/tui"
	"log"
	"os"
)

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
				log.Fatal("There's been an error:", err)
			}
			if next == "choose" {
				err, client, model = chooseModel()
				if model == "" {
					next = "quit"
				}
			}
		}
	} else {
		os.Exit(0)
	}
}

func chooseModel() (error, models.ChatClient, string) {
	var err error
	var client models.ChatClient

	pickerProgram := tea.NewProgram(tui.NewPicker())
	output, err := pickerProgram.Run()
	picker := output.(tui.Picker)
	model := picker.Choice

	if model == "gemini-2.0-flash" {
		client, err = models.NewGemini(context.Background())
	} else if model == "gpt-4o" {
		client, err = models.NewOpenAi()
	} else if model == "claude-opus-4-0" {
		client, err = models.NewClaude()
	}
	return err, client, model
}
