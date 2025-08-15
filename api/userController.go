package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pull2push/api/base"
	"pull2push/entity/request"
)

// UserController 处理所有流相关的请求
type UserController struct {
	*base.BaseController
}

// NewUserController 创建一个新的 UserController
func NewUserController(base *base.BaseController) *UserController {
	return &UserController{
		BaseController: base,
	}
}

// CreateUser 创建新的媒体流
func (c *UserController) CreateUser(ctx *gin.Context) {
	var createUserRequest request.CreateUserRequest
	if err := ctx.ShouldBindJSON(&createUserRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":  "Invalid request format",
			"detail": err.Error(),
		})
		return
	}

	// todo
	// 创建用户
	// userService.CreateUser(createUserRequest)
	res := struct{}{}

	ctx.JSON(http.StatusOK, base.JsonResultSuccess[any](res))

}
