package commands

import (
	"github.com/regmicmahesh/merosharemorelikeidontcare/common"
	"github.com/urfave/cli/v2"
)

func InitializeCredentials(username, password string, clientID int) error {
	_, err := common.Login(username, password, clientID)
	if err != nil {
		return err
	}

	_, err = common.Hydrate(username, password, clientID)
	if err != nil {
		return err
	}
	return nil

}

var InitCommand = &cli.Command{
	Name:    "init",
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
		username := c.String("username")
		password := c.String("password")
		clientID := c.Int("clientID")
		return InitializeCredentials(username, password, clientID)
	},
}
