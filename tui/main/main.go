package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/lbrooks/warehouse"
	"github.com/rivo/tview"
	"go.opentelemetry.io/otel"
)

var apiUrl string

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	apiUrl = os.Getenv("WAREHOUSE_API_URL")
	if apiUrl == "" {
		panic("WAREHOUSE_URL is undefined")
	}
}

func getData() []warehouse.Item {
	sc, span := warehouse.CreateSpan(context.TODO(), "http-req", "fetch-items")
	defer span.End()

	req, err := http.NewRequest("GET", apiUrl+"/api/item", nil)
	if err != nil {
		panic(err.Error())
	}

	otel.GetTextMapPropagator().Inject(sc, req.Header)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var itemList []warehouse.Item
	err = json.Unmarshal(body, &itemList)
	if err != nil {
		panic(err.Error())
	}

	return itemList
}

func main() {
	flush := warehouse.InitializeJaeger("warehouse-tui")
	defer flush()

	items := getData()

	table := tview.NewTable()
	table.SetFixed(1, 0)
	table.SetSelectable(true, false)
	table.SetCell(0, 0, tview.NewTableCell("Barcode").SetSelectable(false))
	table.SetCell(0, 1, tview.NewTableCell("Brand").SetSelectable(false))
	table.SetCell(0, 2, tview.NewTableCell("Name").SetSelectable(false))
	table.SetCell(0, 3, tview.NewTableCell("Quantity").SetSelectable(false))
	table.SetSelectedFunc(func(row, column int) {
		items[row-1].Quantity = items[row-1].Quantity + 1
		table.GetCell(row, 3).SetText(strconv.Itoa(items[row-1].Quantity))
	})

	for i, v := range items {
		table.SetCellSimple(i+1, 0, v.Barcode)
		table.SetCellSimple(i+1, 1, v.Brand)
		table.SetCellSimple(i+1, 2, v.Name)
		table.SetCellSimple(i+1, 3, strconv.Itoa(v.Quantity))
	}

	if err := tview.NewApplication().SetRoot(table, true).Run(); err != nil {
		panic(err)
	}
}
