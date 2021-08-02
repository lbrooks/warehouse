package tui

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/lbrooks/warehouse"
	"github.com/rivo/tview"
	"go.opentelemetry.io/otel/attribute"
)

type Gui struct {
	app *tview.Application

	appFlex *tview.Flex

	itemTable *tview.Table

	formFlex    *tview.Form
	formBarcode *tview.InputField
	formName    *tview.InputField
	formBrand   *tview.InputField

	service warehouse.ItemService
}

func New(s warehouse.ItemService) *Gui {
	return &Gui{
		app:       tview.NewApplication(),
		appFlex:   tview.NewFlex(),
		itemTable: tview.NewTable(),
		formFlex:  tview.NewForm(),
		service:   s,
	}
}

// Start start application
func (g *Gui) Start() error {
	g.initPanels()
	if err := g.app.Run(); err != nil {
		g.app.Stop()
		return err
	}

	return nil
}

func (g *Gui) initPanels() {
	g.appFlex.SetDirection(tview.FlexRow)

	g.itemTable.SetFixed(1, 0)
	g.itemTable.SetSeparator('|')
	g.itemTable.SetSelectable(true, false)
	g.itemTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'r':
			g.refreshItemList(context.TODO())
			return nil
		case 'p':
			row, _ := g.itemTable.GetSelection()
			item := g.itemTable.GetCell(row, 0).GetReference().(*warehouse.Item)
			g.increment(context.TODO(), item)
			return nil
		case 'm':
			row, _ := g.itemTable.GetSelection()
			item := g.itemTable.GetCell(row, 0).GetReference().(*warehouse.Item)
			g.decrement(context.TODO(), item)
			return nil
		case 'n':
			g.createForm()
			g.app.SetFocus(g.formFlex)
			return nil
		default:
			return event
		}
	})
	g.appFlex.AddItem(g.itemTable, 0, 1, true)

	g.formFlex.SetHorizontal(true)

	g.app.SetRoot(g.appFlex, true)

	g.itemTable.Clear()
	g.setTableHeader()
}

func (g *Gui) createForm() {
	g.formFlex.AddInputField("Barcode", "", 30, nil, nil)
	g.formFlex.AddInputField("Name", "", 30, nil, nil)
	g.formFlex.AddInputField("Brand", "", 30, nil, nil)
	g.formFlex.AddButton("Create", func() {
		g.increment(context.TODO(), &warehouse.Item{
			Barcode:  g.formFlex.GetFormItemByLabel("Barcode").(*tview.InputField).GetText(),
			Name:     g.formFlex.GetFormItemByLabel("Name").(*tview.InputField).GetText(),
			Brand:    g.formFlex.GetFormItemByLabel("Brand").(*tview.InputField).GetText(),
			Quantity: 0,
		})
		g.formFlex.Clear(true)
		g.appFlex.RemoveItem(g.formFlex)
		g.app.SetFocus(g.itemTable)
	})

	g.appFlex.AddItem(g.formFlex, 3, 0, true)
}

func (g *Gui) refreshItemList(ctx context.Context) {
	sc, span := warehouse.CreateSpan(ctx, "ui", "fetch-items")
	defer span.End()

	items, err := g.service.Search(sc, warehouse.Item{})
	if err != nil {
		panic(err.Error())
	}
	g.updateTableContents(items)
}

func (g *Gui) increment(ctx context.Context, item *warehouse.Item) {
	sc, span := warehouse.CreateSpan(ctx, "ui", "increment-item")
	defer span.End()

	g.update(sc, warehouse.Item{
		Barcode:  item.Barcode,
		Brand:    item.Brand,
		Name:     item.Name,
		Quantity: item.Quantity + 1,
	})
}

func (g *Gui) decrement(ctx context.Context, item *warehouse.Item) {
	sc, span := warehouse.CreateSpan(ctx, "ui", "decrement-item")
	defer span.End()

	g.update(sc, warehouse.Item{
		Barcode:  item.Barcode,
		Brand:    item.Brand,
		Name:     item.Name,
		Quantity: item.Quantity - 1,
	})
}

func (g *Gui) update(ctx context.Context, item warehouse.Item) {
	sc, span := warehouse.CreateSpan(ctx, "ui", "update-item")
	defer span.End()

	itemLog, err := json.Marshal(item)
	if err != nil {
		panic("couldn't serialize item")
	}
	span.SetAttributes(attribute.String("item", string(itemLog)))

	_, err = g.service.Update(sc, item)
	if err != nil {
		panic(err.Error())
	}

	g.refreshItemList(sc)
}

func (g *Gui) updateTableContents(items []*warehouse.Item) {
	g.itemTable.Clear()
	g.setTableHeader()
	for i, v := range items {
		g.itemTable.SetCell(i+1, 0, tview.NewTableCell(v.Name).SetReference(v))
		g.itemTable.SetCell(i+1, 1, tview.NewTableCell(v.Brand))
		g.itemTable.SetCell(i+1, 2, tview.NewTableCell(strconv.Itoa(v.Quantity)).SetAlign(tview.AlignRight))
	}
}

func (g *Gui) setTableHeader() {
	g.itemTable.SetCell(0, 0, tview.NewTableCell("Name").SetSelectable(false).SetExpansion(1).SetBackgroundColor(tcell.ColorDarkRed))
	g.itemTable.SetCell(0, 1, tview.NewTableCell("Brand").SetSelectable(false).SetExpansion(1).SetBackgroundColor(tcell.ColorDarkRed))
	g.itemTable.SetCell(0, 2, tview.NewTableCell("Quantity").SetSelectable(false).SetAlign(tview.AlignRight).SetBackgroundColor(tcell.ColorDarkRed))
}
