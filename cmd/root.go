package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/algorandfoundation/algorun-tui/api"
	"github.com/algorandfoundation/algorun-tui/cmd/configure"
	"github.com/algorandfoundation/algorun-tui/cmd/node"
	"github.com/algorandfoundation/algorun-tui/internal"
	"github.com/algorandfoundation/algorun-tui/ui"
	"github.com/algorandfoundation/algorun-tui/ui/explanations"
	"github.com/algorandfoundation/algorun-tui/ui/style"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/securityprovider"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"os"
	"runtime"
	"strings"
)

var (
	algod   string
	token   = strings.Repeat("a", 64)
	Version = ""
	rootCmd = &cobra.Command{
		Use:   "algorun",
		Short: "Manage Algorand nodes",
		Long:  style.Purple(style.BANNER) + "\n",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			log.SetOutput(cmd.OutOrStdout())
			initConfig()

			if viper.GetString("algod-endpoint") == "" {
				return fmt.Errorf(style.Red.Render("algod-endpoint is required") + explanations.NodeNotFound)
			}

			if viper.GetString("algod-token") == "" {
				return fmt.Errorf(style.Red.Render("algod-token is required"))
			}

			client, err := getClient()
			cobra.CheckErr(err)

			ctx := context.Background()
			v, err := client.GetStatusWithResponse(ctx)
			if err != nil {
				return fmt.Errorf(
					style.Red.Render("failed to get status: %s")+explanations.Unreachable,
					err)
			} else if v.StatusCode() == 401 {
				return fmt.Errorf(
					style.Red.Render("failed to get status: Unauthorized") + explanations.TokenInvalid)
			} else if v.StatusCode() != 200 {
				return fmt.Errorf(
					style.Red.Render("failed to get status: error code %d")+explanations.TokenNotAdmin,
					v.StatusCode())
			}

			partkeys, err := internal.GetPartKeys(ctx, client)
			if err != nil {
				return fmt.Errorf(
					style.Red.Render("failed to get participation keys: %s")+
						explanations.TokenNotAdmin,
					err)
			}
			state := internal.StateModel{
				Status: internal.StatusModel{
					State:       "INITIALIZING",
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
				ParticipationKeys: partkeys,

				Client:  client,
				Http:    new(internal.HttpPkg),
				Context: ctx,
			}
			state.Accounts, err = internal.AccountsFromState(&state, new(internal.Clock), client)
			cobra.CheckErr(err)
			// Fetch current state
			err = state.Status.Fetch(ctx, client, state.Http)
			cobra.CheckErr(err)

			m, err := ui.NewViewportViewModel(&state, client)
			cobra.CheckErr(err)

			p := tea.NewProgram(
				m,
				tea.WithAltScreen(),
				tea.WithFPS(120),
			)
			go func() {
				state.Watch(func(status *internal.StateModel, err error) {
					if err == nil {
						p.Send(state)
					}
					if err != nil {
						p.Send(state)
						p.Send(err)
					}
				}, ctx, client)
			}()
			_, err = p.Run()
			return err
		},
	}
)

func check(err interface{}) {
	if err != nil {
		panic(err)
	}
}

