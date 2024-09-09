package ginuserlikeitem

import (
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"todolist/common"
	userlikeitemstorage "todolist/modules/userlikeitem/storage"
)

func GetItemLikes(serviceCtx goservice.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		type RequestData struct {
			Ids []int `json:"ids"`
		}

		var data RequestData

		if err := c.ShouldBind(&data); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)
		store := userlikeitemstorage.NewSQLStore(db)
		mapRs, err := store.GetItemLikes(c.Request.Context(), data.Ids)

		if err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(mapRs))
	}
}
