package ginitem

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"todolist/common"
	"todolist/modules/item/biz"
	"todolist/modules/item/model"
	"todolist/modules/item/storage"
)

func CreateItem(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var itemData model.TodoItemCreation

		if err := c.ShouldBind(&itemData); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		store := storage.NewSQLStore(db)
		biz := biz.NewCreateItemBiz(store)

		if err := biz.CreateNewItem(c.Request.Context(), &itemData); err != nil {
			panic(err)

		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(itemData.Id))
	}
}
