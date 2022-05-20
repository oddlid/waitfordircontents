package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

const (
	optEnvVar  = "env-var"
	optDir     = "directories"
	optTimeout = "timeout"
)

var (
	errInvalidDir = errors.New("directory does not exist or is not a directory")
	errNoDirs     = errors.New("no directories given")
)

func splitByColon(input string) []string {
	return strings.Split(input, ":")
}

func getPathsFromEnvVar(name string) []string {
	paths := []string{}
	value, ok := os.LookupEnv(name)
	if ok {
		value = strings.TrimSpace(value)
		if value != "" {
			paths = append(paths, splitByColon(value)...)
		}
	}
	return paths
}

func dirExists(path string) bool {
	path = strings.TrimSpace(path)
	if path == "" {
		return false
	}
	info, err := os.Stat(path)
	if err != nil {
		return !errors.Is(err, fs.ErrNotExist)
	}
	return info.IsDir()
}

func verifyDirs(paths []string) error {
	if paths == nil || len(paths) == 0 {
		return errNoDirs
	}
	for _, path := range paths {
		if !dirExists(path) {
			return fmt.Errorf("%w: %s", errInvalidDir, path)
		}
	}
	return nil
}

func entryPoint(ctx *cli.Context) error {
	paths := []string{}
	dirsFromCli := ctx.StringSlice(optDir)
	dirsFromEnvVar := getPathsFromEnvVar(ctx.String(optEnvVar))
	if len(dirsFromCli) > 0 {
		paths = append(paths, dirsFromCli...)
	}
	if len(dirsFromEnvVar) > 0 {
		paths = append(paths, dirsFromEnvVar...)
	}
	log.Info().Strs("dirs", paths).Msg("Directories to watch")

	// Create context, with timeout if given. Pass this to the watcher, so we have a way to quit
	// if it either times out, or if we get a quit signal from the shell

	return verifyDirs(paths)
}

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
				Name:  optTimeout,
				Usage: "How long to wait before giving up. 0 means wait forever.",
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Send()
	}
}
