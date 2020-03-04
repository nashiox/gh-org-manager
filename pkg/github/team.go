package github

import (
	"context"
	"encoding/json"
)

type Team struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c *Client) ListTeams(ctx context.Context, org string) (teams []Team, err error) {
	query := `
		query listTemas($org: String!, $after: String) {
			organization(login: $org) {
				teams(first: 30, after: $after) {
					nodes {
						id
						name
					}
					pageInfo {
						hasNextPage
						endCursor
					}
				}
			}
		}
	`

	var resp struct {
		Organization struct {
			Teams struct {
				Nodes    []Team `json:"nodes"`
				PageInfo struct {
					HasNextPage bool   `json:"hasNextPage"`
					EndCursor   string `json:"endCursor"`
				} `json:"pageInfo"`
			} `json:"teams"`
		} `json:"organization"`
	}

	vars := map[string]interface{}{
		"org": org,
	}

	for {
		data, err := c.request(ctx, query, vars)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(data, &resp); err != nil {
			return nil, err
		}

		teams = append(teams, resp.Organization.Teams.Nodes...)

		if !resp.Organization.Teams.PageInfo.HasNextPage {
			break
		}

		vars["after"] = resp.Organization.Teams.PageInfo.EndCursor
	}

	return
}
