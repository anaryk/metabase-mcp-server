package metabase

import "fmt"

// GetField returns field details by ID.
func (c *Client) GetField(id int) (*Field, error) {
	var result Field
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/field/%d", id))
	if err != nil {
		return nil, fmt.Errorf("get field: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetFieldValues returns distinct values for a field.
func (c *Client) GetFieldValues(id int) (*FieldValues, error) {
	var result FieldValues
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/field/%d/values", id))
	if err != nil {
		return nil, fmt.Errorf("get field values: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result, nil
}

// SearchFieldValues searches distinct values for a field by prefix.
func (c *Client) SearchFieldValues(id int, query string, limit int) (*FieldValues, error) {
	var result FieldValues
	req := c.httpClient.R().
		SetResult(&result).
		SetQueryParam("value", query)
	if limit > 0 {
		req.SetQueryParam("limit", fmt.Sprintf("%d", limit))
	}
	resp, err := req.Get(fmt.Sprintf("/api/field/%d/search/%d", id, id))
	if err != nil {
		return nil, fmt.Errorf("search field values: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result, nil
}
