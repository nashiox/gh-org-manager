package github

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

const APIv4Endpoint = "https://api.github.com/graphql"

var previewHeaders = []string{
	"application/vnd.github.ocelot-preview+json",
	"application/vnd.github.shadow-cat-preview+json",
}

type Client struct {
	endpoint *url.URL
	client   *http.Client
}

type apiRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type apiResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []apiError      `json:"errors,omitempty"`
}

type apiError struct {
	Message string `json:"message"`
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	endpoint, _ := url.Parse(APIv4Endpoint)
	return &Client{
		endpoint: endpoint,
		client:   httpClient,
	}
}

func (c *Client) request(ctx context.Context, query string, vars map[string]interface{}) ([]byte, error) {
	greq := apiRequest{
		Query:     query,
		Variables: vars,
	}

	data, err := json.Marshal(greq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.endpoint.String(), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	req.Header.Set("Accept", strings.Join(previewHeaders, ","))

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp apiResponse
	if err = json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if len(apiResp.Errors) > 0 {
		return nil, errors.New(apiResp.Errors[0].Message)
	}

	return []byte(apiResp.Data), nil
}
