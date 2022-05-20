package main

import (
	"context"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
)

type watcher struct {
	wg sync.WaitGroup
}

func (w *watcher) start(ctx context.Context, paths []string) error {
	for _, path := range paths {
		if err := w.watch(ctx, path); err != nil {
			return err
		}
	}
	return nil
}

func (w *watcher) watch(ctx context.Context, path string) error {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	err = fsWatcher.Add(path)
	if err != nil {
		return err
	}

	w.wg.Add(1)
	go func() {
		defer fsWatcher.Close()
		defer w.wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-fsWatcher.Events:
				if !ok {
					return
				}
				log.Info().Msg(event.String())
				return
			case err, ok := <-fsWatcher.Errors:
				if !ok {
					return
				}
				log.Info().Err(err).Send()
				return
			}
		}
	}()
	return nil
}

func (w *watcher) wait() {
	w.wg.Wait()
}
