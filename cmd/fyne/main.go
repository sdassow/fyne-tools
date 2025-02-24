// Run a command line helper for various Fyne tools.
package main

import (
	"fmt"
	"os"

	"fyne.io/tools/cmd/fyne/internal/commands"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:        "fyne",
		Usage:       "A command line helper for various Fyne tools.",
		Description: "The fyne command provides tooling for fyne applications and to assist in their development.",
		Commands: []*cli.Command{
			commands.Bundle(),
			commands.Env(),
			commands.Install(),
			commands.Package(),
			commands.Release(),
			commands.Version(),
			commands.Serve(),
			commands.Translate(),
			commands.Build(),
			commands.Init(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
