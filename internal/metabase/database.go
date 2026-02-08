package metabase

import "fmt"

// ListDatabases returns all connected databases.
func (c *Client) ListDatabases() ([]Database, error) {
	var result struct {
		Data []Database `json:"data"`
	}
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get("/api/database")
	if err != nil {
		return nil, fmt.Errorf("list databases: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetDatabase returns a database by ID.
func (c *Client) GetDatabase(id int) (*Database, error) {
	var result Database
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/database/%d", id))
	if err != nil {
		return nil, fmt.Errorf("get database: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetDatabaseMetadata returns full metadata for a database.
func (c *Client) GetDatabaseMetadata(id int) (*Database, error) {
	var result Database
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/database/%d/metadata", id))
	if err != nil {
		return nil, fmt.Errorf("get database metadata: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result, nil
}

// SyncDatabase triggers a schema sync for a database.
func (c *Client) SyncDatabase(id int) error {
	resp, err := c.httpClient.R().
		Post(fmt.Sprintf("/api/database/%d/sync_schema", id))
	if err != nil {
		return fmt.Errorf("sync database: %w", err)
	}
	return checkResponse(resp)
}
