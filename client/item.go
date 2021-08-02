package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/lbrooks/warehouse"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

type itemService struct {
	apiURL string
}

func NewItemService() warehouse.ItemService {
	apiURL := os.Getenv("WAREHOUSE_API_URL")
	if apiURL == "" {
		panic("WAREHOUSE_API_URL is undefined")
	}

	return &itemService{
		apiURL: apiURL,
	}
}

func (s *itemService) GetCounts(ctx context.Context) (map[string]int, error) {
	sc, span := warehouse.CreateSpan(ctx, "service", "count")
	defer span.End()

	req, err := http.NewRequest("GET", s.apiURL+"/api/item/count", nil)
	if err != nil {
		return nil, err
	}

	otel.GetTextMapPropagator().Inject(sc, propagation.HeaderCarrier(req.Header))

	span.AddEvent("Sent Request")
	resp, err := http.DefaultClient.Do(req)
	span.AddEvent("Completed Request")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var itemList map[string]int
	err = json.Unmarshal(body, &itemList)
	if err != nil {
		return nil, err
	}

	return itemList, nil
}

func (s *itemService) Search(ctx context.Context, item warehouse.Item) ([]*warehouse.Item, error) {
	sc, span := warehouse.CreateSpan(ctx, "service", "update")
	defer span.End()

	req, err := http.NewRequest("GET", s.apiURL+"/api/item", nil)
	if err != nil {
		return nil, err
	}

	otel.GetTextMapPropagator().Inject(sc, propagation.HeaderCarrier(req.Header))

	span.AddEvent("Sent Request")
	resp, err := http.DefaultClient.Do(req)
	span.AddEvent("Completed Request")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var itemList []*warehouse.Item
	err = json.Unmarshal(body, &itemList)
	if err != nil {
		return nil, err
	}

	return itemList, nil
}

func (s *itemService) Update(ctx context.Context, item warehouse.Item) (string, error) {
	sc, span := warehouse.CreateSpan(ctx, "service", "update")
	defer span.End()

	itemBytes, _ := json.Marshal(item)
	span.SetAttributes(attribute.String("item", string(itemBytes)))
	req, err := http.NewRequest("POST", s.apiURL+"/api/item", bytes.NewReader(itemBytes))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	otel.GetTextMapPropagator().Inject(sc, propagation.HeaderCarrier(req.Header))

	span.AddEvent("Sent Request")
	_, err = http.DefaultClient.Do(req)
	span.AddEvent("Completed Request")
	if err != nil {
		return "", err
	}

	return "", nil
}
