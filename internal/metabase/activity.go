package metabase

import "fmt"

// GetActivity returns the recent activity log.
func (c *Client) GetActivity() ([]ActivityItem, error) {
	var result []ActivityItem
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get("/api/activity")
	if err != nil {
		return nil, fmt.Errorf("get activity: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return result, nil
}

// GetRecentViews returns recently viewed items.
func (c *Client) GetRecentViews() ([]RecentItem, error) {
	var result []RecentItem
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get("/api/activity/recent_views")
	if err != nil {
		return nil, fmt.Errorf("get recent views: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return result, nil
}
