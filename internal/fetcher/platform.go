package fetcher

import "context"

type Platform interface {
	Name() string
	FetchLinks(ctx context.Context) ([]string, error)
}
