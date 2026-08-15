package handle

import (
	"log"
	"net/http"

	"github.com/Rashomon-code/myblog/internal/model"
	"github.com/Rashomon-code/myblog/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandle struct {
	userService *service.UserService
	postService *service.PostService
}

func NewUserHandle(userService *service.UserService, postService *service.PostService) *UserHandle {
	return &UserHandle{
		userService: userService,
		postService: postService,
	}
}

func (h *UserHandle) MyPageAPI(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"エラー": "ユーザーが見つかりませんでした"})
		return
	}

	userID := userIDVal.(int64)
	posts, err := h.postService.GetPostTitle(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"エラー": "ポスト獲得できませんでした"})
		log.Printf("Bind error: %v", err)
		return
	}

	user, err := h.userService.GetProfileService(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"エラー": "ユーザー情報が獲得できませんでした"})
		return
	}

	userPage := model.MyPageResponse{
		UserProfile: *user,
		Posts:       posts,
	}

	c.JSON(http.StatusOK, userPage)
}
