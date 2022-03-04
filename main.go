package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:        "meroshare-cli",
		Usage:       "meroshare-cli is a command line interface for Meroshare.",
		HelpName:    "meroshare-cli",
		Description: "Meroshare CLI - Reject Angular and Bloated Web and embrace the power of CLI.",
		Commands: []*cli.Command{
			{
				Name:    "details",
				Aliases: []string{"t"},
				Usage:   "Get details of the user.",
				Action: func(c *cli.Context) error {
					_, err := getOwnDetails(true)
					if err != nil {
						fmt.Println("❌", err)
						return err
					}
					return nil
				},
			},
			{
				Name:    "portfolio",
				Aliases: []string{"t"},
				Usage:   "Print portfolio.",
				Action: func(c *cli.Context) error {
					err := printPortfolio()
					if err != nil {
						fmt.Println("❌", err)
						return err
					}
					return nil
				},
			},
			{
				Name:    "test",
				Aliases: []string{"t"},
				Usage:   "Test if credentials are working.",
				Action: func(c *cli.Context) error {
					err := loadToken()
					if err != nil {
						fmt.Println("❌", err)
						return err
					}
					return nil
				},
			},
			{
				Name:    "init",
				Aliases: []string{"i"},
				Usage:   "Initialize credentials for Meroshare.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						EnvVars:  []string{"MEROSHARE_USERNAME"},
						Name:     "username",
						Usage:    "Username for Meroshare.",
						Required: true,
					},
					&cli.StringFlag{
						EnvVars:  []string{"MEROSHARE_PASSWORD"},
						Name:     "password",
						Usage:    "Password for Meroshare.",
						Required: true,
					},
					&cli.IntFlag{
						EnvVars:  []string{"MEROSHARE_CLIENT_ID"},
						Name:     "clientID",
						Usage:    "Client ID for Meroshare.",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {

					fmt.Println("✔️ All flags are set.")

					username := c.String("username")
					password := c.String("password")
					clientID := c.Int("clientID")

					_, err := login(username, password, clientID)

					if err != nil {
						return err
					}

					fmt.Println("✔️ Logged in successfully.")

					data := &struct {
						Username string `json:"username"`
						Password string `json:"password"`
						ClientID int    `json:"clientId"`
					}{
						Username: username,
						Password: password,
						ClientID: clientID,
					}
					homedir, err := os.UserHomeDir()

					if err != nil {
						fmt.Println(err)
					}

					file, err := os.Create(path.Join(homedir, "credentials.json"))
					if err != nil {
						return err
					}

					json.NewEncoder(file).Encode(data)

					defer file.Close()

					fmt.Println("✔️ Successfully Initialized.")

					return nil

				},
			},
		},
	}
	app.Run(os.Args)
}
