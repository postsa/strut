package tui

import "github.com/charmbracelet/lipgloss"

// View renders the TUI.
func (m Model) View() string {
	if m.quitting {
		return "Exiting...\n"
	}
	var bottom string
	bottom = lipgloss.JoinHorizontal(lipgloss.Left, m.resultsViewport.View(), m.previousQuestionsListModel.View())
	if m.loading {
		return lipgloss.JoinVertical(lipgloss.Top, m.textinput.View(), bottom, m.progress.View())
	}
	return lipgloss.JoinVertical(lipgloss.Top, m.textinput.View(), bottom)
}
