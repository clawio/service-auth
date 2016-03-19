package commands

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"net/http"
	"os"
)

// TemplateCommand handles the templates command.
var TemplateCommand = cli.Command{
	Name:      "templates",
	Aliases:   []string{"templ", "temp", "t"},
	Usage:     "Shows a list of available templates",
	UsageText: "Hello Hello Hello Usage text",
	Description: `
This command shows a list of the available templates to apply to the
final .pdf file.

Examples of usage:
pentas_md2pdf templates
			`,
	ArgsUsage: "",
	Action:    templates,
}

func templates(c *cli.Context) {

	serverURI := os.Getenv(HUMSAT20_DMH_WEBSERVICE)
	if serverURI == "" {
		fmt.Println("You must set HUMSAT20_DMH_WEBSERVICE")
		os.Exit(1)
	}

	res, err := http.Get(serverURI + "/templates")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if res.StatusCode != 200 {
		fmt.Println("There is a problem handling your request.")
		fmt.Println("Ensure your input file is a Markdown formatted file.")
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	data := []string{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Available templates:")
	for _, v := range data {
		fmt.Printf("- %s\n", v)
	}
}