// Handle global flags and set usage templates
func init() {
	log.SetReportTimestamp(false)

	// Configure Version
	if Version == "" {
		Version = "unknown (built from source)"
	}
	rootCmd.Version = Version

	// Bindings
	rootCmd.Flags().StringVarP(&algod, "algod-endpoint", "a", "", style.LightBlue("algod endpoint address URI, including http[s]"))
	rootCmd.Flags().StringVarP(&token, "algod-token", "t", "", lipgloss.JoinHorizontal(
		lipgloss.Left,
		style.LightBlue("algod "),
		style.BoldUnderline("admin"),
		style.LightBlue(" token"),
	))
	_ = viper.BindPFlag("algod-endpoint", rootCmd.Flags().Lookup("algod-endpoint"))
	_ = viper.BindPFlag("algod-token", rootCmd.Flags().Lookup("algod-token"))

	// Update Long Text
	rootCmd.Long +=
		style.Magenta("Configuration: ") + viper.GetViper().ConfigFileUsed() + "\n" +
			style.LightBlue("Algod: ") + viper.GetString("algod-endpoint")

	if viper.GetString("data") != "" {
		rootCmd.Long +=
			style.Magenta("\nAlgorand Data: ") + viper.GetString("data")
	}

	// Add Commands
	rootCmd.AddCommand(statusCmd)
	if runtime.GOOS != "windows" {
		rootCmd.AddCommand(node.Cmd)
		rootCmd.AddCommand(configure.Cmd)
	}
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

type AlgodConfig struct {
	EndpointAddress string `json:"EndpointAddress"`
}

func replaceEndpointUrl(s string) string {
	s = strings.Replace(s, "\n", "", 1)
	s = strings.Replace(s, "0.0.0.0", "127.0.0.1", 1)
	s = strings.Replace(s, "[::]", "127.0.0.1", 1)
	return s
}
func hasWildcardEndpointUrl(s string) bool {
	return strings.Contains(s, "0.0.0.0") || strings.Contains(s, "::")
}
func initConfig() {
	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	// Look for paths
	viper.AddConfigPath(".")
	viper.AddConfigPath(home)
	viper.AddConfigPath("/etc/algorun/")

	// Set Config Properties
	viper.SetConfigType("yaml")
	viper.SetConfigName(".algorun")
	viper.SetEnvPrefix("algorun")

	// Load Configurations
	viper.AutomaticEnv()
	_ = viper.ReadInConfig()

	// Check for algod
	loadedAlgod := viper.GetString("algod-endpoint")
	loadedToken := viper.GetString("algod-token")

	// Load ALGORAND_DATA/config.json
	algorandData, exists := os.LookupEnv("ALGORAND_DATA")

	// Load the Algorand Data Configuration
	if exists && algorandData != "" && loadedAlgod == "" {
		// Placeholder for Struct
		var algodConfig AlgodConfig

		dataConfigPath := algorandData + "/config.json"

		// Open the config.json File
		configFile, err := os.Open(dataConfigPath)
		check(err)

		// Read the bytes of the File
		byteValue, _ := io.ReadAll(configFile)
		err = json.Unmarshal(byteValue, &algodConfig)
		check(err)

		// Close the open handle
		err = configFile.Close()
		check(err)

		// Check for endpoint address
		if hasWildcardEndpointUrl(algodConfig.EndpointAddress) {
			algodConfig.EndpointAddress = replaceEndpointUrl(algodConfig.EndpointAddress)
		} else if algodConfig.EndpointAddress == "" {
			// Assume it is not set, try to discover the port from the network file
			networkPath := algorandData + "/algod.net"
			networkFile, err := os.Open(networkPath)
			check(err)

			byteValue, err = io.ReadAll(networkFile)
			check(err)

			if hasWildcardEndpointUrl(string(byteValue)) {
				algodConfig.EndpointAddress = replaceEndpointUrl(string(byteValue))
			} else {
				algodConfig.EndpointAddress = string(byteValue)
			}

		}
		if strings.Contains(algodConfig.EndpointAddress, ":0") {
			algodConfig.EndpointAddress = strings.Replace(algodConfig.EndpointAddress, ":0", ":8080", 1)
		}
		if loadedToken == "" {
			// Handle Token Path
			tokenPath := algorandData + "/algod.admin.token"

			tokenFile, err := os.Open(tokenPath)
			check(err)

			byteValue, err = io.ReadAll(tokenFile)
			check(err)

			viper.Set("algod-token", strings.Replace(string(byteValue), "\n", "", 1))
		}

		// Set the algod configuration
		viper.Set("algod-endpoint", "http://"+strings.Replace(algodConfig.EndpointAddress, "\n", "", 1))
		viper.Set("data", dataConfigPath)
	}

}

func getClient() (*api.ClientWithResponses, error) {
	apiToken, err := securityprovider.NewSecurityProviderApiKey("header", "X-Algo-API-Token", viper.GetString("algod-token"))
	if err != nil {
		return nil, err
	}
	return api.NewClientWithResponses(viper.GetString("algod-endpoint"), api.WithRequestEditorFn(apiToken.Intercept))
}
