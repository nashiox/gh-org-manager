package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/google/go-github/v27/github"
	"github.com/spf13/cobra"
)

var GitHubRoles = []string{
	"admin",
	"direct_member",
	"billing_manager",
}

var cmdMember = &cobra.Command{
	Use:   "member",
	Short: "manage member",
}

var cmdAddMember = &cobra.Command{
	Use:   "add",
	Short: "create invitation for add member",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}

		email, err := cmd.Flags().GetString("email")
		if err != nil {
			return err
		}

		if (name == "" && email == "") || (name != "" && email != "") {
			return errors.New("set either email or name")
		}

		role, err := cmd.Flags().GetString("role")
		if err != nil {
			return err
		}

		if !func(r string) bool {
			for _, v := range GitHubRoles {
				if r == v {
					return true
				}
			}
			return false
		}(role) {
			return errors.New("require ['admin', 'direct_member', 'billing_manager']")
		}

		teams, err := cmd.Flags().GetStringSlice("team")
		if err != nil {
			return err
		}

		listOptions := &github.ListOptions{}

		orgTeams := make([]*github.Team, 0)
		for {
			t, resp, err := GitHubClient.Teams.ListTeams(ctx, GitHubOrg, listOptions)
			if err != nil {
				return err
			}

			orgTeams = append(orgTeams, t...)

			if resp.NextPage == 0 {
				break
			}

			listOptions.Page = resp.NextPage
		}

		teamIDs := make([]int64, 0)
		for _, t := range orgTeams {
			for _, tn := range teams {
				if *t.Name == tn {
					teamIDs = append(teamIDs, *t.ID)
				}
			}
		}

		invitationOpts := &github.CreateOrgInvitationOptions{
			Role:   &role,
			TeamID: teamIDs,
		}

		if name != "" {
			user, _, err := GitHubClient.Users.Get(ctx, name)
			if err != nil {
				return err
			}

			invitationOpts.InviteeID = user.ID
		} else if email != "" {
			invitationOpts.Email = &email
		}

		invitation, _, err := GitHubClient.Organizations.CreateOrgInvitation(ctx, GitHubOrg, invitationOpts)
		if err != nil {
			return err
		}

		fmt.Printf("Login: %#v, ID: %#v, Email: %#v, Role: %#v\n", *invitation.Login, *invitation.ID, *invitation.Email, *invitation.Role)

		return nil
	},
}

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
				if _, err := GitHubClient.Organizations.RemoveMember(ctx, GitHubOrg, user); err != nil {
					return err
				}
			}
		case "id":
			for _, arg := range args {
				id, err := strconv.ParseInt(arg, 10, 64)
				if err != nil {
					return err
				}

				user, _, err := GitHubClient.Users.GetByID(ctx, id)
				if err != nil {
					return err
				}

				if _, err = GitHubClient.Organizations.RemoveMember(ctx, GitHubOrg, *user.Login); err != nil {
					return err
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
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		membersOpts := &github.ListMembersOptions{}

		members := make([]*github.User, 0)
		for {
			users, resp, err := GitHubClient.Organizations.ListMembers(ctx, GitHubOrg, membersOpts)
			if err != nil {
				return err
			}

			members = append(members, users...)

			if resp.NextPage == 0 {
				break
			}

			membersOpts.Page = resp.NextPage
		}

		for _, m := range members {
			fmt.Printf("Login: %#v, ID: %#v\n", *m.Login, *m.ID)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(cmdMember)

	cmdMember.AddCommand(cmdAddMember)
	cmdAddMember.Flags().StringP("name", "n", "", "Required unless you provide email. GitHub user name for the person you are inviting.")
	cmdAddMember.Flags().String("email", "", "Required unless you provide name. Email address of the person you are inviting, which can be an existing GitHub user.")
	cmdAddMember.Flags().StringP("role", "r", "direct_member", "Specify role for new member ['admin', 'direct_member', 'billing_manager']")
	cmdAddMember.Flags().StringSliceP("team", "t", []string{}, "Specify names for the teams you want to invite new members to.")

	cmdMember.AddCommand(cmdDeleteMember)
	cmdDeleteMember.Flags().StringP("mode", "m", "name", "delete by ['name' , 'id']")

	cmdMember.AddCommand(cmdListMember)
}
