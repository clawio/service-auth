package commands

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/clawio/service-auth/sdk"
	"github.com/codegangsta/cli"
)

var VerifyCommand = cli.Command{
	Name:      "verify",
	Usage:     "Verifies an authentication token",
	UsageText: "verify <token>",
	Description: `
This command verifies an uuthentication token.
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

	authSDK, err := sdk.NewSDK(os.Getenv("CLAWIO_AUTH_ADDR"))
	if err != nil {
		log.Fatal(err)
	}
	identity, err := authSDK.Verify(args.Get(0))
	if err != nil {
		log.Fatal(err)
	}
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintf(w, "Username:\t%s\n", identity.Username)
	fmt.Fprintf(w, "Email:\t%s\n", identity.Email)
	fmt.Fprintf(w, "DisplayName:\t%s\n", identity.DisplayName)
	w.Flush()
}
