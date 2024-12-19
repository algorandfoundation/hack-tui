package app

import (
	"github.com/algorandfoundation/algorun-tui/internal/algod"
	tea "github.com/charmbracelet/bubbletea"
)

// AccountSelected is a type alias for algod.Account, representing a selected account during application runtime.
type AccountSelected algod.Account

// EmitAccountSelected waits for and retrieves a new set of table rows from a given channel.
func EmitAccountSelected(account algod.Account) tea.Cmd {
	return func() tea.Msg {
		return AccountSelected(account)
	}
}
