package main

import (
	"log"
	"os"

	"github.com/Rashomon-code/myblog/internal/handle"
	"github.com/Rashomon-code/myblog/internal/middleware"
	"github.com/Rashomon-code/myblog/internal/repository"
	"github.com/Rashomon-code/myblog/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	secret := os.Getenv("JWT_SECRET")

	db, err := repository.InitSQL()
	if err != nil {
		log.Fatalln("データベース初期化失敗:", err)
	}
	defer db.Close()

	jwtService := service.NewJWTService(secret)

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, jwtService)
	authHandle := handle.NewAuthHandle(authService)

	middleware := middleware.NewMiddleware(jwtService)

	postRepo := repository.NewPostRepository(db)
	postService := service.NewPostService(postRepo)
	postHandle := handle.NewPostHandle(postService)

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

	r.GET("/create", func(ctx *gin.Context) {
		ctx.File("templates/post.html")
	})
	r.GET("/mypage", func(ctx *gin.Context) {
		ctx.File("templates/mypage.html")
	})

	r.GET("/post", func(ctx *gin.Context) {
		ctx.File("templates/post_detail.html")
	})
	r.GET("/post/detail", postHandle.PostDetailAPI)

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.POST("/posts", postHandle.CreatePostAPI)
		api.GET("/mypage", postHandle.MyPageAPI)
	}

	r.Run(":8080")
}
