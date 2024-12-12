package main

import (
	"fmt"
	"github.com/algorandfoundation/algorun-tui/daemon/cmd"
)

var version = "development"

// main initializes the Echo framework, hides its default startup banner, prints a custom BANNER,
// registers the HTTP handlers, starts the algod process, and begins listening for HTTP requests on port 1323.
func main() {
	fmt.Println("Algorun TUI")
	fmt.Println("Version:", version)
	err := cmd.Execute()
	if err != nil {
		return
	}
}
