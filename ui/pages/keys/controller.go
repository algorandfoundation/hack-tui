package keys

import (
	"github.com/algorandfoundation/hack-tui/api"
	"github.com/algorandfoundation/hack-tui/internal"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

func (m ViewModel) Init() tea.Cmd {
	return nil
}

func (m ViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.HandleMessage(msg)
}

func (m ViewModel) HandleMessage(msg tea.Msg) (ViewModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tableRows:
		m.table.SetRows(msg)
	case tea.KeyMsg:
		switch msg.String() {
		case "t":
			// TODO: navigation
			//	m.viewportStatus = ViewportStatusTransaction
			//	return m, nil
		case "g":
			// TODO: navigation
			//m.viewportStatus = ViewportStatusGenerate
			params := api.GenerateParticipationKeysParams{
				Dilution: nil,
				First:    0,
				Last:     1000,
			}
			_, err := internal.GenerateKeyPair(m.ctx, m.client, "WAFPLTCSVMCESEIMYPJHRADDGGKLB4LW4PFYCIU6VDCW3GLCJJS6RRWU3E", &params)
			if err != nil {
				log.Fatal(err)
			}
		case "d":
			row := m.table.SelectedRow()
			if row == nil {
				return m, nil
			}
			err := internal.DeletePartKey(m.ctx, m.client, row[0])
			if err != nil {
				log.Fatal(err)
			}
			keys, err := internal.GetPartKeys(m.ctx, m.client)
			if err != nil {
				log.Fatal(err)
			}
			m.table.SetRows(m.makeRows(keys))
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		controlsHeight := lipgloss.Height(m.controls.View())
		m.table.SetWidth(m.ViewWidth - 3)
		m.table.SetHeight(max(0, m.ViewHeight-controlsHeight))
		m.table.SetColumns(m.makeColumns())
	}

	var cmds []tea.Cmd
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)
	m.controls, cmd = m.controls.HandleMessage(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
