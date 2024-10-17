package transaction

import (
	"context"
	"github.com/algorandfoundation/hack-tui/api"
	"github.com/algorandfoundation/hack-tui/ui/controls"
	"github.com/charmbracelet/lipgloss"
)

var green = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))

type ViewModel struct {
	Width      int
	Height     int
	ViewHeight int
	ViewWidth  int

	controls controls.Model

	// TODO: add URL
	// urlTxn   string
	ctx    context.Context
	client *api.ClientWithResponses
}

func New(ctx context.Context, client *api.ClientWithResponses) ViewModel {
	return ViewModel{
		ctx:      ctx,
		client:   client,
		controls: controls.New("(a)ccounts | (k)eys | " + green.Render("(t)xn ")),
	}
}
