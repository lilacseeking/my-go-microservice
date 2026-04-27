// api/v1/router.go
package handler

import (
	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes 注册用户相关路由
func RegisterUserRoutes(r *gin.RouterGroup, userHandler *UserHandler) {
	users := r.Group("/users")
	{
		users.POST("", userHandler.CreateUserHandler) // POST /v1/users
		users.GET("/:id", userHandler.GetUserHandler) // GET /v1/users/:id
	}
}
