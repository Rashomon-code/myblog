package handle

import (
	"log"
	"net/http"
	"strconv"

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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ユーザーが見つかりませんでした"})
		return
	}

	userID := userIDVal.(int64)
	posts, err := h.postService.GetPostTitle(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ポスト獲得できませんでした"})
		log.Printf("Bind error: %v", err)
		return
	}

	user, err := h.userService.GetProfileService(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ユーザー情報が獲得できませんでした"})
		return
	}

	userPage := model.MyPageResponse{
		UserProfile: *user,
		Posts:       posts,
	}

	c.JSON(http.StatusOK, userPage)
}

func (h *UserHandle) UpdateRoleAPI(c *gin.Context) {
	targetUserID, _ := strconv.ParseInt(c.Param("id"), 10, 64) // URL から id を引き出す
	currentUserID := c.GetInt64("userID")                      // JWT から id を引き出す

	var req model.UpdateRoleRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効な形式"})
		return
	}

	err = h.userService.UpdateRoleService(currentUserID, targetUserID, req.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新しました"})
}

func (h *UserHandle) GetAllUsersAPI(c *gin.Context) {
	users, err := h.userService.GetAllUsersService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザーデータが取得できませんでした"})
		return
	}

	c.JSON(http.StatusOK, users)
}
