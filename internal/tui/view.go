package tui

// View renders the TUI.
func (m Model) View() string {
	if m.quitting {
		return "Exiting...\n"
	}
	if m.viewing {
		return m.textarea.View() + "\n" + m.viewport.View()
	} else {
		return m.textarea.View()
	}
}
