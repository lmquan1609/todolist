package ginitem

import (
	goservice "github.com/200Lab-Education/go-sdk"
	"net/http"
	"todolist/common"
	"todolist/modules/item/biz"
	"todolist/modules/item/model"
	"todolist/modules/item/repository"
	"todolist/modules/item/storage"
	userlikeitemstorage "todolist/modules/userlikeitem/storage"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ListItem(serviceCtx goservice.ServiceContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)
		var queryStr struct {
			common.Paging
			model.Filter
		}

		if err := c.ShouldBind(&queryStr); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		store := storage.NewSQLStore(db)
		requester := c.MustGet(common.CurrentUser).(common.Requester)
		likeStore := userlikeitemstorage.NewSQLStore(db)
		repo := repository.NewListItemRepo(store, likeStore, requester)
		biz := biz.NewListItemBiz(repo, requester)

		data, err := biz.ListItemBiz(c.Request.Context(), &queryStr.Filter, &queryStr.Paging)
		if err != nil {
			panic(err)
		}

		for i := range data {
			data[i].Mask()
			if i == len(data)-1 {
				queryStr.Paging.NextCursor = data[len(data)-1].FakeId.String()
			}
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(data, queryStr.Paging, queryStr.Filter))
	}
}
