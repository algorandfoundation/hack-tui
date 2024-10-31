package main

import (
	"github.com/algorandfoundation/hack-tui/daemon/cmd"
)

// main initializes the Echo framework, hides its default startup banner, prints a custom BANNER,
// registers the HTTP handlers, starts the algod process, and begins listening for HTTP requests on port 1323.
func main() {
	err := cmd.Execute()
	if err != nil {
		return
	}
}
