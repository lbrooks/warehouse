package http

import (
	"bytes"
	"context"
	"io/ioutil"
	"strings"

	"github.com/lbrooks/warehouse"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/label"
)

// TracedRouteHandler function to handle request
type TracedRouteHandler func(spanCtx context.Context, ginCtx *gin.Context) (interface{}, error)

func readAndReset(c *gin.Context) string {
	var bodyBytes []byte

	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	return string(bodyBytes)
}

// TracedRoute wrapper to trace route
func TracedRoute(tracerName, operationName string, handler TracedRouteHandler) func(*gin.Context) {
	return func(c *gin.Context) {
		sc, span := warehouse.CreateSpan(c, tracerName, operationName)
		defer span.End()

		for k, v := range c.Request.Header {
			span.SetAttributes(label.String("request.header."+k, strings.Join(v, " | ")))
		}
		span.SetAttributes(label.String("request.query", c.Request.URL.RawQuery))
		span.SetAttributes(label.String("request.body", readAndReset(c)))

		msg, err := handler(sc, c)

		if err != nil {
			span.SetAttributes(label.Int("http.status_code", 500))
			span.RecordError(err)

			c.JSON(500, err.Error())
			return
		}

		span.SetAttributes(label.Int("http.status_code", 200))
		c.JSON(200, msg)
	}
}
