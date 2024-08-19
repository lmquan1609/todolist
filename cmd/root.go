package cmd

import (
	"fmt"
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"os"
	"todolist/common"
	"todolist/middleware"
	"todolist/modules/item/transport/ginitem"
	"todolist/modules/upload"
	userstorage "todolist/modules/user/storage"
	ginuser "todolist/modules/user/transport/gin"
	ginuserlikeitem "todolist/modules/userlikeitem/transport/gin"
	"todolist/plugin/sdkgorm"
	"todolist/plugin/tokenprovider/jwt"
)

func newService() goservice.Service {
	service := goservice.New(
		goservice.WithName("todolist"),
		goservice.WithVersion("1.0.0"),
		goservice.WithInitRunnable(sdkgorm.NewGormDB("main", common.PluginDBMain)),
		goservice.WithInitRunnable(jwt.NewJWTProvider(common.PluginJWT)),
	)
	return service
}

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Start social TODO service",
	Run: func(cmd *cobra.Command, args []string) {
		service := newService()

		serviceLogger := service.Logger("service")

		if err := service.Init(); err != nil {
			serviceLogger.Fatalln(err)
		}

		service.HTTPServer().AddHandler(func(e *gin.Engine) {
			e.Use(middleware.Recover())

			db := service.MustGet(common.PluginDBMain).(*gorm.DB)

			authStore := userstorage.NewSQLStore(db)
			middlewareAuth := middleware.RequiredAuth(authStore, service)

			v1 := e.Group("/v1")
			{
				items := v1.Group("/items")
				{
					items.POST("", middlewareAuth, ginitem.CreateItem(service))
					items.GET("/:id", ginitem.GetItemById(service))
					items.PATCH("/:id", middlewareAuth, ginitem.UpdateItem(service))
					items.DELETE("/:id", middlewareAuth, ginitem.DeleteItem(service))
					items.GET("", middlewareAuth, ginitem.ListItem(service))

					items.POST("/:id/like", middlewareAuth, ginuserlikeitem.LikeItem(service))
					items.DELETE("/:id/unlike", middlewareAuth, ginuserlikeitem.UnlikeItem(service))
					items.GET("/:id/liked-users", ginuserlikeitem.ListUserLiked(service))
				}
				v1.PUT("/upload", upload.Upload(db))

				v1.POST("/register", ginuser.Register(service))
				v1.POST("/login", ginuser.Login(service))
				v1.GET("/profile", middlewareAuth, ginuser.Profile())
			}
		})
		if err := service.Start(); err != nil {
			serviceLogger.Fatalln(err)
		}
	},
}

func Execute() {
	rootCmd.AddCommand(outEnvCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
