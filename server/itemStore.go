package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lbrooks/warehouse"

	"go.opentelemetry.io/otel/attribute"
)

type ItemStore struct {
	items []*warehouse.Item
}

func (m *ItemStore) findItemsMatching(ctx context.Context, filter warehouse.Item) (items []*warehouse.Item) {
	_, span := warehouse.CreateSpan(ctx, "ItemStore", "findItemsMatching")
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

func (m *ItemStore) initialize(ctx context.Context) {
	sc, span := warehouse.CreateSpan(ctx, "ItemStore", "initialize")
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

// NewItemStore Create In Memory Storage
func NewItemStore(ctx context.Context, initalizeData bool) warehouse.ItemService {
	sc, span := warehouse.CreateSpan(ctx, "ItemStore", "NewItemStore")
	defer span.End()

	span.SetAttributes(attribute.Bool("initializeData", initalizeData))

	m := &ItemStore{
		items: make([]*warehouse.Item, 0),
	}
	span.AddEvent("Created Store")

	if initalizeData {
		m.initialize(sc)
		span.AddEvent("Data Injected")
	}

	return m
}

func (m *ItemStore) GetCounts(ctx context.Context) (map[string]int, error) {
	_, span := warehouse.CreateSpan(ctx, "ItemStore", "GetCounts")
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

func (m *ItemStore) Search(ctx context.Context, item warehouse.Item) (items []*warehouse.Item, err error) {
	sc, span := warehouse.CreateSpan(ctx, "ItemStore", "Search")
	defer span.End()

	items = m.findItemsMatching(sc, item)
	warehouse.SortItems(items)

	return
}

func (m *ItemStore) Update(ctx context.Context, item warehouse.Item) (string, error) {
	sc, span := warehouse.CreateSpan(ctx, "ItemStore", "Update")
	defer span.End()

	j, _ := json.Marshal(item)
	span.SetAttributes(attribute.String("item", string(j)))

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
