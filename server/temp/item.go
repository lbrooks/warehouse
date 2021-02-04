package temp

import (
	"context"
	"fmt"

	"github.com/lbrooks/warehouse"
	"github.com/lbrooks/warehouse/server"

	"go.opentelemetry.io/otel/label"
)

type itemService struct {
	items []*warehouse.Item
}

func (m *itemService) findItemsMatching(ctx context.Context, filter warehouse.Item) (items []*warehouse.Item) {
	_, span := warehouse.CreateSpan(ctx, "itemService", "findItemsMatching")
	defer span.End()

	ignoreEmpty := filter.Barcode == "" && filter.Name == "" && filter.Brand == ""

	items = make([]*warehouse.Item, 0)
	for _, v := range m.items {
		if filter.Matches(v) {
			if !ignoreEmpty || v.Quantity > 0 {
				items = append(items, v)
			}
		}
	}
	return
}

func (m *itemService) initialize(ctx context.Context) {
	sc, span := warehouse.CreateSpan(ctx, "itemService", "initialize")
	defer span.End()

	_, err := m.Update(sc, warehouse.Item{Barcode: "1", Name: "Toilet Paper", Brand: "Charmin", Quantity: 1})
	if err != nil {
		span.RecordError(err)
		return
	}
	_, err = m.Update(sc, warehouse.Item{Barcode: "2", Name: "Toilet Paper", Brand: "Sandpaper", Quantity: 1})
	if err != nil {
		span.RecordError(err)
		return
	}
	_, err = m.Update(sc, warehouse.Item{Barcode: "1", Name: "Paper Towels", Brand: "Bounty", Quantity: 1})
	if err != nil {
		span.RecordError(err)
		return
	}
	_, err = m.Update(sc, warehouse.Item{Barcode: "3", Name: "Gallon Bag", Brand: "Ziploc", Quantity: 1})
	if err != nil {
		span.RecordError(err)
		return
	}
	_, err = m.Update(sc, warehouse.Item{Barcode: "4", Name: "Quart Bag", Brand: "Ziploc", Quantity: 3})
	if err != nil {
		span.RecordError(err)
		return
	}
	span.AddEvent("Mock Data Initalized")
}

// NewItemService Create In Memory Storage
func NewItemService(ctx context.Context, initalizeData bool) server.ItemService {
	sc, span := warehouse.CreateSpan(ctx, "itemService", "NewItemService")
	defer span.End()

	span.SetAttributes(label.Bool("initializeData", initalizeData))

	m := &itemService{
		items: make([]*warehouse.Item, 0),
	}
	span.AddEvent("Created Store")

	if initalizeData {
		m.initialize(sc)
		span.AddEvent("Data Injected")
	}

	return m
}

func (m *itemService) GetCounts(ctx context.Context) (map[string]int, error) {
	_, span := warehouse.CreateSpan(ctx, "itemService", "GetCounts")
	defer span.End()

	counts := make(map[string]int)
	for _, i := range m.items {
		if _, found := counts[i.Name]; !found {
			counts[i.Name] = i.Quantity
		} else {
			counts[i.Name] += i.Quantity
		}
	}
	return counts, nil
}

func (m *itemService) Search(ctx context.Context, item warehouse.Item) (items []*warehouse.Item, err error) {
	sc, span := warehouse.CreateSpan(ctx, "itemService", "Search")
	defer span.End()

	items = m.findItemsMatching(sc, item)
	warehouse.SortItems(items)

	return
}

func (m *itemService) Update(ctx context.Context, item warehouse.Item) (string, error) {
	sc, span := warehouse.CreateSpan(ctx, "itemService", "Update")
	defer span.End()

	span.SetAttributes(
		label.String("barcode", item.Barcode),
		label.String("name", item.Name),
		label.Int("quantity", item.Quantity),
	)

	matching := m.findItemsMatching(sc, item)
	if len(matching) == 0 {
		m.items = append(m.items, &item)
	} else if len(matching) > 1 {
		err := fmt.Errorf("Multiple Items Matched")
		span.RecordError(err)
		return "", err
	} else {
		matching[0].Quantity = item.Quantity
		if matching[0].Quantity < 0 {
			matching[0].Quantity = 0
		}
	}

	return "", nil
}
