// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 alibaba/open-code-review Contributors

package main

import (
	"fmt"
	"os"

	"github.com/alibaba/open-code-review/internal/llm"
)

func main() {
	llm.AppVersion = Version
	os.Exit(run())
}

func run() int {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}
	return 0
}
