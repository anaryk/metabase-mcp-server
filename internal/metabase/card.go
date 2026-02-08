package metabase

import "fmt"

// ListCards returns all saved questions/cards.
func (c *Client) ListCards() ([]Card, error) {
	var result []Card
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get("/api/card")
	if err != nil {
		return nil, fmt.Errorf("list cards: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return result, nil
}

// GetCard returns a card by ID.
func (c *Client) GetCard(id int) (*Card, error) {
	var result Card
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/card/%d", id))
	if err != nil {
		return nil, fmt.Errorf("get card: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateCard creates a new saved question/card.
func (c *Client) CreateCard(card *Card) (*Card, error) {
	var result Card
	resp, err := c.httpClient.R().
		SetBody(card).
		SetResult(&result).
		Post("/api/card")
	if err != nil {
		return nil, fmt.Errorf("create card: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateCard updates an existing card.
func (c *Client) UpdateCard(id int, card *Card) (*Card, error) {
	var result Card
	resp, err := c.httpClient.R().
		SetBody(card).
		SetResult(&result).
		Put(fmt.Sprintf("/api/card/%d", id))
	if err != nil {
		return nil, fmt.Errorf("update card: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteCard archives/deletes a card.
func (c *Client) DeleteCard(id int) error {
	resp, err := c.httpClient.R().
		Delete(fmt.Sprintf("/api/card/%d", id))
	if err != nil {
		return fmt.Errorf("delete card: %w", err)
	}
	return checkResponse(resp)
}

// ExecuteCardQuery runs a card's saved query and returns results.
func (c *Client) ExecuteCardQuery(id int, parameters map[string]any) (*DatasetQueryResponse, error) {
	var result DatasetQueryResponse
	req := c.httpClient.R().SetResult(&result)
	if len(parameters) > 0 {
		req.SetBody(parameters)
	}
	resp, err := req.Post(fmt.Sprintf("/api/card/%d/query", id))
	if err != nil {
		return nil, fmt.Errorf("execute card query: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result, nil
}
