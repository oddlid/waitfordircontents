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

type eventFilterFunc func(e fsnotify.Event) bool

func (w *watcher) start(ctx context.Context, paths []string, eventFilter eventFilterFunc) (<-chan error, error) {
	errChan := make(chan error, 1)
	for _, path := range paths {
		empty, err := dirIsEmpty(path)
		if err != nil {
			return errChan, err
		}
		if !empty {
			log.Warn().Str("dir", path).Msg("Directory has content, skipping")
			continue
		}
		if err := w.watch(ctx, path, errChan, eventFilter); err != nil {
			return errChan, err
		}
	}
	return errChan, nil
}

func (w *watcher) watch(ctx context.Context, path string, errChan chan<- error, eventFilter eventFilterFunc) error {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	if err := fsWatcher.Add(path); err != nil {
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
				if !eventFilter(event) {
					log.Info().Str("event", event.String()).Msg("Skipped event")
					continue
				}
				log.Info().Str("event", event.String()).Msg("Event fulfilled, returning")
				return
			case err, ok := <-fsWatcher.Errors:
				if !ok {
					return
				}
				errChan <- err
				return
			}
		}
	}()

	return nil
}

func (w *watcher) wait() <-chan struct{} {
	done := make(chan struct{})
	go func() {
		w.wg.Wait()
		close(done)
	}()
	return done
}
