package ginitem

import (
	"net/http"
	"todolist/common"
	"todolist/modules/item/biz"
	"todolist/modules/item/model"
	"todolist/modules/item/storage"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ListItem(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var queryStr struct {
			common.Paging
			model.Filter
		}

		if err := c.ShouldBind(&queryStr); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		store := storage.NewSQLStore(db)
		requester := c.MustGet(common.CurrentUser).(common.Requester)
		biz := biz.NewListItemBiz(store, requester)

		data, err := biz.ListItemBiz(c.Request.Context(), &queryStr.Filter, &queryStr.Paging)
		if err != nil {
			panic(err)
		}

		for i := range data {
			data[i].Mask()
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(data, queryStr.Paging, queryStr.Filter))
	}
}
