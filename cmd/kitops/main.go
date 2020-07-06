package main

import (
	"log"
	"os"

	"github.com/300481/kitops/pkg/kitops"
	cli "github.com/urfave/cli/v2"
)

var (
	app = cli.NewApp()
)

func init() {
}

func info() {
	app.Name = "Kitops"
	app.Usage = ""
	app.Version = "0.1.0"
}

func commands() {
	app.Commands = []*cli.Command{
		{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "Run in server mode",
			Action: func(c *cli.Context) error {
				kitops.New().Serve()
				return nil
			},
		},
	}
}

func main() {
	info()
	commands()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
