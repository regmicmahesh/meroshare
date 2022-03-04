package main

import (
	"os"

	"github.com/regmicmahesh/merosharemorelikeidontcare/commands"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:     "meroshare-cli",
		Usage:    "meroshare-cli is a command line interface for Meroshare.",
		HelpName: "meroshare-cli",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "output",
				Usage: "output in json or ascii table",
				Value: "ascii",
			},
		},
		Description: "Meroshare CLI - Reject Angular and Bloated Web and embrace the power of CLI.",
		Commands: []*cli.Command{
			commands.DetailsCommand,
			commands.PortfolioCommand,
			commands.InitCommand,
		},
	}
	app.Run(os.Args)

}
