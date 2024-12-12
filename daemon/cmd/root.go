package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/algorand/go-algorand/config"
	"github.com/algorandfoundation/algorun-tui/daemon/fortiter"
	"github.com/algorandfoundation/algorun-tui/daemon/rpc"
	"github.com/algorandfoundation/algorun-tui/ui/style"
	"github.com/charmbracelet/log"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"os"
	"strings"
)

const BANNER = `
 ______   ______     ______     ______   __     ______   ______     ______    
/\  ___\ /\  __ \   /\  == \   /\__  _\ /\ \   /\__  _\ /\  ___\   /\  == \   
\ \  __\ \ \ \/\ \  \ \  __<   \/_/\ \/ \ \ \  \/_/\ \/ \ \  __\   \ \  __<   
 \ \_\    \ \_____\  \ \_\ \_\    \ \_\  \ \_\    \ \_\  \ \_____\  \ \_\ \_\ 
  \/_/     \/_____/   \/_/ /_/     \/_/   \/_/     \/_/   \/_____/   \/_/ /_/ 
`

var version = ""
var (
	Version = version
	sqlFile string
	rootCmd = &cobra.Command{
		Version: Version,
		Use:     "fortiter",
		Short:   "Consume all the data",
		Long:    style.Purple(BANNER) + "\n",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.SetOutput(cmd.OutOrStdout())

			e := echo.New()
			e.HideBanner = true

			// TODO: handle user interaction
			e.Use(middleware.Logger())
			e.Use(middleware.Recover())
			e.Use(echoprometheus.NewMiddleware("fortiter"))
			e.Static("/", "public")

			fmt.Println(style.Magenta(BANNER))
			fmt.Println(style.LightBlue("Database: ") + viper.GetString("database"))

			algodConfig := viper.GetString("data")
			if algodConfig != "" {
				fmt.Println(style.LightBlue("Configuration: ") + algodConfig)
			}
			logFile := viper.GetString("log")
			if logFile != "" {
				fmt.Println(style.LightBlue("Log: ") + logFile)
			}
			var si = fortiter.Handlers{PrometheusHandler: echoprometheus.NewHandler()}
			rpc.RegisterHandlers(e, si)

			db, err := sqlx.Connect("sqlite3", viper.GetString("database"))
			cobra.CheckErr(err)

			// exec the schema or fail; multi-statement Exec behavior varies between
			// database drivers;  pq will exec them all, sqlite3 won't, ymmv
			db.MustExec(fortiter.Schema)
			db.MustExec(fortiter.StatsSchema)

			err = fortiter.Sync(context.Background(), logFile, db)
			cobra.CheckErr(err)

			return e.Start(":1337")
		},
	}
)

func check(err interface{}) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

// Handle global flags and set usage templates
func init() {
	log.SetReportTimestamp(false)
	initConfig()
	// Configure Version
	if Version == "" {
		Version = "unknown (built from source)"
	}
	rootCmd.Version = Version

	// Bindings
	rootCmd.PersistentFlags().StringVar(&sqlFile, "database", "fortiter.db", style.LightBlue("database file location"))
	_ = viper.BindPFlag("database", rootCmd.PersistentFlags().Lookup("database"))
	// Update Long Text
	rootCmd.Long +=
		//style.Magenta("Database: ") + viper.GetViper().ConfigFileUsed() + "\n" +
		style.LightBlue("Database: ") + viper.GetString("database")

	if viper.GetString("data") != "" {
		rootCmd.Long +=
			style.Magenta("\nAlgorand Data: ") + viper.GetString("data")
	}
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

type AlgodConfig struct {
	EndpointAddress string `json:"EndpointAddress"`
}

func initExistingAlgod() {
	algorandData, exists := os.LookupEnv("ALGORAND_DATA")

	// Load the Algorand Data Configuration
	if exists && algorandData != "" {
		// Placeholder for Struct
		var algodConfig config.Local

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

		// Replace catchall address with localhost
		if strings.Contains(algodConfig.EndpointAddress, "0.0.0.0") {
			algodConfig.EndpointAddress = strings.Replace(algodConfig.EndpointAddress, "0.0.0.0", "127.0.0.1", 1)
		}

		// Handle Token Path
		tokenPath := algorandData + "/algod.admin.token"

		tokenFile, err := os.Open(tokenPath)
		check(err)

		byteValue, err = io.ReadAll(tokenFile)
		check(err)
	}
}
func initConfig() {
	// Load ALGORAND_DATA/config.json
	algorandData, exists := os.LookupEnv("ALGORAND_DATA")

	// Load the Algorand Data Configuration
	if exists && algorandData != "" {
		// Placeholder for Struct
		var algodConfig config.Local

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

		// Find the log file
		logPath, _ := algodConfig.ResolveLogPaths(algorandData)

		// Handle Token Path
		tokenPath := algorandData + "/algod.admin.token"

		tokenFile, err := os.Open(tokenPath)
		check(err)

		byteValue, err = io.ReadAll(tokenFile)
		check(err)

		// Set the server configuration
		viper.Set("server", "http://"+algodConfig.EndpointAddress)
		viper.Set("token", string(byteValue))
		viper.Set("data", dataConfigPath)
		viper.Set("log", logPath)
	}

}
