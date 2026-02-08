package metabase

import "fmt"

// ListSettings returns all Metabase settings.
func (c *Client) ListSettings() ([]Setting, error) {
	var result []Setting
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get("/api/setting")
	if err != nil {
		return nil, fmt.Errorf("list settings: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return result, nil
}

// GetSetting returns a specific setting by key.
func (c *Client) GetSetting(key string) (any, error) {
	var result any
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/setting/%s", key))
	if err != nil {
		return nil, fmt.Errorf("get setting: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return result, nil
}
