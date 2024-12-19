package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/algorandfoundation/algorun-tui/internal"
	"github.com/algorandfoundation/algorun-tui/ui"
	"github.com/algorandfoundation/algorun-tui/ui/style"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// statusCmd is the main entrypoint for the `status` cobra.Command with a tea.Program
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the node status",
	Long:  style.Purple(style.BANNER) + "\n" + style.LightBlue("View the node status"),
	RunE: func(cmd *cobra.Command, args []string) error {
		initConfig()
		if viper.GetString("algod-endpoint") == "" {
			return errors.New(style.Magenta("algod-endpoint is required"))
		}

		// Get Algod from configuration
		client, err := getClient()
		cobra.CheckErr(err)
		state := internal.StateModel{
			Status: internal.StatusModel{
				State:       "SYNCING",
				Version:     "N/A",
				Network:     "N/A",
				Voting:      false,
				NeedsUpdate: true,
				LastRound:   0,
			},
			Metrics: internal.MetricsModel{
				RoundTime: 0,
				TPS:       0,
				RX:        0,
				TX:        0,
			},
			ParticipationKeys: nil,
		}
		err = state.Status.Fetch(context.Background(), client, new(internal.HttpPkg))
		cobra.CheckErr(err)
		// Create the TUI
		view := ui.MakeStatusViewModel(&state)

		p := tea.NewProgram(view, tea.WithAltScreen())
		go func() {
			state.Watch(func(status *internal.StateModel, err error) {
				cobra.CheckErr(err)
				p.Send(state)
			}, context.Background(), client)
		}()
		// Execute the Command
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
		return nil
	},
}
