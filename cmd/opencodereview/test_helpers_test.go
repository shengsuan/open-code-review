// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 alibaba/open-code-review Contributors

package main

import (
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/alibaba/open-code-review/internal/config/template"
	"github.com/spf13/cobra"
)

func loadTestTemplate(t *testing.T) *template.Template {
	t.Helper()
	tpl, err := template.LoadDefault()
	if err != nil {
		t.Fatalf("LoadDefault: %v", err)
	}
	return tpl
}

func captureDelegateStdout(t *testing.T, fn func()) []byte {
	t.Helper()
	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("create stdout pipe: %v", err)
	}
	os.Stdout = w
	defer func() { os.Stdout = orig }()

	var out []byte
	var readErr error
	done := make(chan struct{})
	go func() {
		defer close(done)
		out, readErr = io.ReadAll(r)
	}()
	fn()
	if err := w.Close(); err != nil {
		t.Fatalf("close stdout writer: %v", err)
	}
	<-done
	if readErr != nil {
		t.Fatalf("read stdout: %v", readErr)
	}
	_ = r.Close()
	return out
}

func gitCommitFile(t *testing.T, dir, name, content, msg string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
	for _, args := range [][]string{{"add", "."}, {"commit", "-m", msg}} {
		cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", args, err, out)
		}
	}
}

func silenceStdout(t *testing.T, fn func()) {
	t.Helper()
	orig := os.Stdout
	devnull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		t.Fatalf("open devnull: %v", err)
	}
	os.Stdout = devnull
	defer func() {
		os.Stdout = orig
		_ = devnull.Close()
	}()
	fn()
}

func setTestHome(t *testing.T, dir string) {
	t.Helper()
	t.Setenv("HOME", dir)
	t.Setenv("USERPROFILE", dir)
}

func freshOCRHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	setTestHome(t, home)
	t.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
	t.Setenv("GIT_CONFIG_SYSTEM", os.DevNull)
	return home
}

func assertNoSessionStore(t *testing.T, home string) {
	t.Helper()
	store := filepath.Join(home, ".opencodereview", "sessions")
	if _, err := os.Stat(store); !os.IsNotExist(err) {
		t.Errorf("preview created the session store at %s (stat err = %v)", store, err)
	}
}

// runReview provides test compatibility for the pre-cobra runReview(args)
// entry point: parse flags, then execute the review.
func runReview(args []string) error {
	opts, err := parseReviewFlags(args)
	if err != nil {
		return err
	}
	return executeReview(opts)
}

// runConfig provides test compatibility: dispatches config subcommands through a
// fresh cobra command tree that mirrors the production config command, so the
// test compat layer cannot drift from production routing.
func runConfig(args []string) error {
	cmd := &cobra.Command{Use: "config", SilenceUsage: true, SilenceErrors: true}
	setCmd := &cobra.Command{
		Use: "set", Args: cobra.ExactArgs(2),
		SilenceUsage: true, SilenceErrors: true,
		RunE: func(cmd *cobra.Command, a []string) error { return runConfigSet(a[0], a[1]) },
	}
	unsetCmd := &cobra.Command{
		Use: "unset", Args: cobra.ExactArgs(1),
		SilenceUsage: true, SilenceErrors: true,
		RunE: func(cmd *cobra.Command, a []string) error { return runConfigUnset(a[0]) },
	}
	providerCmd := &cobra.Command{
		Use: "provider", Args: cobra.NoArgs,
		SilenceUsage: true, SilenceErrors: true,
		RunE: func(cmd *cobra.Command, a []string) error { return runConfigProvider() },
	}
	modelCmd := &cobra.Command{
		Use: "model", Args: cobra.NoArgs,
		SilenceUsage: true, SilenceErrors: true,
		RunE: func(cmd *cobra.Command, a []string) error { return runConfigModel() },
	}
	cmd.AddCommand(setCmd, unsetCmd, providerCmd, modelCmd)
	// A nil args slice makes cobra fall back to os.Args; pass an explicit empty
	// slice so `runConfig(nil)` shows usage instead of reading the test binary's args.
	if args == nil {
		args = []string{}
	}
	cmd.SetArgs(args)
	return cmd.Execute()
}

func parseReviewFlags(args []string) (reviewOptions, error) {
	var opts reviewOptions
	cmd := &cobra.Command{
		Use:           "review",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return validateReviewOptions(&opts)
		},
	}
	registerReviewFlags(cmd, &opts)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return opts, err
}

// runSessionListCompat provides test compatibility for old-style runSessionList([]string{...}) calls.
func runSessionListCompat(args []string) error {
	cmd := &cobra.Command{
		Use:           "list",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, a []string) error {
			return runSessionList()
		},
	}
	cmd.Flags().StringVar(&sessionListRepoDir, "repo", "", "")
	cmd.Flags().BoolVar(&sessionListJSON, "json", false, "")
	cmd.Flags().IntVar(&sessionListLimit, "limit", 20, "")
	cmd.SetArgs(args)
	return cmd.Execute()
}

// runSessionShowCompat provides test compatibility for old-style runSessionShow([]string{...}) calls.
func runSessionShowCompat(args []string) error {
	cmd := &cobra.Command{
		Use:           "show",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, a []string) error {
			return runSessionShow(a[0])
		},
	}
	cmd.Flags().StringVar(&sessionShowRepoDir, "repo", "", "")
	cmd.Flags().BoolVar(&sessionShowJSON, "json", false, "")
	cmd.SetArgs(args)
	return cmd.Execute()
}

// runSessionCommentsCompat provides test compatibility for old-style runSessionComments([]string{...}) calls.
func runSessionCommentsCompat(args []string) error {
	cmd := &cobra.Command{
		Use:           "comments",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, a []string) error {
			return runSessionComments(a[0])
		},
	}
	cmd.Flags().StringVar(&sessionCommentsRepoDir, "repo", "", "")
	cmd.Flags().BoolVar(&sessionCommentsJSON, "json", false, "")
	cmd.Flags().StringVar(&sessionCommentsSeverity, "severity", "", "")
	cmd.Flags().StringVar(&sessionCommentsCategory, "category", "", "")
	cmd.SetArgs(args)
	return cmd.Execute()
}

// runSession provides test compatibility for the old dispatch function.
func runSession(args []string) error {
	cmd := &cobra.Command{Use: "session", SilenceUsage: true, SilenceErrors: true}
	listCmd := &cobra.Command{
		Use: "list", Aliases: []string{"ls"},
		SilenceUsage: true, SilenceErrors: true,
		RunE: func(cmd *cobra.Command, a []string) error { return runSessionList() },
	}
	listCmd.Flags().StringVar(&sessionListRepoDir, "repo", "", "")
	listCmd.Flags().BoolVar(&sessionListJSON, "json", false, "")
	listCmd.Flags().IntVar(&sessionListLimit, "limit", 20, "")
	showCmd := &cobra.Command{
		Use: "show", Args: cobra.ExactArgs(1),
		SilenceUsage: true, SilenceErrors: true,
		RunE: func(cmd *cobra.Command, a []string) error { return runSessionShow(a[0]) },
	}
	showCmd.Flags().StringVar(&sessionShowRepoDir, "repo", "", "")
	showCmd.Flags().BoolVar(&sessionShowJSON, "json", false, "")
	cmd.AddCommand(listCmd, showCmd)
	cmd.SetArgs(args)
	return cmd.Execute()
}

func executeReview(opts reviewOptions) error {
	return executeReviewContext(context.Background(), opts)
}
