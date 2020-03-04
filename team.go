package main

import (
	"fmt"

	"github.com/google/go-github/v27/github"
	"github.com/spf13/cobra"
)

var cmdTeam = &cobra.Command{
	Use:   "team",
	Short: "manage team",
}

var cmdListTeam = &cobra.Command{
	Use:   "list",
	Short: "get team list",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		listOptions := &github.ListOptions{}

		var teams []*github.Team
		for {
			t, resp, err := GitHubClient.Teams.ListTeams(ctx, GitHubOrg, listOptions)
			if err != nil {
				return err
			}

			teams = append(teams, t...)

			if resp.NextPage == 0 {
				break
			}

			listOptions.Page = resp.NextPage
		}

		for _, t := range teams {
			if t.Name == nil {
				*t.Name = ""
			}

			if t.ID == nil {
				*t.ID = 0
			}

			if t.Description == nil {
				*t.Description = ""
			}

			fmt.Printf("Name: %s, ID: %d, Description: %s\n", *t.Name, *t.ID, *t.Description)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(cmdTeam)

	cmdTeam.AddCommand(cmdListTeam)
}
