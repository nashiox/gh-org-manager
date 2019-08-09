package main

import (
	"context"
	"os"

	"github.com/google/go-github/v27/github"
	"github.com/spf13/cobra"

	"golang.org/x/oauth2"
)

var rootCmd = &cobra.Command{
	Use:   "ghom",
	Short: "ghom is GitHub Organization management cli.",
	Long:  "ghom is GitHub Organization management cli.",
}

var GitHubOrg string

var GitHubClient *github.Client
var ctx context.Context

func init() {
	rootCmd.Version = version

	rootCmd.PersistentFlags().StringVarP(&GitHubOrg, "organization", "o", "", "GitHub Organization name")

	ctx = context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	GitHubClient = github.NewClient(tc)
}

func execute() error {
	return rootCmd.Execute()
}
