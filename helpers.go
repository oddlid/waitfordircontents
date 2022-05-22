package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/signal"
	"strings"
	"syscall"
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

func dirIsEmpty(path string) (bool, error) {
	d, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer d.Close()

	_, err = d.Readdirnames(1)
	if err != nil && errors.Is(err, io.EOF) {
		return true, nil
	}
	return false, err
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

func setupSignalListening() <-chan os.Signal {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	return signals
}
