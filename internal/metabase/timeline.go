package metabase

import "fmt"

// ListTimelines returns all timelines.
func (c *Client) ListTimelines(collectionID *int) ([]Timeline, error) {
	var result []Timeline
	req := c.httpClient.R().SetResult(&result)
	if collectionID != nil {
		req.SetQueryParam("collection_id", fmt.Sprintf("%d", *collectionID))
	}
	resp, err := req.Get("/api/timeline")
	if err != nil {
		return nil, fmt.Errorf("list timelines: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return result, nil
}

// GetTimeline returns a timeline by ID.
func (c *Client) GetTimeline(id int) (*Timeline, error) {
	var result Timeline
	resp, err := c.httpClient.R().
		SetResult(&result).
		SetQueryParam("include", "events").
		Get(fmt.Sprintf("/api/timeline/%d", id))
	if err != nil {
		return nil, fmt.Errorf("get timeline: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result, nil
}
