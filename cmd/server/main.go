package main

import (
	"log"
	"os"

	"github.com/Rashomon-code/myblog/internal/handle"
	"github.com/Rashomon-code/myblog/internal/middleware"
	"github.com/Rashomon-code/myblog/internal/repository"
	"github.com/Rashomon-code/myblog/internal/router"
	"github.com/Rashomon-code/myblog/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	secret := os.Getenv("JWT_SECRET")

	db, err := repository.InitSQL()
	if err != nil {
		log.Fatalln("データベース初期化失敗:", err)
	}
	defer db.Close()

	jwtService := service.NewJWTService(secret)

	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo, jwtService)
	authHandle := handle.NewAuthHandle(authService)
	middleware := middleware.NewMiddleware(jwtService)

	postRepo := repository.NewPostRepository(db)
	postService := service.NewPostService(postRepo)
	postHandle := handle.NewPostHandle(postService)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandle := handle.NewUserHandle(userService, postService)

	r := router.SetupRouter(authHandle, postHandle, middleware, userHandle)

	r.Run(":8080")
}
