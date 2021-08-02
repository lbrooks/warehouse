package server

import (
	"github.com/lbrooks/warehouse"

	"github.com/gin-gonic/gin"
)

type ItemController struct {
	itemService warehouse.ItemService
}

func NewItemController(itemService warehouse.ItemService) *ItemController {
	return &ItemController{
		itemService: itemService,
	}
}

// Counts godoc
// @Summary Get All Counts
// @Description Get all counts
// @ID count-items
// @Accept json
// @Produce json
// @Success 200 {object} map[string]int
// @Router /api/item/counts [get]
func (c *ItemController) Counts(ginCtx *gin.Context) {
	v, err := c.itemService.GetCounts(ginCtx.Request.Context())
	if err != nil {
		ginCtx.JSON(500, err.Error())
	} else {
		ginCtx.JSON(200, v)
	}
}

// Search godoc
// @Summary Get All Items Matching Filter
// @Description Get all items matching filter
// @ID search-items
// @Accept json
// @Produce json
// @Param filter path warehouse.Item true "Filter"
// @Success 200 {array} warehouse.Item
// @Router /api/item [get]
func (c *ItemController) Search(ginCtx *gin.Context) {
	var item warehouse.Item
	err := ginCtx.Bind(&item)
	if err != nil {
		ginCtx.JSON(500, err.Error())
	} else {
		v, err := c.itemService.Search(ginCtx.Request.Context(), item)
		if err != nil {
			ginCtx.JSON(500, err.Error())
		} else {
			ginCtx.JSON(200, v)
		}
	}
}

// Update godoc
// @Summary Add or Update Item
// @Description Add or Update Item
// @ID update-item
// @Accept json
// @Produce json
// @Param item path warehouse.Item true "Item"
// @Success 200 {string}
// @Router /api/item [post]
func (c *ItemController) Update(ginCtx *gin.Context) {
	var item warehouse.Item
	err := ginCtx.Bind(&item)
	if err != nil {
		ginCtx.JSON(500, err.Error())
	} else {
		v, err := c.itemService.Update(ginCtx.Request.Context(), item)
		if err != nil {
			ginCtx.JSON(500, err.Error())
		} else {
			ginCtx.JSON(200, v)
		}
	}
}
