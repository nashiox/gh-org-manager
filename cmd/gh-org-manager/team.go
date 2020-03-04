package cmd

import (
	"fmt"

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
		teams, err := GitHubClient.ListTeams(ctx, GitHubOrg)
		if err != nil {
			return err
		}

		for _, t := range teams {
			fmt.Printf("Name: %s, ID: %s, Description: %s\n", t.Name, t.ID, t.Description)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(cmdTeam)

	cmdTeam.AddCommand(cmdListTeam)
}
