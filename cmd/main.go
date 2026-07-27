package main

import (
	"log"

	"github.com/Rashomon-code/myblog/internal/handle"
	"github.com/Rashomon-code/myblog/internal/repository"
	"github.com/Rashomon-code/myblog/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := repository.InitSQL()
	if err != nil {
		log.Fatalln("データベース初始化失敗:", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandle := handle.NewAuthHandle(authService)

	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.File("templates/index.html")
	})
	r.GET("/register", func(ctx *gin.Context) {
		ctx.File("templates/register.html")
	})
	r.POST("/register", authHandle.RegisterAPI)

	r.GET("/login", func(ctx *gin.Context) {
		ctx.File("templates/login.html")
	})
	r.POST("/login", authHandle.LoginAPI)

	r.Run(":8080")
}
