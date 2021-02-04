package server

import (
	"context"

	"github.com/lbrooks/warehouse"
)

// ItemService Item Service
type ItemService interface {
	GetCounts(ctx context.Context) (map[string]int, error)
	Search(ctx context.Context, item warehouse.Item) ([]*warehouse.Item, error)
	Update(ctx context.Context, item warehouse.Item) (string, error)
}
