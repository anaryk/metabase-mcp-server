package metabase

import "fmt"

// ExecuteQuery executes a dataset query (native SQL or MBQL).
func (c *Client) ExecuteQuery(req *DatasetQueryRequest) (*DatasetQueryResponse, error) {
	var result DatasetQueryResponse
	resp, err := c.httpClient.R().
		SetBody(req).
		SetResult(&result).
		Post("/api/dataset")
	if err != nil {
		return nil, fmt.Errorf("execute query: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result, nil
}

// ExportQueryResults exports query results in the given format (csv, json, xlsx).
func (c *Client) ExportQueryResults(req *DatasetQueryRequest, format string) ([]byte, error) {
	resp, err := c.httpClient.R().
		SetBody(req).
		Post(fmt.Sprintf("/api/dataset/%s", format))
	if err != nil {
		return nil, fmt.Errorf("export query results: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return resp.Body(), nil
}
