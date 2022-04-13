/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"time"

	"chainguard.dev/exitdir"
)

func main() {
	fmt.Println("[Work] Doing work...")
	time.Sleep(5 * time.Second)
	fmt.Println("[Work] Exiting...")
	if err := exitdir.Exit(); err != nil {
		panic(err)
	}
}
