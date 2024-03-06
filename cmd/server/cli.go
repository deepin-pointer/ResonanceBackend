package server

import (
	"encoding/hex"

	"golang.org/x/crypto/sha3"

	"rsbackend/internal/model"

	"github.com/desertbit/grumble"
	"github.com/spf13/viper"
)

func Run() error {

	var s *Server

	model.InitConfig()

	var app = grumble.New(&grumble.Config{
		Name:        "rsbackend",
		Description: "Backend server for Resonance Solitice game stats platform",

		Flags: func(f *grumble.Flags) {
			f.String("c", "config", "config.json", "Path to config file")
		},
	})

	app.OnInit(func(a *grumble.App, flags grumble.FlagMap) error {
		model.InitConfig()
		path := flags.String("config")
		if path != "" {
			viper.SetConfigFile(path)
		}
		s = newServer()
		go s.Run(viper.GetString("bind"))
		return nil
	})

	app.OnClosing(func() error {
		s.Shutdown()
		return nil
	})

	app.AddCommand(&grumble.Command{
		Name: "adduser",
		Help: "Add user for login",

		Args: func(a *grumble.Args) {
			a.String("name", "User name")
			a.String("password", "Password for user")
		},

		Run: func(c *grumble.Context) error {
			name := c.Args.String("name")
			pass := c.Args.String("password")
			if len(name) > 0 && len(pass) > 0 {
				digest := sha3.Sum512([]byte(pass))
				viper.GetStringMapString("users")[name] = hex.EncodeToString(digest[:])
			}
			viper.WriteConfig()
			return nil
		},
	})

	return app.Run()
}
