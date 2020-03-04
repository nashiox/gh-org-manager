package github

import (
	"context"
	"encoding/json"
)

type Member struct {
	ID    uint64 `json:"databaseID"`
	Login string `json:"login"`
}

func (c *Client) ListMembers(ctx context.Context, org string) (members []Member, err error) {
	query := `
		query listMembers($org: String!, $after: String) {
			organization(login: $org) {
				membersWithRole(first: 30, after: $after) {
					nodes {
						databaseId
						login
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
			Members struct {
				Nodes    []Member `json:"nodes"`
				PageInfo struct {
					HasNextPage bool   `json:"hasNextPage"`
					EndCursor   string `json:"endCursor"`
				} `json:"pageInfo"`
			} `json:"membersWithRole"`
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

		members = append(members, resp.Organization.Members.Nodes...)

		if !resp.Organization.Members.PageInfo.HasNextPage {
			break
		}

		vars["after"] = resp.Organization.Members.PageInfo.EndCursor
	}

	return
}
