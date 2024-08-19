package ginuserlikeitem

import (
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"todolist/common"
	userlikeitembiz "todolist/modules/userlikeitem/biz"
	userlikeitemstorage "todolist/modules/userlikeitem/storage"
)

func ListUserLiked(serviceCtx goservice.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := common.FromBase58(c.Param("id"))

		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		var queryStr struct {
			common.Paging
		}

		if err := c.ShouldBind(&queryStr); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		queryStr.Paging.Process()

		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)
		store := userlikeitemstorage.NewSQLStore(db)
		biz := userlikeitembiz.NewListUserLikeItemBiz(store)

		result, err := biz.ListUserLikedItem(c.Request.Context(), int(id.GetLocalID()), &queryStr.Paging)
		if err != nil {
			panic(err)
		}

		for i := range result {
			result[i].Mask()
		}
		c.JSON(http.StatusOK, common.NewSuccessResponse(result, queryStr.Paging, nil))
	}
}
