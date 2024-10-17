package keys

import (
	"context"
	"github.com/algorandfoundation/hack-tui/api"
	"github.com/algorandfoundation/hack-tui/internal"
	"github.com/algorandfoundation/hack-tui/ui/controls"
	"github.com/algorandfoundation/hack-tui/ui/utils"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"sort"
)

var green = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))

type tableRows = []table.Row

type ViewModel struct {
	Address    string
	Width      int
	Height     int
	ViewHeight int
	ViewWidth  int
	table      table.Model
	controls   controls.Model

	rowsChannel chan tableRows

	ctx    context.Context
	client *api.ClientWithResponses
}

//func (m ViewModel) SelectedParticipationKey() api.ParticipationKey {
//	id := m.table.SelectedRow()[0]
//}

func New(ctx context.Context, client *api.ClientWithResponses) (ViewModel, error) {
	m := ViewModel{
		Address:    "WAFPLTCSVMCESEIMYPJHRADDGGKLB4LW4PFYCIU6VDCW3GLCJJS6RRWU3E",
		Width:      80,
		Height:     24,
		ViewHeight: 24,
		ViewWidth:  80,

		controls: controls.New("(a)ccounts | " + green.Render("(k)eys") + " | (t)xn | (d)elete | (g)enerate "),

		table:       table.New(),
		ctx:         ctx,
		client:      client,
		rowsChannel: make(chan tableRows),
	}
	keys, err := internal.GetPartKeys(m.ctx, m.client)

	if err != nil {
		return m, err
	}

	m.table = table.New(
		table.WithColumns(m.makeColumns()),
		table.WithRows(m.makeRows(keys)),
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

	return m, nil
}
func (m ViewModel) makeColumns() []table.Column {
	// TODO: refine responsiveness
	fillSize := max(0, (m.Width-49)/2)
	return []table.Column{
		{Title: "ID", Width: 10},
		{Title: "Address", Width: fillSize},
		{Title: "SelectionKey", Width: fillSize},
		{Title: "SelectionKey", Width: fillSize},
		{Title: "SelectionKey", Width: fillSize},
		{Title: "EffectiveLastValid", Width: hidden(20, fillSize)},
		{Title: "EffectiveFirstValid", Width: 15},
		{Title: "LastVote", Width: 10},
		{Title: "LastBlockProposal", Width: 10},
		{Title: "LastStateProof", Width: 10},
	}
}

func (m ViewModel) makeRows(keys *[]api.ParticipationKey) []table.Row {
	rows := make([]table.Row, 0)
	for _, key := range *keys {
		if key.Address == m.Address {
			rows = append(rows, table.Row{
				key.Id,
				key.Address,
				*utils.UrlEncodeBytesPtrOrNil(key.Key.SelectionParticipationKey[:]),
				*utils.UrlEncodeBytesPtrOrNil(key.Key.VoteParticipationKey[:]),
				*utils.UrlEncodeBytesPtrOrNil(*key.Key.StateProofKey),
				//utils.StrOrNA(key.Key.VoteKeyDilution),
				//utils.StrOrNA(key.Key.StateProofKey),
				utils.StrOrNA(key.EffectiveLastValid),
				utils.StrOrNA(key.EffectiveFirstValid),
				utils.StrOrNA(key.LastVote),
				utils.StrOrNA(key.LastBlockProposal),
				utils.StrOrNA(key.LastStateProof),
			})
		}
	}
	sort.SliceStable(rows, func(i, j int) bool {
		return rows[i][0] < rows[j][0]
	})
	return rows
}
