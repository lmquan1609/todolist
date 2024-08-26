package ginuserlikeitem

import (
	"fmt"
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/btcsuite/btcutil/base58"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"todolist/common"
	userlikeitembiz "todolist/modules/userlikeitem/biz"
	userlikeitemstorage "todolist/modules/userlikeitem/storage"
)

const timeLayout = "2006-01-02T15:04:05.999999"

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
			if i == len(result)-1 {
				cursorStr := base58.Encode([]byte(fmt.Sprintf("%v", result[i].CreatedAt.Format(timeLayout))))
				queryStr.Paging.NextCursor = cursorStr
			}
		}
		c.JSON(http.StatusOK, common.NewSuccessResponse(result, queryStr.Paging, nil))
	}
}
