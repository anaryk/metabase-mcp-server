package metabase

import "fmt"

// ListTables returns all tables for a given database.
func (c *Client) ListTables(databaseID int) ([]Table, error) {
	var result []Table
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/database/%d/metadata/tables", databaseID))
	if err != nil {
		return nil, fmt.Errorf("list tables: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return result, nil
}

// GetTable returns a table by ID.
func (c *Client) GetTable(id int) (*Table, error) {
	var result Table
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/table/%d", id))
	if err != nil {
		return nil, fmt.Errorf("get table: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetTableMetadata returns table metadata including all fields.
func (c *Client) GetTableMetadata(id int) (*Table, error) {
	var result Table
	resp, err := c.httpClient.R().
		SetResult(&result).
		SetQueryParam("include_hidden_fields", "true").
		Get(fmt.Sprintf("/api/table/%d/query_metadata", id))
	if err != nil {
		return nil, fmt.Errorf("get table metadata: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetTableForeignKeys returns foreign key relationships for a table.
func (c *Client) GetTableForeignKeys(id int) ([]ForeignKey, error) {
	var result []ForeignKey
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/table/%d/fks", id))
	if err != nil {
		return nil, fmt.Errorf("get table fks: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return result, nil
}
