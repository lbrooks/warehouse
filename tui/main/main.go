package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/lbrooks/warehouse"

	"github.com/rivo/tview"
)

const WAREHOUSE_SERVER_URL = "http://localhost:8080"

func getData() []warehouse.Item {
	resp, err := http.Get(WAREHOUSE_SERVER_URL + "/api/item")
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var itemList []warehouse.Item
	err = json.Unmarshal(body, &itemList)
	if err != nil {
		// handle error
	}

	return itemList
}

func main() {
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
