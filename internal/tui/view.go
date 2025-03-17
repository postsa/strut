package tui

import "github.com/charmbracelet/lipgloss"

// View renders the TUI.
func (m Model) View() string {
	if m.quitting {
		return "Exiting...\n"
	}
	header := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("228")).Render("\n  Strut LLM CLI\n")
	var bottom string
	bottom = lipgloss.JoinHorizontal(lipgloss.Top, m.resultsViewport.View(), " ", m.previousQuestionsListModel.View())
	if m.loading && &m.progress != nil {
		return lipgloss.JoinVertical(lipgloss.Left, header, bottom, "\n   "+m.progress.View())
	}
	return lipgloss.JoinVertical(lipgloss.Left, header, m.textinput.View(), bottom)
}
