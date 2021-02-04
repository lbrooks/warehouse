package http

import (
	"github.com/lbrooks/warehouse"
	"github.com/lbrooks/warehouse/server"

	"github.com/gin-gonic/gin"
)

// AddRoutes Add Testing Routes
func AddRoutes(api *gin.RouterGroup, s server.ItemService) {
	itemAPIRoutes := api.Group("item")

	// itemAPIRoutes.GET("", TracedRoute("routes", "get", func(spanCtx context.Context, ginCtx *gin.Context) (interface{}, error) {
	itemAPIRoutes.GET("", func(ginCtx *gin.Context) {
		var item warehouse.Item
		err := ginCtx.Bind(&item)
		if err != nil {
			ginCtx.JSON(500, err.Error())
		} else {
			v, err := s.Search(ginCtx.Request.Context(), item)
			if err != nil {
				ginCtx.JSON(500, err.Error())
			} else {
				ginCtx.JSON(200, v)
			}
		}
	})

	// itemAPIRoutes.POST("update", TracedRoute("routes", "update", func(spanCtx context.Context, ginCtx *gin.Context) (interface{}, error) {
	itemAPIRoutes.POST("update", func(ginCtx *gin.Context) {
		var item warehouse.Item
		err := ginCtx.Bind(&item)
		if err != nil {
			ginCtx.JSON(500, err.Error())
		} else {
			v, err := s.Update(ginCtx.Request.Context(), item)
			if err != nil {
				ginCtx.JSON(500, err.Error())
			} else {
				ginCtx.JSON(200, v)
			}
		}
	})
}
