package main

import (
	"context"

	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func entryPoint(c *cli.Context) error {
	timeout := c.Duration(optTimeout)
	dirsFromCli := c.StringSlice(optDir)
	dirsFromEnvVar := getPathsFromEnvVar(c.String(optEnvVar))
	exitOnWatchFail := c.Bool(optExitOnWatchFail)
	paths := []string{}

	if len(dirsFromCli) > 0 {
		paths = append(paths, dirsFromCli...)
	}
	if len(dirsFromEnvVar) > 0 {
		paths = append(paths, dirsFromEnvVar...)
	}

	log.Info().Strs("dirs", paths).Msg("Directories to watch")

	if err := verifyDirs(paths); err != nil {
		return err
	}

	var ctx context.Context
	var cancel context.CancelFunc
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel()

	signals := setupSignalListening()
	eventFilter := func(e fsnotify.Event) bool {
		return e.Op&fsnotify.Create != 0 || e.Op&fsnotify.Write != 0
	}
	dirWatcher := watcher{}
	errChan, err := dirWatcher.start(ctx, paths, eventFilter)
	if err != nil {
		return err
	}
	done := dirWatcher.wait()

	for {
		select {
		case sig := <-signals:
			log.Info().Str("signal", sig.String()).Msg("Got signal, exiting")
			cancel()
		case err := <-errChan:
			log.Error().Err(err).Msg("Watch error")
			if exitOnWatchFail {
				log.Warn().Msg("Set to exit on watch failure, bailing out")
				cancel()
				<-done
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		case <-done:
			log.Info().Msg("All done")
			return nil
		}
	}
}
