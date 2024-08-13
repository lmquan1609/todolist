package ginuser

import (
	"net/http"
	"todolist/common"

	"github.com/gin-gonic/gin"
)

func Profile() func(c *gin.Context) {
	return func(c *gin.Context) {
		user := c.MustGet(common.CurrentUser)
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(user))
	}
}
