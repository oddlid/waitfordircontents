package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

const (
	optEnvVar          = "env-var"
	optDir             = "directories"
	optTimeout         = "timeout"
	optExitOnWatchFail = "exit-on-watch-failure"
)

func main() {
	app := &cli.App{
		Name:                 "waitfordircontents",
		Usage:                "Wait until given directories are not empty",
		Copyright:            "(C) 2022 Odd Eivind Ebbesen",
		EnableBashCompletion: true,
		Authors: []*cli.Author{
			{
				Name:  "Odd Eivind Ebbesen",
				Email: "oddebb@gmail.com",
			},
		},
		Before: func(ctx *cli.Context) error {
			zerolog.TimeFieldFormat = "2006-01-02T15:04:05.999-07:00"
			return nil
		},
		Action: entryPoint,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    optEnvVar,
				Usage:   "Environment `variable` to get paths from. Values should be colon separated.",
				Aliases: []string{"e"},
			},
			&cli.StringSliceFlag{
				Name:    optDir,
				Usage:   "Directories to watch. Separate by commas, or specify multiple times.",
				Aliases: []string{"d"},
			},
			&cli.DurationFlag{
				Name:    optTimeout,
				Usage:   "How long to wait before giving up. 0 means wait forever.",
				Aliases: []string{"t"},
			},
			&cli.BoolFlag{
				Name:    optExitOnWatchFail,
				Usage:   "Exit with error if a watch failure happens",
				Aliases: []string{"x"},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Send()
	}
}
