package warehouse

import (
	"context"
)

// ItemService Item Service
type ItemService interface {
	GetCounts(ctx context.Context) (map[string]int, error)
	Search(ctx context.Context, item Item) ([]*Item, error)
	Update(ctx context.Context, item Item) (string, error)
}
