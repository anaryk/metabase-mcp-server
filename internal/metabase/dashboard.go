package metabase

import "fmt"

// ListDashboards returns all dashboards.
func (c *Client) ListDashboards() ([]Dashboard, error) {
	var result []Dashboard
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get("/api/dashboard")
	if err != nil {
		return nil, fmt.Errorf("list dashboards: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return result, nil
}

// GetDashboard returns a dashboard by ID.
func (c *Client) GetDashboard(id int) (*Dashboard, error) {
	var result Dashboard
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/dashboard/%d", id))
	if err != nil {
		return nil, fmt.Errorf("get dashboard: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateDashboard creates a new dashboard.
func (c *Client) CreateDashboard(dashboard *Dashboard) (*Dashboard, error) {
	var result Dashboard
	resp, err := c.httpClient.R().
		SetBody(dashboard).
		SetResult(&result).
		Post("/api/dashboard")
	if err != nil {
		return nil, fmt.Errorf("create dashboard: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateDashboard updates an existing dashboard.
func (c *Client) UpdateDashboard(id int, dashboard *Dashboard) (*Dashboard, error) {
	var result Dashboard
	resp, err := c.httpClient.R().
		SetBody(dashboard).
		SetResult(&result).
		Put(fmt.Sprintf("/api/dashboard/%d", id))
	if err != nil {
		return nil, fmt.Errorf("update dashboard: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteDashboard deletes a dashboard.
func (c *Client) DeleteDashboard(id int) error {
	resp, err := c.httpClient.R().
		Delete(fmt.Sprintf("/api/dashboard/%d", id))
	if err != nil {
		return fmt.Errorf("delete dashboard: %w", err)
	}
	return checkResponse(resp)
}

// AddCardToDashboard adds a card to a dashboard.
func (c *Client) AddCardToDashboard(dashboardID int, dashCard *DashCard) (*DashCard, error) {
	var result DashCard
	resp, err := c.httpClient.R().
		SetBody(dashCard).
		SetResult(&result).
		Post(fmt.Sprintf("/api/dashboard/%d/cards", dashboardID))
	if err != nil {
		return nil, fmt.Errorf("add card to dashboard: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result, nil
}

// RemoveCardFromDashboard removes a dashcard from a dashboard.
func (c *Client) RemoveCardFromDashboard(dashboardID, dashCardID int) error {
	resp, err := c.httpClient.R().
		Delete(fmt.Sprintf("/api/dashboard/%d/cards?dashcardId=%d", dashboardID, dashCardID))
	if err != nil {
		return fmt.Errorf("remove card from dashboard: %w", err)
	}
	return checkResponse(resp)
}

// UpdateDashboardCards updates the layout/positions of cards on a dashboard.
func (c *Client) UpdateDashboardCards(dashboardID int, cards []DashCard) error {
	body := map[string]any{"cards": cards}
	resp, err := c.httpClient.R().
		SetBody(body).
		Put(fmt.Sprintf("/api/dashboard/%d/cards", dashboardID))
	if err != nil {
		return fmt.Errorf("update dashboard cards: %w", err)
	}
	return checkResponse(resp)
}

// CopyDashboard copies a dashboard to a new collection.
func (c *Client) CopyDashboard(id int, name string, description *string, collectionID *int) (*Dashboard, error) {
	body := map[string]any{"name": name}
	if description != nil {
		body["description"] = *description
	}
	if collectionID != nil {
		body["collection_id"] = *collectionID
	}
	var result Dashboard
	resp, err := c.httpClient.R().
		SetBody(body).
		SetResult(&result).
		Post(fmt.Sprintf("/api/dashboard/%d/copy", id))
	if err != nil {
		return nil, fmt.Errorf("copy dashboard: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result, nil
}
