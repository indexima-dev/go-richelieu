package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/estebgonza/go-richelieu/constants"
	"github.com/estebgonza/go-richelieu/generator"
	"github.com/urfave/cli/v2"
)

const helpTemplate = `
Usage: {{.HelpName}} [command]

{{if .Commands}}Commands:

{{range .Commands}}{{if not .HideHelp}}{{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}
`

func main() {
	cli.AppHelpTemplate = fmt.Sprintf(helpTemplate)
	app := cli.NewApp()
	app.Name = constants.AppName
	app.Usage = constants.AppDescription
	app.Version = constants.AppVersion

	app.Commands = []*cli.Command{
		{
			Name:    "generate",
			Usage:   "Generate the dataset from the plan.json input",
			Aliases: []string{"g"},
			Action:  func(c *cli.Context) error { return generator.Generate() },
		},
		{
			Name:    "readFromColumn",
			Usage:   "Create a plan.json from column list argument",
			Aliases: []string{"rc"},
			Action:  func(c *cli.Context) error { return generator.CreateFromColumn(c.Args()) },
		},
	}

	app.CommandNotFound = func(c *cli.Context, command string) {
		fmt.Printf("Command not found: %v\n", command)
		cli.ShowAppHelp(c)
	}

	log.Println("Starting...")
	startTime := time.Now()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Done in", time.Since(startTime))
}
