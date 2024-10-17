package generate

import (
	"github.com/algorandfoundation/hack-tui/api"
	"github.com/algorandfoundation/hack-tui/internal"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

func (m ViewModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m ViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.HandleMessage(msg)
}

func (m ViewModel) HandleMessage(msg tea.Msg) (ViewModel, tea.Cmd) {

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.keyTable.SelectedRow() != nil {
				return m, tea.Batch(
					//m.viewport.SetContent()
					tea.Printf("Let's go to %s!", m.keyTable.SelectedRow()[1]),
				)
			}
			params := api.GenerateParticipationKeysParams{
				Dilution: nil,
				First:    0,
				Last:     1000,
			}
			_, err := internal.GenerateKeyPair(m.ctx, m.client, "WAFPLTCSVMCESEIMYPJHRADDGGKLB4LW4PFYCIU6VDCW3GLCJJS6RRWU3E", &params)
			if err != nil {
				log.Fatal(err)
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case error:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}
