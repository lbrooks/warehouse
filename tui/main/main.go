package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/lbrooks/warehouse"
	"github.com/lbrooks/warehouse/tui/gui"
	"github.com/lbrooks/warehouse/tui/web"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	flush := warehouse.InitializeJaeger("warehouse-tui")
	defer flush()

	service := web.NewItemService()

	gui := gui.New(service)
	if err := gui.Start(); err != nil {
		panic(err)
	}
}
