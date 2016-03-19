package commands

import (
	"bytes"
	"fmt"
	"github.com/codegangsta/cli"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

const HUMSAT20_DMH_WEBSERVICE string = "HUMSAT20_DMH_WEBSERVICE"

var ConvertCommand = cli.Command{
	Name:      "convert",
	Aliases:   []string{"conv", "con", "c"},
	Usage:     "Converts a markdown file to a PDF using a LaTeX template",
	UsageText: "Hello Hello Hello Usage text",
	Description: `
This command converts a file formatted in Pandoc´s Markdown to a PDF
applying a custom template. In order to see the list of available
templates launch the list-templates command. If -o is not specified
the data will be output to stdout. If -o is set the output will be written to
a local file.

Examples of usage:
pentas_md2pdf convert input.md pentas
pentas_md2pdf convert input.md pentas > /tmp/output.pdf
pentas_md2pdf convert input.md pentas -o /tmp/output.pdf
			`,
	ArgsUsage: "INPUT.MD TEMPLATE",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "output, o",
			Usage: "Write PDF to a local file",
		},
	},
	Action: convert,
}

func convert(c *cli.Context) {
	// validate input file and template args are non-empty
	args := c.Args()
	if len(args) < 2 || args.Get(0) == "" || args.Get(1) == "" {
		fmt.Println("Usage: cmd convert [command options] INPUT.MD TEMPLATE")
		os.Exit(1)
	}

	serverURI := os.Getenv(HUMSAT20_DMH_WEBSERVICE)
	if serverURI == "" {
		fmt.Println("You must set HUMSAT20_DMH_WEBSERVICE")
		os.Exit(1)
	}

	input := args.Get(0)
	template := args.Get(1)

	params := map[string]string{
		"template": template,
	}

	req, err := newfileUploadRequest(serverURI+"/convert", params, "srcfile", input)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if res.StatusCode != 200 {
		fmt.Println("There is a problem handling your request.")
		fmt.Println("Ensure your input file is a Markdown formatted file.")
		os.Exit(1)
	}

	var out io.Writer
	if c.String("output") == "" {
		out = os.Stdout
	} else {
		fd, err := os.Create(c.String("output"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		out = fd
	}
	_, err = io.Copy(out, res.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err

	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err

	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)

	}
	err = writer.Close()
	if err != nil {
		return nil, err

	}
	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}
