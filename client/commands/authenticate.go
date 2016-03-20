package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/clawio/service-auth/sdk"
	"github.com/codegangsta/cli"
)

var AuthenticateCommand = cli.Command{
	Name:      "authenticate",
	Usage:     "Authenticates a user with username/password against a ClawIO Service Auth",
	UsageText: "authenticate john johnpasswd",
	Description: `
This command authenticates a user with a username and password.
			`,
	ArgsUsage: "<username> <password>",
	Action:    authenticate,
}

func authenticate(c *cli.Context) {
	args := c.Args()
	if len(args) < 2 || args.Get(0) == "" || args.Get(1) == "" {
		fmt.Println("Usage: authenticate  <username> <password>")
		os.Exit(1)
	}

	authSDK, err := sdk.NewSDK(os.Getenv("CLAWIO_AUTH_ADDR"))
	if err != nil {
		log.Fatal(err)
	}
	token, err := authSDK.Authenticate(args.Get(0), args.Get(1))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("API Token: ", token)
}
