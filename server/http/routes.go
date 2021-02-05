package http

import (
	"github.com/lbrooks/warehouse"

	"github.com/gin-gonic/gin"
)

// AddRoutes Add Testing Routes
func AddRoutes(api *gin.RouterGroup, s warehouse.ItemService) {
	itemAPIRoutes := api.Group("item")

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

	itemAPIRoutes.POST("", func(ginCtx *gin.Context) {
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
