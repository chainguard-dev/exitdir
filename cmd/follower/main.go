/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"chainguard.dev/exitdir"
)

func main() {
	ctx := exitdir.Aware(context.Background())
	ticker := time.NewTicker(1 * time.Second)
	for i := 0; true; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("[Follower] Exiting...")
			os.Exit(0)
		case <-ticker.C:
			fmt.Println("[Follower] Tick", i)
		}
	}
}
