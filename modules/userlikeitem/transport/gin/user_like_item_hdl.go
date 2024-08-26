package ginuserlikeitem

import (
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"todolist/common"
	"todolist/modules/item/storage"
	userlikeitembiz "todolist/modules/userlikeitem/biz"
	userlikeitemmodel "todolist/modules/userlikeitem/model"
	userlikeitemstorage "todolist/modules/userlikeitem/storage"
)

func LikeItem(serviceCtx goservice.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := common.FromBase58(c.Param("id"))

		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		requester := c.MustGet(common.CurrentUser).(common.Requester)
		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)

		store := userlikeitemstorage.NewSQLStore(db)
		itemStore := storage.NewSQLStore(db)
		biz := userlikeitembiz.NewUserLikeItemBiz(store, itemStore)

		if err := biz.LikeItem(c.Request.Context(), &userlikeitemmodel.Like{
			ItemId: int(id.GetLocalID()),
			UserId: requester.GetUserId(),
		}); err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
