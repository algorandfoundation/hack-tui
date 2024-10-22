package main

import (
	"fmt"
	ptyDevice "github.com/creack/pty"
	"github.com/labstack/echo/v4"
	"os"
	"os/exec"
)

const BANNER = `
   _____  .__                __________              
  /  _  \ |  |    ____   ____\______   \__ __  ____  
 /  /_\  \|  |   / ___\ /  _ \|       _/  |  \/    \ 
/    |    \  |__/ /_/  >  <_> )    |   \  |  /   |  \
\____|__  /____/\___  / \____/|____|_  /____/|___|  /
        \/     /_____/               \/           \/ 
`

// printPty continuously reads from the provided *os.File and prints the output as byte count and string to standard output.
func printPty(f *os.File) {
	b1 := make([]byte, 2048)
	for {
		n1, _ := f.Read(b1)
		fmt.Printf("%d bytes: %s\n", n1, string(b1[:n1]))
	}
}

// runAlgod starts the `algod` process using the default system path and pipes its output to the printPty function.
func runAlgod() {
	p, _ := exec.LookPath("algod")
	cmd := exec.Command(p, "-o")
	//ptyDevice.Open()
	f, _ := ptyDevice.Start(cmd)
	ptyDevice.Open()
	go printPty(f)
}

// main initializes the Echo framework, hides its default startup banner, prints a custom BANNER,
// registers the HTTP handlers, starts the algod process, and begins listening for HTTP requests on port 1323.
func main() {
	e := echo.New()
	e.HideBanner = true
	fmt.Println(BANNER)

	var si = Handlers{}
	RegisterHandlers(e, si)

	runAlgod()

	e.Logger.Fatal(e.Start(":1323"))
}
