package metabase

import "fmt"

// Search searches across all entities.
func (c *Client) Search(query string, models []string) (*SearchResponse, error) {
	var result SearchResponse
	req := c.httpClient.R().
		SetResult(&result).
		SetQueryParam("q", query)
	for _, m := range models {
		req.SetQueryParam("models", m)
	}
	resp, err := req.Get("/api/search")
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result, nil
}
