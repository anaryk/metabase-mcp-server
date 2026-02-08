package metabase

import "fmt"

// ListActions returns actions for a model.
func (c *Client) ListActions(modelID int) ([]Action, error) {
	var result []Action
	resp, err := c.httpClient.R().
		SetResult(&result).
		SetQueryParam("model-id", fmt.Sprintf("%d", modelID)).
		Get("/api/action")
	if err != nil {
		return nil, fmt.Errorf("list actions: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return result, nil
}

// GetAction returns an action by ID.
func (c *Client) GetAction(id int) (*Action, error) {
	var result Action
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/action/%d", id))
	if err != nil {
		return nil, fmt.Errorf("get action: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result, nil
}
