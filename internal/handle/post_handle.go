package handle

import (
	"net/http"

	"github.com/Rashomon-code/myblog/internal/model"
	"github.com/Rashomon-code/myblog/internal/service"
	"github.com/gin-gonic/gin"
)

type PostHandle struct {
	postService *service.PostService
}

func NewPostHandle(s *service.PostService) *PostHandle {
	return &PostHandle{postService: s}
}

func (h *PostHandle) CreatePostAPI(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"エラー": "ユーザーが見つかりませんでした"})
		return
	}
	userID := userIDVal.(int64)

	var req model.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"エラー": "入力に誤りがございます: " + err.Error()})
		return
	}

	err := h.postService.CreatePostService(userID, req.Title, req.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"エラー": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"メッセージ": "投稿しました。"})
}
