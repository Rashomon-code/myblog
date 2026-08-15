package router

import (
	"github.com/Rashomon-code/myblog/internal/handle"
	"github.com/Rashomon-code/myblog/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(authHandle *handle.AuthHandle, postHandle *handle.PostHandle, mw *middleware.Middleware, userHandle *handle.UserHandle) *gin.Engine {
	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.File("templates/index.html")
	})

	r.GET("/register", func(ctx *gin.Context) {
		ctx.File("templates/register.html")
	})

	r.GET("/login", func(ctx *gin.Context) {
		ctx.File("templates/login.html")
	})

	r.GET("/create", func(ctx *gin.Context) {
		ctx.File("templates/post.html")
	})
	r.GET("/mypage", func(ctx *gin.Context) {
		ctx.File("templates/mypage.html")
	})

	r.GET("/posts/:id", func(ctx *gin.Context) {
		ctx.File("templates/post_detail.html")
	})

	r.GET("/posts/:id/edit", func(ctx *gin.Context) {
		ctx.File("templates/edit.html")
	})

	api := r.Group("/api")
	{
		api.GET("/posts", postHandle.PostListAPI)
		api.GET("/posts/:id", postHandle.PostDetailAPI)

		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandle.RegisterAPI)
			auth.POST("/login", authHandle.LoginAPI)
		}

		protected := api.Group("")
		protected.Use(mw.AuthMiddleware())
		{
			protected.POST("/posts", postHandle.CreatePostAPI)
			protected.PUT("/posts/:id", postHandle.EditPostAPI)
			protected.DELETE("/posts/:id", postHandle.DeletePostAPI)
			protected.GET("/me/posts", userHandle.MyPageAPI)
		}
	}

	return r
}
