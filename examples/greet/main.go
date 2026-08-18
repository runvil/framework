// Example CLI application demonstrating the Runvil CLI framework.
//
// Run with:
//
//	go run ./examples/greet hello --name Alice
//	GREET_GREETING=Halo go run ./examples/greet hello --name Alice
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/runvil/runvil-framework/cli"
	"github.com/runvil/runvil-libs/core"
)

func main() {
	app := cli.NewApp("greet", "0.1.0").
		Command(cli.NewCommand("hello", "greet a user", registerHello, hello))

	os.Exit(int(app.Run(os.Args[1:])))
}

func registerHello(fs *flag.FlagSet) {
	fs.String("name", "", "name of the user to greet")
}

func hello(fs *flag.FlagSet) core.ExitCode {
	name := fs.Lookup("name").Value.String()
	greeting := os.Getenv("GREET_GREETING")
	if greeting == "" {
		greeting = "Hello"
	}
	if name == "" {
		fmt.Fprintln(os.Stderr, "usage: greet hello --name <name>")
		return core.ExitCodeUsage
	}
	fmt.Printf("%s, %s!\n", greeting, name)
	return core.ExitCodeSuccess
}
