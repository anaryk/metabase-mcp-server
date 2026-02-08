package metabase

import "fmt"

// ListPermissionGroups returns all permission groups.
func (c *Client) ListPermissionGroups() ([]PermissionGroup, error) {
	var result []PermissionGroup
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get("/api/permissions/group")
	if err != nil {
		return nil, fmt.Errorf("list permission groups: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return result, nil
}

// GetPermissionGroup returns a permission group by ID.
func (c *Client) GetPermissionGroup(id int) (*PermissionGroup, error) {
	var result PermissionGroup
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/permissions/group/%d", id))
	if err != nil {
		return nil, fmt.Errorf("get permission group: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPermissionsGraph returns the full permissions graph.
func (c *Client) GetPermissionsGraph() (map[string]any, error) {
	var result map[string]any
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get("/api/permissions/graph")
	if err != nil {
		return nil, fmt.Errorf("get permissions graph: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return result, nil
}
