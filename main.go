package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"todolist/component/tokenprovider/jwt"
	"todolist/middleware"
	"todolist/modules/item/transport/ginitem"
	"todolist/modules/upload"
	userstorage "todolist/modules/user/storage"
	ginuser "todolist/modules/user/transport/gin"
)

func main() {
	dsn := os.Getenv("DBConnectionStr")
	secret := os.Getenv("SECRET")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("DB connection:", db)

	r := gin.Default()
	r.Use(middleware.Recover())

	authStore := userstorage.NewSQLStore(db)
	tokenProvider := jwt.NewTokenJWTProvider("jwt", secret)
	middlewareAuth := middleware.RequiredAuth(authStore, tokenProvider)

	v1 := r.Group("/v1")
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
	if err := r.Run(":3000"); err != nil {
		log.Fatalln(err)
	}
}
