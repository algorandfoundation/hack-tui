package generate

import (
	"context"
	"github.com/algorandfoundation/hack-tui/api"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
)

type ViewModel struct {
	Address   string
	keyTable  table.Model
	textInput textinput.Model
	err       error
	ctx       context.Context
	client    *api.ClientWithResponses
}

func New(ctx context.Context, client *api.ClientWithResponses) ViewModel {
	ti := textinput.New()
	ti.Placeholder = "Pikachu"
	ti.Focus()
	ti.CharLimit = 56
	ti.Width = 56

	return ViewModel{
		textInput: ti,
		err:       nil,
		ctx:       ctx,
		client:    client,
	}
}
