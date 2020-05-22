package turnipexchange

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type IslandsRequest struct {
	Islander string `json:"islander"`
	Fee      int    `json:"fee"`
	Category string `json:"category"`
}

type IslandsResponse struct {
	Success bool
	Message string
	Islands []Island
}

func (c *Client) Islands(Islander string, Category string, Fee int) ([]Island, *http.Response, error) {
	url := fmt.Sprintf("%s/islands/", c.BaseURL)
	payload := &IslandsRequest{
		Islander: Islander,
		Category: Category,
		Fee:      Fee,
	}
	req, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, err
	}

	resp, err := http.Post(url, c.ContentType, bytes.NewBuffer(req))
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()

	var data IslandsResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, nil, err
	}

	if !data.Success {
		return nil, nil, &ErrorResponseNotSuccess{Success: data.Success, Message: data.Message}
	}

	return data.Islands, resp, nil
}
