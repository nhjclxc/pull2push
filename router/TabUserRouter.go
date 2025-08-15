package router

import (
	"github.com/gin-gonic/gin"
)

// TabUserRouter 用户 路由Router层
type TabUserRouter struct {
	tabUserApi api.TabUserApi
}

// InitTabUserRouter 初始化 TabUserRouter 路由
func (this *TabUserRouter) InitTabUserRouter(privateRouterOrigin *gin.RouterGroup, publicRouterOrigin *gin.RouterGroup) {
	privateRouter := privateRouterOrigin.Group("/tab/user")
	{
		// PrivateRouter 下是一些必须进行登录的接口
		// http://localhost:8080/private

		privateRouter.POST("", this.tabUserApi.InsertTabUser)              // 新增用户
		privateRouter.PUT("", this.tabUserApi.UpdateTabUser)               // 修改用户
		privateRouter.DELETE("/:idList", this.tabUserApi.DeleteTabUser)    // 删除用户
		privateRouter.GET("/:id", this.tabUserApi.GetTabUserById)          // 获取用户详细信息
		privateRouter.GET("/list", this.tabUserApi.GetTabUserList)         // 查询用户列表
		privateRouter.GET("/pageList", this.tabUserApi.GetTabUserPageList) // 分页查询用户列表
		privateRouter.GET("/export", this.tabUserApi.ExportTabUser)        // 导出用户列表
	}

	publicRouter := publicRouterOrigin.Group("/tab/user")
	{
		// PublicRouter 下是一些无需登录的接口，可以直接访问，无须经过授权操作
		// http://localhost:8080/public

		publicRouter.POST("", this.tabUserApi.InsertTabUser)              // 新增用户
		publicRouter.PUT("", this.tabUserApi.UpdateTabUser)               // 修改用户
		publicRouter.DELETE("/:id", this.tabUserApi.DeleteTabUser)        // 删除用户
		publicRouter.GET("/:id", this.tabUserApi.GetTabUserById)          // 获取用户详细信息
		publicRouter.GET("/list", this.tabUserApi.GetTabUserList)         // 查询用户列表
		publicRouter.GET("/pageList", this.tabUserApi.GetTabUserPageList) // 分页查询用户列表
		publicRouter.GET("/export", this.tabUserApi.ExportTabUser)        // 导出用户列表
	}
}
