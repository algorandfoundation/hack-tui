package catchup

import (
	"context"
	"github.com/algorandfoundation/algorun-tui/api"
	"github.com/algorandfoundation/algorun-tui/cmd/utils"
	"github.com/algorandfoundation/algorun-tui/internal/algod"
	"github.com/algorandfoundation/algorun-tui/ui/style"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var startCmdLong = lipgloss.JoinVertical(
	lipgloss.Left,
	style.Purple(style.BANNER),
	"",
	style.Bold("Catchup the node to the latest catchpoint."),
	"",
	style.BoldUnderline("Overview:"),
	"Starting a catchup will sync the node to the latest catchpoint.",
	"Actual sync times may vary depending on the number of accounts, number of blocks and the network.",
	"",
	style.Yellow.Render("Note: Not all networks support Fast-Catchup."),
)

// startCmd is a Cobra command used to check the node's sync status and initiate a fast catchup when necessary.
var startCmd = utils.WithAlgodFlags(&cobra.Command{
	Use:          "start",
	Short:        "Get the latest catchpoint and start catching up.",
	Long:         startCmdLong,
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		err := utils.InitConfig()
		if err != nil {
			log.Fatal(err)
		}

		endpoint := viper.GetString("algod-endpoint")
		token := viper.GetString("algod-token")
		if endpoint == "" {
			log.Fatal("algod-endpoint is required")
		}
		if token == "" {
			log.Fatal("algod-token is required")
		}

		ctx := context.Background()
		httpPkg := new(api.HttpPkg)
		client, err := algod.GetClient(endpoint, token)
		cobra.CheckErr(err)

		status, response, err := algod.NewStatus(ctx, client, httpPkg)
		utils.WithInvalidResponsesExplanations(err, response, cmd.UsageString())

		if status.State == algod.FastCatchupState {
			log.Fatal(style.Red.Render("Node is currently catching up. Use --abort to cancel."))
		}

		// Get the latest catchpoint
		catchpoint, _, err := algod.GetLatestCatchpoint(httpPkg, status.Network)
		if err != nil && err.Error() == api.InvalidNetworkParamMsg {
			log.Fatal("This network does not support fast-catchup.")
		} else {
			log.Info(style.Green.Render("Latest Catchpoint: " + catchpoint))
		}

		// Start catchup
		res, _, err := algod.StartCatchup(ctx, client, catchpoint, nil)
		if err != nil {
			log.Fatal(err)
		}

		log.Info(style.Green.Render(res))
	},
}, &endpoint, &token)

func init() {
	startCmd.Flags().BoolVarP(&force, "force", "f", false, style.Yellow.Render("forcefully catchup the node"))
}
