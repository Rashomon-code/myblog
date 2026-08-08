package main

import (
	"log"
	"os"

	"github.com/Rashomon-code/myblog/internal/handle"
	"github.com/Rashomon-code/myblog/internal/middleware"
	"github.com/Rashomon-code/myblog/internal/repository"
	"github.com/Rashomon-code/myblog/internal/router"
	"github.com/Rashomon-code/myblog/internal/service"
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

	r := router.SetupRouter(authHandle, postHandle, middleware)

	r.Run(":8080")
}
