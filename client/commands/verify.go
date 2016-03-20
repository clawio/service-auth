package commands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

var VerifyCommand = cli.Command{
	Name:      "verify",
	Usage:     "Verifies an authentication token and shows user identity",
	UsageText: "verify <token>",
	Description: `
Verifies an authentication token and shows user identity
`,
	ArgsUsage: "<token>",
	Action:    verify,
}

func verify(c *cli.Context) {
	args := c.Args()
	if len(args) < 1 || args.Get(0) == "" {
		fmt.Println("Usage: verify  <token>")
		os.Exit(1)
	}
}
