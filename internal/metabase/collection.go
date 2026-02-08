package metabase

import "fmt"

// ListCollections returns all collections.
func (c *Client) ListCollections(namespace string) ([]Collection, error) {
	var result []Collection
	req := c.httpClient.R().SetResult(&result)
	if namespace != "" {
		req.SetQueryParam("namespace", namespace)
	}
	resp, err := req.Get("/api/collection")
	if err != nil {
		return nil, fmt.Errorf("list collections: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return result, nil
}

// GetCollection returns a collection by ID.
func (c *Client) GetCollection(id string) (*Collection, error) {
	var result Collection
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/collection/%s", id))
	if err != nil {
		return nil, fmt.Errorf("get collection: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateCollection creates a new collection.
func (c *Client) CreateCollection(collection *Collection) (*Collection, error) {
	var result Collection
	resp, err := c.httpClient.R().
		SetBody(collection).
		SetResult(&result).
		Post("/api/collection")
	if err != nil {
		return nil, fmt.Errorf("create collection: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateCollection updates an existing collection.
func (c *Client) UpdateCollection(id int, collection *Collection) (*Collection, error) {
	var result Collection
	resp, err := c.httpClient.R().
		SetBody(collection).
		SetResult(&result).
		Put(fmt.Sprintf("/api/collection/%d", id))
	if err != nil {
		return nil, fmt.Errorf("update collection: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListCollectionItems returns items in a collection.
func (c *Client) ListCollectionItems(id string, models []string) ([]CollectionItem, error) {
	var result struct {
		Data []CollectionItem `json:"data"`
	}
	req := c.httpClient.R().SetResult(&result)
	for _, m := range models {
		req.SetQueryParam("models", m)
	}
	resp, err := req.Get(fmt.Sprintf("/api/collection/%s/items", id))
	if err != nil {
		return nil, fmt.Errorf("list collection items: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return result.Data, nil
}
