package server

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func Run() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "data",
				Value: "data.json",
				Usage: "Path to static data file",
			},
			&cli.StringFlag{
				Name:  "listen",
				Value: ":8000",
				Usage: "Address to listen on",
			},
		},
		Action: func(c *cli.Context) error {
			s := newServer(c.String("data"))
			s.Run(c.String("listen"))
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
