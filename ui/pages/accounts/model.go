package accounts

import (
	"context"
	"github.com/algorandfoundation/hack-tui/api"
	"github.com/algorandfoundation/hack-tui/internal"
	"github.com/algorandfoundation/hack-tui/ui/controls"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"sort"
	"strconv"
	"time"
)

var green = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))

type ViewModel struct {
	Width      int
	Height     int
	ViewHeight int
	ViewWidth  int
	table      table.Model
	controls   controls.Model

	rowsChannel chan []table.Row

	ctx    context.Context
	client *api.ClientWithResponses
}

func (m ViewModel) SelectedAccount() string {
	row := m.table.SelectedRow()
	return row[0]
}

// TODO: remove this polling in favor of upstream purposed algod response
func (m ViewModel) pollForKeyChanges(interval time.Duration) error {
	// Sleep then try again
	time.Sleep(interval)
	// Fetch the latest keys
	currentKeys, err := internal.GetPartKeys(m.ctx, m.client)
	if err != nil {
		return err
	}
	if len(*currentKeys) != len(m.table.Rows()) {
		rows := *m.makeRows(currentKeys)
		m.rowsChannel <- rows
		m.table.SetRows(rows)
	}

	return m.pollForKeyChanges(interval)
}

func New(ctx context.Context, client *api.ClientWithResponses) (ViewModel, error) {
	m := ViewModel{
		Width:      80,
		Height:     24,
		ViewHeight: 24,
		ViewWidth:  80,

		controls: controls.New(green.Render(" (a)ccunts") + " | (k)eys | (t)xn "),

		table:       table.New(),
		ctx:         ctx,
		client:      client,
		rowsChannel: make(chan []table.Row),
	}
	keys, err := internal.GetPartKeys(m.ctx, m.client)

	if err != nil {
		return m, err
	}

	m.table = table.New(
		table.WithColumns(m.makeColumns()),
		table.WithRows(*m.makeRows(keys)),
		table.WithFocused(true),
		table.WithHeight(m.Height-lipgloss.Height(m.controls.View())-1),
		table.WithWidth(m.Width),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	m.table.SetStyles(s)

	// Watch for Key Changes
	//go func() {
	//	// TODO: get algod to update the generate endpoint
	//	err := m.pollForKeyChanges(1 * time.Second)
	//	if err != nil {
	//		//	panic(err)
	//	}
	//}()

	return m, nil
}
func (m ViewModel) makeColumns() []table.Column {
	// TODO: refine responsiveness
	fillSize := max(0, (m.Width-49)/2)
	return []table.Column{
		{Title: "ID", Width: 10},
		{Title: "Account", Width: fillSize},
		{Title: "Status", Width: hidden(20, fillSize)},
		{Title: "Expires", Width: 15},
		{Title: "Balance", Width: fillSize},
	}
}

func (m ViewModel) makeRows(keys *[]api.ParticipationKey) *[]table.Row {
	rows := make([]table.Row, 0)
	values := internal.AccountsFromParticipationKeys(keys)

	for key := range values {
		rows = append(rows, table.Row{
			values[key].Address,
			strconv.Itoa(values[key].Keys),
			values[key].Status,
			strconv.Itoa(values[key].Balance),
		})
	}
	sort.SliceStable(rows, func(i, j int) bool {
		return rows[i][0] < rows[j][0]
	})
	return &rows
}
