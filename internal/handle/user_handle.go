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

	user, err := h.userService.GetProfileService(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ユーザー情報が獲得できませんでした"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandle) UserProfileAPI(c *gin.Context) {
	ID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ユーザーが見つかりませんでした"})
		return
	}

	var userID int64
	userIDVal, exists := c.Get("userID")

	if exists {
		switch v := userIDVal.(type) {
		case int64:
			userID = v
		case float64:
			userID = int64(v)
		case int:
			userID = int64(v)
		}
	}

	isMe := userID != 0 && userID == ID

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.Query("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	posts, totalCount, err := h.postService.GetPostTitle(ID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ポスト獲得できませんでした"})
		log.Printf("Bind error: %v", err)
		return
	}

	user, err := h.userService.GetProfileService(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ユーザー情報が獲得できませんでした"})
		return
	}

	userPage := model.PageResult{
		Posts:      posts,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}

	var response model.UserProfileResponse

	response.IsMe = isMe
	response.Data.UserProfile = *user
	response.Data.PageResult = userPage

	c.JSON(http.StatusOK, response)
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

func (h *UserHandle) UpdateProfileAPI(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ユーザーが見つかりませんでした"})
		return
	}
	userID := userIDVal.(int64)

	var profile model.ProfileRequest
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userService.UpdateProfileService(userID, profile.DisplayName, profile.Bio)
	if err != nil {
		if err.Error() == "更新できませんでした" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新できませんでした: " + err.Error()})
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
