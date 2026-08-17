package handle

import (
	"net/http"
	"strconv"

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

	var req model.UpdatePostRequest
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

func (h *PostHandle) PostDetailAPI(c *gin.Context) {
	idStr := c.Param("id")
	postID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"エラー": "IDが間違っています"})
		return
	}

	post, err := h.postService.PostDetailService(postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"エラー": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

func (h *PostHandle) DeletePostAPI(c *gin.Context) {
	idStr := c.Param("id")
	postID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"エラー": "IDが間違っています"})
		return
	}

	userID := c.GetInt64("userID")
	userRole := c.GetString("role")

	err = h.postService.DeletePostService(postID, userID, userRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"エラー": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"メッセージ": "削除しました"})
}

func (h *PostHandle) EditPostAPI(c *gin.Context) {
	idStr := c.Param("id")
	postID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"エラー": "IDが間違っています"})
		return
	}

	var req model.UpdatePostRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"エラー": "入力に誤りがございます"})
		return
	}

	userID := c.GetInt64("userID")
	userRole := c.GetString("role")

	err = h.postService.EditPostService(postID, req.Title, req.Content, userID, userRole)
	if err != nil {
		if err.Error() == "更新できませんでした" {
			c.JSON(http.StatusNotFound, gin.H{"エラー": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"エラー": "更新できませんでした: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"メッセージ": "更新しました"})
}

func (h *PostHandle) PostListAPI(c *gin.Context) {
	posts, err := h.postService.PostHomeService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"エラー": "ポスト獲得できませんでした"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})
}
