package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"

	gcli "github.com/grafana/grafana/pkg/cmd/grafana-cli/commands"
	gsrv "github.com/grafana/grafana/pkg/cmd/grafana-server/commands"
	"github.com/grafana/grafana/pkg/cmd/grafana/apiserver"
)

// The following variables cannot be constants, since they can be overridden through the -X link flag
var version = "9.2.0"
var commit = gcli.DefaultCommitValue
var enterpriseCommit = gcli.DefaultCommitValue
var buildBranch = "main"
var buildstamp string

func main() {
	app := &cli.App{
		Name:  "grafana",
		Usage: "Grafana server and command line interface",
		Authors: []*cli.Author{
			{
				Name:  "Grafana Project",
				Email: "hello@grafana.com",
			},
		},
		Version: version,
		Commands: []*cli.Command{
			gcli.CLICommand(version),
			gsrv.ServerCommand(version, commit, enterpriseCommit, buildBranch, buildstamp),
			{
				// The kubernetes standalone apiserver service runner
				// The command line is actually managed by cobra
				Name:  "apiserver",
				Usage: "run a standalone api service (experimental)",
				Action: func(context *cli.Context) error {
					return nil // not actually used
				},
			},
		},
		CommandNotFound:      cmdNotFound,
		EnableBashCompletion: true,
	}

	//
	if len(os.Args) > 1 && os.Args[1] == "apiserver" {
		apiserver.RunMain()
		return
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("%s: %s %s\n", color.RedString("Error"), color.RedString("✗"), err)
		os.Exit(1)
	}

	os.Exit(0)
}

func cmdNotFound(c *cli.Context, command string) {
	fmt.Printf(
		"%s: '%s' is not a %s command. See '%s --help'.\n",
		c.App.Name,
		command,
		c.App.Name,
		os.Args[0],
	)
	os.Exit(1)
}
