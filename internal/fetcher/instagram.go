package fetcher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type InstagramFetcher struct {
	RapidAPIKey  string
	RapidAPIHost string
	Username     string
	Start        time.Time
	End          time.Time
}

type instagramReelsResponse struct {
	Result struct {
		Edges []struct {
			Node struct {
				Media struct {
					Code string `json:"code"`
					Pk   string `json:"pk"`
				} `json:"media"`
			} `json:"node"`
			Cursor string `json:"cursor"`
		} `json:"edges"`
	} `json:"result"`
}

func (i *InstagramFetcher) Name() string {
	return "Instagram"
}

func (i *InstagramFetcher) FetchLinks(ctx context.Context) ([]string, error) {
	var links []string
	maxID := ""

	for {
		url := fmt.Sprintf("https://%s/api/instagram/reels", i.RapidAPIHost)

		payload := map[string]string{
			"username": i.Username,
			"maxId":    maxID,
		}

		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}

		req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonPayload))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Add("x-rapidapi-key", i.RapidAPIKey)
		req.Header.Add("x-rapidapi-host", i.RapidAPIHost)
		req.Header.Add("Content-Type", "application/json")

		client := &http.Client{Timeout: 10 * time.Second}
		res, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to make request: %w", err)
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			return nil, fmt.Errorf("instagram API returned status: %d", res.StatusCode)
		}

		var result instagramReelsResponse
		if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode json: %w", err)
		}

		if len(result.Result.Edges) == 0 {
			break
		}

		for _, edge := range result.Result.Edges {
			if edge.Node.Media.Code != "" {
				// Instagram links can use /p/ for both posts and reels
				link := fmt.Sprintf("https://www.instagram.com/p/%s/", edge.Node.Media.Code)
				links = append(links, link)
			}
		}

		// Check to continue pagination
		// If the last cursor is empty or same as maxID, don't continue
		lastCursor := result.Result.Edges[len(result.Result.Edges)-1].Cursor
		if lastCursor == "" || lastCursor == maxID {
			break
		}

		maxID = lastCursor
	}

	return links, nil
}
