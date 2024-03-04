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
				Name:  "static",
				Value: "static.json",
				Usage: "Path to static data file",
			},
			&cli.StringFlag{
				Name:  "dynamic",
				Value: "dynamic.bin",
				Usage: "Path to dynamic data log file",
			},
			&cli.StringFlag{
				Name:  "listen",
				Value: ":8000",
				Usage: "Address to listen on",
			},
		},
		Action: func(c *cli.Context) error {
			s := newServer(c.String("static"), c.String("dynamic"))
			s.Run(c.String("listen"))
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
