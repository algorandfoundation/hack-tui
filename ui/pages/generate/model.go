package generate

import (
	"github.com/algorandfoundation/hack-tui/api"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
)

type ViewModel struct {
	Address   string
	keyTable  table.Model
	textInput textinput.Model
	err       error
}

func New(address string, partkeys *[]api.ParticipationKey) ViewModel {
	ti := textinput.New()
	ti.Placeholder = "Pikachu"
	ti.Focus()
	ti.CharLimit = 56
	ti.Width = 56

	return ViewModel{
		textInput: ti,
		err:       nil,
	}
}
