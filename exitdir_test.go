/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package exitdir_test

import (
	"context"
	"os"
	"testing"
	"time"

	"chainguard.dev/exitdir"
)

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
