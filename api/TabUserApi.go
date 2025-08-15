package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

// TabUserApi 用户 api 层
type TabUserApi struct {
	tabUserService service.TabUserService
}

// InsertTabUser 新增用户
// @Tags 用户模块
// @Summary 新增用户-Summary
// @Description 新增用户-Description
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param tabUser body model.TabUser true "修改用户实体类"
// @Success 200 {object} commonModel.JsonResult "新增用户响应数据"
// @Failure 401 {object} commonModel.JsonResult "未授权"
// @Failure 500 {object} commonModel.JsonResult "服务器异常"
// @Router /tab/user [post]
func (this *TabUserApi) InsertTabUser(c *gin.Context) {
	var tabUser model.TabUser
	c.ShouldBindJSON(&tabUser)

	res, err := this.tabUserService.InsertTabUser(&tabUser)

	if err != nil {
		c.JSON(http.StatusInternalServerError, commonModel.JsonResultError("新增用户失败："+err.Error()))
		return
	}
	c.JSON(http.StatusOK, commonModel.JsonResultSuccess[any](res))
}

// UpdateTabUser 修改用户
// @Tags 用户模块
// @Summary 修改用户-Summary
// @Description 修改用户-Description
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param tabUser body model.TabUser true "修改用户实体类"
// @Success 200 {object} commonModel.JsonResult "修改用户响应数据"
// @Failure 401 {object} commonModel.JsonResult "未授权"
// @Failure 500 {object} commonModel.JsonResult "服务器异常"
// @Router /tab/user [put]
func (this *TabUserApi) UpdateTabUser(c *gin.Context) {
	var tabUser model.TabUser
	c.ShouldBindJSON(&tabUser)

	res, err := this.tabUserService.UpdateTabUser(&tabUser)

	if err != nil {
		c.JSON(http.StatusInternalServerError, commonModel.JsonResultError("修改用户失败："+err.Error()))
		return
	}
	c.JSON(http.StatusOK, commonModel.JsonResultSuccess[any](res))
}

// DeleteTabUser 删除用户
// @Tags 用户模块
// @Summary 删除用户-Summary
// @Description 删除用户-Description
// @Security BearerAuth
// @Accept path
// @Produce path
// @Param idList path string true "用户主键List"
// @Success 200 {object} commonModel.JsonResult "删除用户响应数据"
// @Failure 401 {object} commonModel.JsonResult "未授权"
// @Failure 500 {object} commonModel.JsonResult "服务器异常"
// @Router /tab/user/:idList [delete]
func (this *TabUserApi) DeleteTabUser(c *gin.Context) {
	idListStr := c.Param("idList") // 例如: "1,2,3"
	idList, err := commonUtils.ParseIds(idListStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, commonModel.JsonResultError("参数错误："+err.Error()))
		return
	}

	res, err := this.tabUserService.DeleteTabUser(idList)

	if err != nil {
		c.JSON(http.StatusInternalServerError, commonModel.JsonResultError("删除用户失败："+err.Error()))
		return
	}
	c.JSON(http.StatusOK, commonModel.JsonResultSuccess[any](res))
}

// GetTabUserById 获取用户详细信息
// @Tags 用户模块
// @Summary 获取用户详细信息-Summary
// @Description 获取用户详细信息-Description
// @Security BearerAuth
// @Accept path
// @Produce path
// @Param id path int64 true "用户主键List"
// @Success 200 {object} commonModel.JsonResult "获取用户详细信息"
// @Failure 401 {object} commonModel.JsonResult "未授权"
// @Failure 500 {object} commonModel.JsonResult "服务器异常"
// @Router /tab/user/:id [get]
func (this *TabUserApi) GetTabUserById(c *gin.Context) {
	idStr := c.Param("id") // 例如: "1"
	id, err := commonUtils.ParseId(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, commonModel.JsonResultError("参数错误："+err.Error()))
		return
	}

	res, err := this.tabUserService.GetTabUserById(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, commonModel.JsonResultError("查询用户失败："+err.Error()))
		return
	}
	c.JSON(http.StatusOK, commonModel.JsonResultSuccess[any](res))
}

// GetTabUserList 查询用户列表
// @Tags 用户模块
// @Summary 查询用户列表-Summary
// @Description 查询用户列表-Description
// @Security BearerAuth
// @Accept param
// @Produce param
// @Param tabUserDto body model.TabUserDto true "用户实体Dto"
// @Success 200 {object} commonModel.JsonResult "查询用户列表响应数据"
// @Failure 401 {object} commonModel.JsonResult "未授权"
// @Failure 500 {object} commonModel.JsonResult "服务器异常"
// @Router /tab/user/list [get]
func (this *TabUserApi) GetTabUserList(c *gin.Context) {
	var tabUserDto dto.TabUserDto
	c.ShouldBindQuery(&tabUserDto)

	res, err := this.tabUserService.GetTabUserList(&tabUserDto)

	if err != nil {
		c.JSON(http.StatusInternalServerError, commonModel.JsonResultError("查询用户列表失败："+err.Error()))
		return
	}
	c.JSON(http.StatusOK, commonModel.JsonResultSuccess[any](res))
}

// GetTabUserPageList 分页查询用户列表
// @Tags 用户模块
// @Summary 分页查询用户列表-Summary
// @Description 分页查询用户列表-Description
// @Security BearerAuth
// @Accept param
// @Produce param
// @Param tabUserDto body model.TabUserDto true "用户实体Dto"
// @Success 200 {object} commonModel.JsonResult "分页查询用户列表响应数据"
// @Failure 401 {object} commonModel.JsonResult "未授权"
// @Failure 500 {object} commonModel.JsonResult "服务器异常"
// @Router /tab/user/pageList [get]
func (this *TabUserApi) GetTabUserPageList(c *gin.Context) {
	var tabUserDto dto.TabUserDto
	c.ShouldBindQuery(&tabUserDto)

	res, err := this.tabUserService.GetTabUserPageList(&tabUserDto)

	if err != nil {
		c.JSON(http.StatusInternalServerError, commonModel.JsonResultError("查询用户列表失败："+err.Error()))
		return
	}
	c.JSON(http.StatusOK, commonModel.JsonResultSuccess[any](res))
}

// ExportTabUser 导出用户列表
// @Tags 用户模块
// @Summary 导出用户列表-Summary
// @Description 导出用户列表-Description
// @Security BearerAuth
// @Accept param
// @Produce param
// @Param tabUserDto body model.TabUserDto true "用户实体Dto"
// @Success 200 {object} commonModel.JsonResult "导出用户列表响应数据"
// @Failure 401 {object} commonModel.JsonResult "未授权"
// @Failure 500 {object} commonModel.JsonResult "服务器异常"
// @Router /tab/user/export [get]
func (this *TabUserApi) ExportTabUser(c *gin.Context) {
	var tabUserDto dto.TabUserDto
	c.ShouldBindQuery(&tabUserDto)

	res, err := this.tabUserService.ExportTabUser(&tabUserDto)

	if err != nil {
		c.JSON(http.StatusInternalServerError, commonModel.JsonResultError("导出用户列表失败："+err.Error()))
		return
	}
	c.JSON(http.StatusOK, commonModel.JsonResultSuccess[any](res))
}
