package app

import (
	"github.com/algorandfoundation/algorun-tui/internal/nodekit"
	tea "github.com/charmbracelet/bubbletea"
)

type AccountSelected nodekit.Account

// EmitAccountSelected waits for and retrieves a new set of table rows from a given channel.
func EmitAccountSelected(account nodekit.Account) tea.Cmd {
	return func() tea.Msg {
		return AccountSelected(account)
	}
}
