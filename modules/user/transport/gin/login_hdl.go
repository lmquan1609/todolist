package ginuser

import (
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"todolist/common"
	userbiz "todolist/modules/user/biz"
	usermodel "todolist/modules/user/model"
	userstorage "todolist/modules/user/storage"
	"todolist/plugin/tokenprovider"
)

func Login(serviceCtx goservice.ServiceContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)
		tokenProvider := serviceCtx.MustGet(common.PluginJWT).(tokenprovider.Provider)

		var loginUserData usermodel.UserLogin

		if err := c.ShouldBind(&loginUserData); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		store := userstorage.NewSQLStore(db)
		md5 := common.NewMd5Hash()

		biz := userbiz.NewLoginBiz(store, tokenProvider, md5, 60*60*24*30)
		account, err := biz.Login(c.Request.Context(), &loginUserData)

		if err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(account))
	}
}
