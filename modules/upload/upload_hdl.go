package upload

import (
	"fmt"
	"net/http"
	"todolist/common"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Upload(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		fileHeader, err := c.FormFile("file")

		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		dst := fmt.Sprintf("./static/%s", fileHeader.Filename)

		if err := c.SaveUploadedFile(fileHeader, dst); err != nil {

		}
		img := common.Image{
			Id:        0,
			Url:       dst,
			Width:     100,
			Height:    100,
			CloudName: "local",
			Extension: "",
		}
		img.Fulfill("http://localhost:3000")
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(img))
	}
}
