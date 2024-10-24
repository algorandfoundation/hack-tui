package keys

import (
	"sort"

	"github.com/algorandfoundation/hack-tui/api"
	"github.com/algorandfoundation/hack-tui/ui/controls"
	"github.com/algorandfoundation/hack-tui/ui/pages"
	"github.com/algorandfoundation/hack-tui/ui/utils"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

type ViewModel struct {
	Address string
	Data    *[]api.ParticipationKey
	Width   int
	Height  int

	table    table.Model
	controls controls.Model
}

func New(address string, keys *[]api.ParticipationKey) ViewModel {
	m := ViewModel{
		Address: address,
		Data:    keys,
		Width:   80,
		Height:  24,

		controls: controls.New(" (g)enerate | (a)ccounts | " + green.Render("(k)eys") + " | (t)xn | (d)elete "),

		table: table.New(),
	}
	m.table = table.New(
		table.WithColumns(m.makeColumns(80)),
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

	return m
}

func (m ViewModel) SelectedKey() *api.ParticipationKey {
	if m.Data == nil {
		return nil
	}
	var partkey *api.ParticipationKey
	for _, key := range *m.Data {
		selected := m.table.SelectedRow()
		if len(selected) > 0 && key.Id == selected[0] {
			partkey = &key
		}
	}
	return partkey
}
func (m ViewModel) makeColumns(width int) []table.Column {
	// TODO: refine responsiveness
	avgWidth := (width - lipgloss.Width(pages.Padding1("")) - 14) / 13
	return []table.Column{
		{Title: "ID", Width: avgWidth},
		{Title: "Address", Width: avgWidth},
		{Title: "SelectionParticipationKey", Width: avgWidth},
		{Title: "VoteParticipationKey", Width: avgWidth},
		{Title: "StateProofKey", Width: avgWidth},
		{Title: "VoteFirstValid", Width: avgWidth},
		{Title: "VoteLastValid", Width: avgWidth},
		{Title: "VoteKeyDilution", Width: avgWidth},
		{Title: "EffectiveLastValid", Width: avgWidth},
		{Title: "EffectiveFirstValid", Width: avgWidth},
		{Title: "LastVote", Width: avgWidth},
		{Title: "LastBlockProposal", Width: avgWidth},
		{Title: "LastStateProof", Width: avgWidth},
	}
}

func (m ViewModel) makeRows(keys *[]api.ParticipationKey) []table.Row {
	rows := make([]table.Row, 0)
	if keys == nil {
		return rows
	}
	for _, key := range *keys {
		if key.Address == m.Address {
			rows = append(rows, table.Row{
				key.Id,
				key.Address,
				*utils.UrlEncodeBytesPtrOrNil(key.Key.SelectionParticipationKey[:]),
				*utils.UrlEncodeBytesPtrOrNil(key.Key.VoteParticipationKey[:]),
				*utils.UrlEncodeBytesPtrOrNil(*key.Key.StateProofKey),
				utils.IntToStr(key.Key.VoteFirstValid),
				utils.IntToStr(key.Key.VoteLastValid),
				utils.IntToStr(key.Key.VoteKeyDilution),
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