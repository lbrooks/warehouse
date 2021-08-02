package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swagFiles "github.com/swaggo/files"
	swagGin "github.com/swaggo/gin-swagger"
	prom "github.com/zsais/go-gin-prometheus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/lbrooks/warehouse"
	"github.com/lbrooks/warehouse/server"
)

func init() {
	existingVars := make(map[string]struct{})
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		existingVars[pair[0]] = struct{}{}
	}

	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	if sEnv, ok := os.LookupEnv("SERVER_ENV"); ok {
		if err := godotenv.Load(sEnv); err != nil {
			log.Printf("No [%s] file found", sEnv)
		}
	}

	log.Default().Println("----- Dot Env Loaded Environmental Vars -----")
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if _, ok := existingVars[pair[0]]; !ok {
			log.Default().Printf("\t%s\n", e)
		}
	}
	log.Default().Println("---------------------------------------------")
}

// @title Warehouse API
// @version 1.0
// @description A simple inventory system
// @termsOfService http://swagger.io/terms/

// @contact.name Lawrence Brooks
// @contact.url http://www.github.com/lbrooks

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
func main() {
	itemStore := server.NewItemStore(context.Background(), true)
	itemController := server.NewItemController(itemStore)

	webServer := gin.Default()
	initializeSwagger(webServer)
	flush := initializeJaeger(webServer)
	defer flush()
	initializePrometheus(webServer)
	initializeRoutes(webServer, itemController)

	err := webServer.Run()
	if err != nil {
		os.Exit(1)
	}
}

func initializeSwagger(webServer *gin.Engine) {
	url := swagGin.URL("http://localhost:8080/swagger/doc.json") // The url pointing to API definition
	webServer.GET("/swagger/*any", swagGin.WrapHandler(swagFiles.Handler, url))
}

func initializeJaeger(webServer *gin.Engine) func() {
	flush, ok := warehouse.InitializeJaeger()
	if ok {
		webServer.Use(otelgin.Middleware(os.Getenv("JAEGER_SERVICE_NAME")))
		return flush
	}
	return func() {}
}

func initializePrometheus(webServer *gin.Engine) *prom.Prometheus {
	p := prom.NewPrometheus("gin")
	p.Use(webServer)
	return p
}

func initializeRoutes(route *gin.Engine, c *server.ItemController) {
	apiRoute := route.Group("api")
	itemAPIRoutes := apiRoute.Group("item")
	itemAPIRoutes.GET("", c.Search)
	itemAPIRoutes.POST("", c.Update)
	itemAPIRoutes.GET("count", c.Counts)
}
