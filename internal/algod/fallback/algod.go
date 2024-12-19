package fallback

import (
	"errors"
	"fmt"
	"github.com/algorandfoundation/algorun-tui/internal/algod/msgs"
	"github.com/algorandfoundation/algorun-tui/internal/algod/utils"
	"github.com/algorandfoundation/algorun-tui/internal/system"
	"github.com/charmbracelet/log"
	"os"
	"os/exec"
	"syscall"
)

// Install executes a series of commands to set up the Algorand node and development tools on a Unix environment.
// TODO: Allow for changing of the paths
func Install() error {
	return system.RunAll(system.CmdsList{
		{"mkdir", "~/node"},
		{"sh", "-c", "cd ~/node"},
		{"wget", "https://raw.githubusercontent.com/algorand/go-algorand/rel/stable/cmd/updater/update.sh"},
		{"chmod", "744", "update.sh"},
		{"sh", "-c", "./update.sh -i -c stable -p ~/node -d ~/node/data -n"},
	})

}

func Start() error {
	path, err := exec.LookPath("algod")
	log.Debug("Starting algod", "path", path)

	// Check if ALGORAND_DATA environment variable is set
	log.Info("Checking if ALGORAND_DATA env var is set...")
	algorandData := os.Getenv("ALGORAND_DATA")

	if !utils.IsDataDir(algorandData) {
		return errors.New(msgs.InvalidDataDirectory)
	}

	log.Info("ALGORAND_DATA env var set to valid directory: " + algorandData)

	cmd := exec.Command("algod")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("Failed to start algod: %v\n", err)
	}
	return nil
}

func Stop() error {
	log.Debug("Manually shutting down algod")
	// Find the process ID of algod
	pid, err := findAlgodPID()
	if err != nil {
		return err
	}

	// Send SIGTERM to the process
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		return err
	}

	return nil
}

func findAlgodPID() (int, error) {
	log.Debug("Scanning for algod process")
	cmd := exec.Command("pgrep", "algod")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	var pid int
	_, err = fmt.Sscanf(string(output), "%d", &pid)
	if err != nil {
		return 0, err
	}

	return pid, nil
}
