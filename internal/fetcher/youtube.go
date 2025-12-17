package fetcher

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YoutubeFetcher struct {
	APIKey    string
	ChannelID string
	Start     time.Time
	End       time.Time
}

func (y *YoutubeFetcher) Name() string {
	return "YouTube"
}

func (y *YoutubeFetcher) FetchLinks(ctx context.Context) ([]string, error) {
	service, err := youtube.NewService(ctx, option.WithAPIKey(y.APIKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create youtube service: %w", err)
	}

	var links []string
	pageToken := ""

	for {
		call := service.Search.List([]string{"id"}).
			ChannelId(y.ChannelID).
			PublishedAfter(y.Start.Format(time.RFC3339)).
			PublishedBefore(y.End.Format(time.RFC3339)).
			Type("video").
			VideoDuration("short"). // filters for videos < 4 minutes
			MaxResults(50).
			PageToken(pageToken)

		response, err := call.Do()
		if err != nil {
			return nil, fmt.Errorf("youtube search failed: %w", err)
		}

		for _, item := range response.Items {
			if item.Id.VideoId == "" {
				continue
			}

			link := fmt.Sprintf("https://www.youtube.com/shorts/%s", item.Id.VideoId)
			links = append(links, link)
		}

		pageToken = response.NextPageToken
		if pageToken == "" {
			break
		}
	}

	return links, nil
}
