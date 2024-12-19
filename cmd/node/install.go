package node

import (
	"github.com/algorandfoundation/algorun-tui/internal/algod"
	"github.com/algorandfoundation/algorun-tui/ui/style"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"os"
	"time"
)

const InstallMsg = "Installing Algorand"
const InstallExistsMsg = "algod is already installed"

var installCmd = &cobra.Command{
	Use:          "install",
	Short:        "Install the algorand daemon",
	Long:         style.Purple(style.BANNER) + "\n" + style.LightBlue("Install the algorand daemon on your local machine"),
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: yes flag

		// TODO: get expected version
		log.Info(style.Green.Render(InstallMsg + " vX.X.X"))
		// Warn user for prompt
		log.Warn(style.Yellow.Render(SudoWarningMsg))

		// TODO: compare expected version to existing version
		if algod.IsInstalled() && !force {
			log.Error(InstallExistsMsg)
			os.Exit(1)
		}

		// Run the installation
		err := algod.Install()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		time.Sleep(5 * time.Second)

		// If it's not running, start the daemon (can happen)
		if !algod.IsRunning() {
			err = algod.Start()
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
		}

		log.Info(style.Green.Render("Algorand installed successfully 🎉"))
	},
}

func init() {
	installCmd.Flags().BoolVarP(&force, "force", "f", false, style.Yellow.Render("forcefully install the node"))
}
