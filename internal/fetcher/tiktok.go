package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type TikTokFetcher struct {
	RapidAPIKey  string
	RapidAPIHost string
	Username     string
	Start        time.Time
	End          time.Time
}

type tiktokPostsResponse struct {
	Data struct {
		Cursor   string `json:"cursor"`
		HasMore  bool   `json:"hasMore"`
		ItemList []struct {
			ID         string `json:"id"`
			CreateTime int64  `json:"createTime"`
		} `json:"itemList"`
	} `json:"data"`
}

type tiktokInfoResponse struct {
	UserInfo struct {
		User struct {
			SecUID string `json:"secUid"`
			ID     string `json:"id"`
		} `json:"user"`
	} `json:"userInfo"`
}

func (t *TikTokFetcher) Name() string {
	return "TikTok"
}

func (t *TikTokFetcher) resolveSecUID(ctx context.Context) (string, error) {
	fmt.Printf("Resolving SecUID for user '%s'...\n", t.Username)

	url := fmt.Sprintf("https://%s/api/user/info?uniqueId=%s", t.RapidAPIHost, t.Username)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("x-rapidapi-host", t.RapidAPIHost)
	req.Header.Add("x-rapidapi-key", t.RapidAPIKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("user info API returned status: %d", res.StatusCode)
	}

	var result tiktokInfoResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode user info: %w", err)
	}

	uid := result.UserInfo.User.SecUID
	if uid == "" {
		return "", fmt.Errorf("could not find secUid for user: %s", t.Username)
	}

	fmt.Println("Found SecUID:", uid)
	return uid, nil
}

func (t *TikTokFetcher) FetchLinks(ctx context.Context) ([]string, error) {
	// get the SecUID
	secUID, err := t.resolveSecUID(ctx)
	if err != nil {
		return nil, err
	}

	var links []string
	cursor := "0"
	hasMore := true

	for hasMore {
		url := fmt.Sprintf("https://%s/api/user/posts?secUid=%s&count=35&cursor=%s",
			t.RapidAPIHost, secUID, cursor)

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Add("x-rapidapi-host", t.RapidAPIHost)
		req.Header.Add("x-rapidapi-key", t.RapidAPIKey)

		client := &http.Client{Timeout: 10 * time.Second}
		res, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			return nil, fmt.Errorf("tiktok api returned status: %d", res.StatusCode)
		}

		var result tiktokPostsResponse
		if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode json: %w", err)
		}

		if len(result.Data.ItemList) == 0 {
			break
		}

		for _, v := range result.Data.ItemList {
			ts := time.Unix(v.CreateTime, 0)

			if ts.After(t.End) {
				continue
			}
			if ts.Before(t.Start) {
				hasMore = false
				break
			}

			if v.ID != "" {
				link := fmt.Sprintf("https://www.tiktok.com/@%s/video/%s", t.Username, v.ID)
				links = append(links, link)
			}
		}

		cursor = result.Data.Cursor
		if !result.Data.HasMore || cursor == "0" || cursor == "-1" || cursor == "" {
			hasMore = false
		}
	}

	return links, nil
}
