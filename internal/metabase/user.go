package metabase

import "fmt"

// ListUsers returns all users.
func (c *Client) ListUsers() ([]User, error) {
	var result struct {
		Data []User `json:"data"`
	}
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get("/api/user")
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetUser returns a user by ID.
func (c *Client) GetUser(id int) (*User, error) {
	var result User
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/user/%d", id))
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetCurrentUser returns the currently authenticated user.
func (c *Client) GetCurrentUser() (*User, error) {
	var result User
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get("/api/user/current")
	if err != nil {
		return nil, fmt.Errorf("get current user: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result, nil
}
