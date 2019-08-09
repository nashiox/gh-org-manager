package main

import (
	"errors"
	"fmt"

	"github.com/google/go-github/v27/github"
	"github.com/spf13/cobra"
)

var githubMembersOpts *github.ListMembersOptions

var cmdMember = &cobra.Command{
	Use:   "member",
	Short: "manage member",
}

var cmdAddMember = &cobra.Command{}

var cmdDeleteMember = &cobra.Command{
	Use:   "delete",
	Short: "delete member",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mode, err := cmd.Flags().GetString("mode")
		if err != nil {
			return err
		}

		switch mode {
		case "name":
			for _, user := range args {
				if _, err := GitHubClient.Organizations.RemoveMember(ctx, user, GitHubOrg); err != nil {
					return nil
				}
			}
		case "id":
			for _, id := range args {
				user, err := GitHubClient.Users.GetByID(ctx, id)
				if err != nil {
					return nil
				}

				if _, err = GitHubClient.Organizations.RemoveMember(ctx, *user.Login, GitHubOrg); err != nil {
					return nil
				}
			}
		default:
			return errors.New("expect mode 'name' or 'id'")
		}

		return nil
	},
}

var cmdListMember = &cobra.Command{
	Use:   "list",
	Short: "get member list",
	RunE: func(cmd *cobra.Command, args []string) error {
		var members []*github.User
		for {
			users, resp, err := GitHubClient.Organizations.ListMembers(ctx, GitHubOrg, githubMembersOpts)
			if err != nil {
				return err
			}

			members = append(members, users...)

			if resp.NextPage == 0 {
				break
			}

			githubMembersOpts.Page = resp.NextPage
		}

		for _, m := range members {
			fmt.Printf("Login: %s, ID: %d\n", *m.Login, *m.ID)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(cmdMember)

	cmdMember.AddCommand(cmdAddMember)

	cmdMember.AddCommand(cmdDeleteMember)
	cmdDeleteMember.Flags().StringP("mode", "m", "name", "delete by ['name' , 'id]")

	cmdMember.AddCommand(cmdListMember)

	githubMembersOpts = &github.ListMembersOptions{
		ListOptions: github.ListOptions{PerPage: 30},
	}
}
