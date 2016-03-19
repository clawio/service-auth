package main

import (
	"github.com/codegangsta/cli"
	"github.com/labkode/Humsat20-dmh-webservice/pentas_md2pdf/commands"
	"os"
)

var VERSION string

func main() {

	app := cli.NewApp()
	app.Version = VERSION
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Hugo González Labrador",
			Email: "contact@hugo.labkode.com",
		},
	}
	app.Copyright = `
	Universidade de Vigo owns the copyright of this tool which is supplied
	in confidence and which shall not be used for any purpose other than that
	for which it is supplied and shall not in whole or in part be reproduced,
	copied or communicated to any person without permission from the owner.
	`

	app.Name = "Pentas Markdown to PDF converter"
	app.Usage = `

	This tool converts an input file formatted in Pandoc's Markdown
	to a PDF file using a custom LaTeX template.

	This tool send requests to a Humsat20-dmh-webservice server in order to
	perform the conversion among other functions.
	Therefore, to stablish a connection with the server the HUMSAT20_DMH_WEBSERVICE
	environmental variable must be set.

	Example:
	export HUMSAT20_DMH_WEBSERVICE=http://localhost:57008
	`
	app.EnableBashCompletion = true

	app.Commands = []cli.Command{
		commands.ConvertCommand,
		commands.TemplateCommand,
	}

	app.Run(os.Args)
}
