// api/v1/user_handler.go
package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"my-go-microservice/internal/model"
	"my-go-microservice/internal/service"
)

// UserHandler 处理用户相关的HTTP请求
type UserHandler struct {
	UserService service.UserService
}

// CreateUserRequest 定义创建用户请求的数据结构
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required,min=2,max=50"`
	Email string `json:"email" binding:"required,email"`
}

// CreateUserHandler 处理创建用户请求
func (h *UserHandler) CreateUserHandler(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数绑定失败: " + err.Error(),
		})
		return
	}

	user := &model.User{
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: time.Now().UTC(),
	}

	if err := h.UserService.CreateUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建用户失败",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"message": "用户创建成功",
		"data":    map[string]string{"user_id": user.ID},
	})
}

// GetUserHandler 处理获取用户请求
func (h *UserHandler) GetUserHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "用户ID不能为空",
		})
		return
	}

	user, err := h.UserService.GetUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "用户不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "查询成功",
		"data":    user,
	})
}
