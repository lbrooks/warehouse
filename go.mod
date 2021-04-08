module github.com/lbrooks/warehouse

go 1.15

require (
	github.com/gdamore/tcell/v2 v2.0.1-0.20201017141208-acf90d56d591
	github.com/gin-gonic/gin v1.6.3
	github.com/joho/godotenv v1.3.0
	github.com/prometheus/client_golang v1.9.0 // indirect
	github.com/rivo/tview v0.0.0-20210125085121-dbc1f32bb1d0
	github.com/sirupsen/logrus v1.7.0 // indirect
	github.com/zsais/go-gin-prometheus v0.1.0
	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.16.0
	go.opentelemetry.io/contrib/propagators v0.16.0
	go.opentelemetry.io/otel v0.16.0
	go.opentelemetry.io/otel/exporters/trace/jaeger v0.16.0
	go.opentelemetry.io/otel/sdk v0.16.0
	google.golang.org/api v0.39.0 // indirect
)
