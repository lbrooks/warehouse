package main

import (
	"context"
	"log"
	"os"

	"github.com/lbrooks/warehouse"
	"github.com/lbrooks/warehouse/server/http"
	"github.com/lbrooks/warehouse/server/temp"
	"github.com/lbrooks/warehouse/server/web"

	"github.com/gin-contrib/static"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	flush := warehouse.InitializeJaeger("warehouse-server")
	defer flush()

	webServer := web.NewWebServer()
	webServer.Use(static.Serve("/", static.LocalFile("/Users/lawrence/projects/inv-client-react/build", false)))
	webServer.Use(otelgin.Middleware("warehouse-server"))
	apiRoutes := webServer.Group("api")

	itemService := temp.NewItemService(context.Background(), true)
	http.AddRoutes(apiRoutes, itemService)

	err := webServer.Run() // listen and serve on 0.0.0.0:8080
	if err != nil {
		os.Exit(1)
	}
}
