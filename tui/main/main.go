package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/lbrooks/warehouse"
	"github.com/lbrooks/warehouse/client"
	"github.com/lbrooks/warehouse/tui"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	flush, ok := warehouse.InitializeJaeger()
	if ok {
		defer flush()
	}

	service := client.NewItemService()

	tui := tui.New(service)
	if err := tui.Start(); err != nil {
		panic(err)
	}
}
