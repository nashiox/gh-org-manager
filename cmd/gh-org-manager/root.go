package cmd

import (
	"context"
	"os"

	"github.com/nashiox/gh-org-manager/pkg/github"
	"github.com/spf13/cobra"

	"golang.org/x/oauth2"
)

var rootCmd = &cobra.Command{
	Use:   "ghom",
	Short: "ghom is GitHub Organization management cli.",
}

var ctx context.Context

var GitHubOrg string
var GitHubClient *github.Client

func GetRootCmd(version string) *cobra.Command {
	rootCmd.Version = version

	rootCmd.PersistentFlags().StringVarP(&GitHubOrg, "organization", "o", "", "GitHub Organization name")

	ctx = context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	GitHubClient = github.NewClient(tc)

	return rootCmd
}
