package cmd

import (
	"fmt"
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"os"
	"todolist/common"
	"todolist/component/tokenprovider/jwt"
	"todolist/middleware"
	"todolist/modules/item/transport/ginitem"
	"todolist/modules/upload"
	userstorage "todolist/modules/user/storage"
	ginuser "todolist/modules/user/transport/gin"
	"todolist/plugin/sdkgorm"
)

func newService() goservice.Service {
	service := goservice.New(
		goservice.WithName("todolist"),
		goservice.WithVersion("1.0.0"),
		goservice.WithInitRunnable(sdkgorm.NewGormDB("main", common.PluginDBMain)),
	)
	return service
}

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Start social TODO service",
	Run: func(cmd *cobra.Command, args []string) {
		secret := os.Getenv("SECRET")

		service := newService()

		serviceLogger := service.Logger("service")

		if err := service.Init(); err != nil {
			serviceLogger.Fatalln(err)
		}

		service.HTTPServer().AddHandler(func(e *gin.Engine) {
			e.Use(middleware.Recover())

			db := service.MustGet(common.PluginDBMain).(*gorm.DB)

			authStore := userstorage.NewSQLStore(db)
			tokenProvider := jwt.NewTokenJWTProvider("jwt", secret)
			middlewareAuth := middleware.RequiredAuth(authStore, tokenProvider)

			v1 := e.Group("/v1")
			{
				items := v1.Group("/items")
				{
					items.POST("", middlewareAuth, ginitem.CreateItem(db))
					items.GET("/:id", ginitem.GetItemById(db))
					items.PATCH("/:id", middlewareAuth, ginitem.UpdateItem(db))
					items.DELETE("/:id", middlewareAuth, ginitem.DeleteItem(db))
					items.GET("", middlewareAuth, ginitem.ListItem(db))
				}
				v1.PUT("/upload", upload.Upload(db))

				v1.POST("/register", ginuser.Register(db))
				v1.POST("/login", ginuser.Login(db, tokenProvider))
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
