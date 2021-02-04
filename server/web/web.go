package web

import (
	"github.com/gin-gonic/gin"
)

// NewWebServer create a new web server
func NewWebServer() *gin.Engine {
	r := gin.Default()
	return r
}
