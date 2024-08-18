package ginuser

import (
	goservice "github.com/200Lab-Education/go-sdk"
	"net/http"
	"todolist/common"
	userbiz "todolist/modules/user/biz"
	usermodel "todolist/modules/user/model"
	userstorage "todolist/modules/user/storage"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(serviceCtx goservice.ServiceContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)
		var data usermodel.UserCreate

		if err := c.ShouldBind(&data); err != nil {
			panic(err)
		}

		store := userstorage.NewSQLStore(db)
		md5 := common.NewMd5Hash()
		biz := userbiz.NewRegisterBiz(store, md5)

		if err := biz.Register(c.Request.Context(), &data); err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data.Id))
	}
}
