package metabase

import "fmt"

// InvalidateCache invalidates the Metabase cache.
func (c *Client) InvalidateCache() error {
	resp, err := c.httpClient.R().
		Post("/api/cache/invalidate")
	if err != nil {
		return fmt.Errorf("invalidate cache: %w", err)
	}
	return checkResponse(resp)
}
