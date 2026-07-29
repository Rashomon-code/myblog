package handle

import (
	"github.com/Rashomon-code/myblog/internal/service"
	"github.com/gin-gonic/gin"
)

type PostHandle struct {
	authService *service.PostService
}

func NewPostHandle(s *service.PostService) *PostHandle {
	return &PostHandle{authService: s}
}

func (h *PostHandle) CreatePostHandle(c *gin.Context) {

}
