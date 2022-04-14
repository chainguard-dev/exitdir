/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package exitdir_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"chainguard.dev/exitdir"
)

func Example_aware_exit() {
	tempdir, err := ioutil.TempDir(os.TempDir(), "example_aware_exit*")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tempdir)
	if err := os.Setenv("EXIT_DIR", tempdir); err != nil {
		panic(err)
	}

	time.Sleep(1 * time.Millisecond)

	go func() {
		ctx := exitdir.Aware(context.Background())
		ticker := time.NewTicker(5 * time.Millisecond)
		for i := 0; true; i++ {
			select {
			case <-ctx.Done():
				fmt.Println("[Follower] Exiting...")
				return
			case <-ticker.C:
				fmt.Println("[Follower] Tick", i)
			}
		}
	}()

	fmt.Println("[Leader] Doing work...")
	time.Sleep(13 * time.Millisecond)
	fmt.Println("[Leader] Exiting...")
	if err := exitdir.Exit(); err != nil {
		panic(err)
	}
	time.Sleep(1 * time.Millisecond)

	// Output:
	// [Leader] Doing work...
	// [Follower] Tick 0
	// [Follower] Tick 1
	// [Leader] Exiting...
	// [Follower] Exiting...
}

func TestAwareExit_empty(t *testing.T) {
	if err := os.Setenv("EXIT_DIR", t.TempDir()); err != nil {
		t.Fatal(err)
	}

	ctx := exitdir.Aware(context.TODO())
	ended := false

	go func() {
		<-ctx.Done()
		ended = true
	}()

	if err := exitdir.Exit(); err != nil {
		t.Fatal(err)
	}

	// Let the filesystem catch up.
	time.Sleep(1 * time.Millisecond)

	if !ended {
		t.Errorf("expected the thread to end")
	}
}

func TestAwareExit_existingFile(t *testing.T) {
	if err := os.Setenv("EXIT_DIR", t.TempDir()); err != nil {
		t.Fatal(err)
	}

	if err := exitdir.Exit(); err != nil {
		t.Fatal(err)
	}

	ctx := exitdir.Aware(context.TODO())
	ended := false

	go func() {
		<-ctx.Done()
		ended = true
	}()

	// Let the filesystem catch up.
	time.Sleep(1 * time.Millisecond)

	if !ended {
		t.Errorf("expected the thread to end")
	}
}
