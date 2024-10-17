package transaction

import (
	"github.com/charmbracelet/lipgloss"
)

func (m ViewModel) View() string {
	qrCode := "TODO"

	return lipgloss.JoinVertical(lipgloss.Center, qrCode, m.controls.View())
}
