package ginitem

import (
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"todolist/common"
	"todolist/modules/item/biz"
	"todolist/modules/item/model"
	"todolist/modules/item/storage"
)

func CreateItem(serviceCtx goservice.ServiceContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)
		var itemData model.TodoItemCreation

		if err := c.ShouldBind(&itemData); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		requester := c.MustGet(common.CurrentUser).(common.Requester)
		itemData.UserId = requester.GetUserId()

		store := storage.NewSQLStore(db)
		biz := biz.NewCreateItemBiz(store)

		if err := biz.CreateNewItem(c.Request.Context(), &itemData); err != nil {
			panic(err)

		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(itemData.Id))
	}
}
