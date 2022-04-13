/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"chainguard.dev/exitdir"
	"context"
	"fmt"
	"os"
	"time"
)

func main() {
	ctx := exitdir.Aware(context.Background())
	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("[Busy] Exiting...")
			os.Exit(0)
		case t := <-ticker.C:
			fmt.Println("[Busy] Tick at", t)
		}
	}
}
