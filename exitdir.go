/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package exitdir

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

// Aware watches for an exitFile and will cancel the returned context if
// `EDIT_DIR` has been set on the environment, otherwise context is returned.
func Aware(ctx context.Context) context.Context {
	exitDir, found := exitDir()
	if !found {
		return ctx
	}

	// Add WithCancel to our normal signal catching so that if we are
	// watching an exitfile that we can then cleanly exit.
	ctxWithCancel, cancel := context.WithCancel(ctx)

	go watchForExitDirAndExit(exitDir, cancel)
	return ctxWithCancel
}

// Exit creates an exitFile if `EXIT_DIR` has been set on the environment.
func Exit() error {
	exitDir, found := exitDir()
	if !found {
		return nil
	}
	_, err := os.Create(filepath.Join(exitDir, "exitFile"))
	return err
}

// get `EXIT_DIR` from environment.
func exitDir() (string, bool) {
	exitDir, found := os.LookupEnv("EXIT_DIR")
	if !found || exitDir == "" {
		return "", false
	}
	// Make sure exitDir a directory we're watching and not a file.
	return strings.TrimSuffix(exitDir, "/") + "/", true
}

// Watches for FS events for a given directory and if a file is created there
// will call the cancel function => exiting.
func watchForExitDirAndExit(exitDir string, cancel context.CancelFunc) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(fmt.Errorf("failed to create watcher: %w", err))
	}

	if err := watcher.Add(exitDir); err != nil {
		panic(fmt.Errorf("failed to add exit dir to watcher: %w", err))
	}

	// See if there are entries already there, say main container already
	// wrote the exit file.
	files, err := ioutil.ReadDir(exitDir)
	if err != nil {
		panic(fmt.Errorf("failed to read exitdir for files: %w", err))
	}

	if len(files) > 0 {
		cancel()
	} else {
		defer watcher.Close()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					// watcher channel was closed.
					return
				}
				// if we see a file create event, call cancel.
				if event.Op&fsnotify.Create == fsnotify.Create {
					cancel()
					return
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					// watcher error channel closed.
					return
				}
				_, _ = fmt.Fprintf(os.Stderr, "error from watcher.Errors: %v", err)
			}
		}
	}
}
