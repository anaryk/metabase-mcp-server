package metabase

import "fmt"

// ListAlerts returns all alerts.
func (c *Client) ListAlerts() ([]Alert, error) {
	var result []Alert
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get("/api/alert")
	if err != nil {
		return nil, fmt.Errorf("list alerts: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return result, nil
}

// GetAlert returns an alert by ID.
func (c *Client) GetAlert(id int) (*Alert, error) {
	var result Alert
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/alert/%d", id))
	if err != nil {
		return nil, fmt.Errorf("get alert: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateAlert creates a new alert.
func (c *Client) CreateAlert(alert *Alert) (*Alert, error) {
	var result Alert
	resp, err := c.httpClient.R().
		SetBody(alert).
		SetResult(&result).
		Post("/api/alert")
	if err != nil {
		return nil, fmt.Errorf("create alert: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result, nil
}
