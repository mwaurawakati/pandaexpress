package taobao

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	apiHost1688 = "1688-datahub.p.rapidapi.com"
)
func SearchImageOn1688(imageURL string, pageSize int) (*ApiResponse, error) {
	url := fmt.Sprintf("https://%s/item_search_image?imgUrl=%s&page=1&sort=default", apiHost, imageURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("x-rapidapi-key", apiKey)
	req.Header.Add("x-rapidapi-host", apiHost1688)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	var resp ApiResponse
	err = json.Unmarshal(body, &resp)
	return &resp, err
}

func SearchItemDetails1688(itemID string) (*Response, error) {
	url := fmt.Sprintf("https://%s/item_detail?itemId=%s", apiHost1688, itemID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("x-rapidapi-key", apiKey)
	req.Header.Add("x-rapidapi-host", apiHost1688)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	var resp Response
	err = json.Unmarshal(body, &resp)
	return &resp, err
}
