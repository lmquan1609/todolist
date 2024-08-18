package ginitem

import (
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"todolist/common"
	"todolist/modules/item/biz"
	"todolist/modules/item/model"
	"todolist/modules/item/storage"
)

func UpdateItem(serviceCtx goservice.ServiceContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		var data model.TodoItemUpdate
		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
		}

		store := storage.NewSQLStore(db)
		requester := c.MustGet(common.CurrentUser).(common.Requester)
		biz := biz.NewUpdateItemBiz(store, requester)
		if err := biz.UpdateItem(c.Request.Context(), id, &data); err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
