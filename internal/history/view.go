package history

func (m HistoryModel) View() string {
	return m.listModel.View()
}
