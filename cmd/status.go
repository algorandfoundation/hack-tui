package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/algorandfoundation/algorun-tui/api"
	"github.com/algorandfoundation/algorun-tui/cmd/utils"
	"github.com/algorandfoundation/algorun-tui/cmd/utils/explanations"
	"github.com/algorandfoundation/algorun-tui/internal"
	"github.com/algorandfoundation/algorun-tui/internal/algod"
	"github.com/algorandfoundation/algorun-tui/ui"
	"github.com/algorandfoundation/algorun-tui/ui/style"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// statusCmd is the main entrypoint for the `status` cobra.Command with a tea.Program
var statusCmd = utils.WithAlgodFlags(&cobra.Command{
	Use:   "status",
	Short: "Get the node status",
	Long:  style.Purple(style.BANNER) + "\n" + style.LightBlue("View the node status"),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := utils.InitConfig()
		if err != nil {
			return err
		}
		if viper.GetString("algod-endpoint") == "" {
			return errors.New(style.Magenta("algod-endpoint is required"))
		}
		ctx := context.Background()
		httpPkg := new(api.HttpPkg)
		// Get Algod from configuration
		client, err := algod.GetClient(viper.GetString("algod-endpoint"), viper.GetString("algod-token"))
		cobra.CheckErr(err)

		// Fetch the state and handle any creation errors
		state, stateResponse, err := internal.NewStateModel(ctx, client, httpPkg)
		if stateResponse.StatusCode() == 401 {
			return fmt.Errorf(
				style.Red.Render("failed to get status: Unauthorized") + explanations.TokenInvalid)
		}
		if stateResponse.StatusCode() > 300 {
			return fmt.Errorf(
				style.Red.Render("failed to get status: error code %d")+explanations.TokenNotAdmin,
				stateResponse.StatusCode())
		}
		cobra.CheckErr(err)

		// Create the TUI
		view := ui.MakeStatusViewModel(state)

		p := tea.NewProgram(view, tea.WithAltScreen())
		go func() {
			// Watch for State Changes
			state.Watch(func(status *internal.StateModel, err error) {
				if err != nil {
					state.Stop()
				}
				cobra.CheckErr(err)
				// Send the state to the TUI
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
}, &algodEndpoint, &algodToken)
