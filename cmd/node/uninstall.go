package node

import (
	"github.com/algorandfoundation/algorun-tui/internal/algod"
	"github.com/algorandfoundation/algorun-tui/ui/style"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

// UninstallWarningMsg provides a warning message to inform users they may be prompted for their password during uninstallation.
const UninstallWarningMsg = "(You may be prompted for your password to uninstall)"

// uninstallCmd defines a Cobra command used to uninstall the Algorand node (Algod) and related binaries from the system.
var uninstallCmd = &cobra.Command{
	Use:              "uninstall",
	Short:            "Uninstall Algorand node (Algod)",
	Long:             "Uninstall Algorand node (Algod) and other binaries on your system installed by this tool.",
	SilenceUsage:     true,
	PersistentPreRun: NeedsToBeStopped,
	Run: func(cmd *cobra.Command, args []string) {
		if force {
			log.Warn(style.Red.Render("Uninstalling Algorand (forcefully)"))
		}
		// Warn user for prompt
		log.Warn(style.Yellow.Render(UninstallWarningMsg))

		err := algod.Uninstall(force)
		if err != nil {
			log.Fatal(err)
		}
	},
}

// init initializes the uninstall command's flags, including the "force" flag for forced uninstallation.
func init() {
	uninstallCmd.Flags().BoolVarP(&force, "force", "f", false, style.Yellow.Render("forcefully uninstall the node"))
}
