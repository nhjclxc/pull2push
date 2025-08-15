package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

// TabUserRoleApi 用户角色关联 api 层
type TabUserRoleApi struct {
	tabUserRoleService service.TabUserRoleService
}

// InsertTabUserRole 新增用户角色关联
// @Tags 用户角色关联模块
// @Summary 新增用户角色关联-Summary
// @Description 新增用户角色关联-Description
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param tabUserRole body model.TabUserRole true "修改用户角色关联实体类"
// @Success 200 {object} commonModel.JsonResult "新增用户角色关联响应数据"
// @Failure 401 {object} commonModel.JsonResult "未授权"
// @Failure 500 {object} commonModel.JsonResult "服务器异常"
// @Router /tab/user/role [post]
func (this *TabUserRoleApi) InsertTabUserRole(c *gin.Context) {
	var tabUserRole model.TabUserRole
	c.ShouldBindJSON(&tabUserRole)

	res, err := this.tabUserRoleService.InsertTabUserRole(&tabUserRole)

	if err != nil {
		c.JSON(http.StatusInternalServerError, commonModel.JsonResultError("新增用户角色关联失败："+err.Error()))
		return
	}
	c.JSON(http.StatusOK, commonModel.JsonResultSuccess[any](res))
}

// UpdateTabUserRole 修改用户角色关联
// @Tags 用户角色关联模块
// @Summary 修改用户角色关联-Summary
// @Description 修改用户角色关联-Description
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param tabUserRole body model.TabUserRole true "修改用户角色关联实体类"
// @Success 200 {object} commonModel.JsonResult "修改用户角色关联响应数据"
// @Failure 401 {object} commonModel.JsonResult "未授权"
// @Failure 500 {object} commonModel.JsonResult "服务器异常"
// @Router /tab/user/role [put]
func (this *TabUserRoleApi) UpdateTabUserRole(c *gin.Context) {
	var tabUserRole model.TabUserRole
	c.ShouldBindJSON(&tabUserRole)

	res, err := this.tabUserRoleService.UpdateTabUserRole(&tabUserRole)

	if err != nil {
		c.JSON(http.StatusInternalServerError, commonModel.JsonResultError("修改用户角色关联失败："+err.Error()))
		return
	}
	c.JSON(http.StatusOK, commonModel.JsonResultSuccess[any](res))
}

// DeleteTabUserRole 删除用户角色关联
// @Tags 用户角色关联模块
// @Summary 删除用户角色关联-Summary
// @Description 删除用户角色关联-Description
// @Security BearerAuth
// @Accept path
// @Produce path
// @Param idList path string true "用户角色关联主键List"
// @Success 200 {object} commonModel.JsonResult "删除用户角色关联响应数据"
// @Failure 401 {object} commonModel.JsonResult "未授权"
// @Failure 500 {object} commonModel.JsonResult "服务器异常"
// @Router /tab/user/role/:idList [delete]
func (this *TabUserRoleApi) DeleteTabUserRole(c *gin.Context) {
	idListStr := c.Param("idList") // 例如: "1,2,3"
	idList, err := commonUtils.ParseIds(idListStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, commonModel.JsonResultError("参数错误："+err.Error()))
		return
	}

	res, err := this.tabUserRoleService.DeleteTabUserRole(idList)

	if err != nil {
		c.JSON(http.StatusInternalServerError, commonModel.JsonResultError("删除用户角色关联失败："+err.Error()))
		return
	}
	c.JSON(http.StatusOK, commonModel.JsonResultSuccess[any](res))
}

// GetTabUserRoleById 获取用户角色关联详细信息
// @Tags 用户角色关联模块
// @Summary 获取用户角色关联详细信息-Summary
// @Description 获取用户角色关联详细信息-Description
// @Security BearerAuth
// @Accept path
// @Produce path
// @Param id path int64 true "用户角色关联主键List"
// @Success 200 {object} commonModel.JsonResult "获取用户角色关联详细信息"
// @Failure 401 {object} commonModel.JsonResult "未授权"
// @Failure 500 {object} commonModel.JsonResult "服务器异常"
// @Router /tab/user/role/:id [get]
func (this *TabUserRoleApi) GetTabUserRoleById(c *gin.Context) {
	idStr := c.Param("id") // 例如: "1"
	id, err := commonUtils.ParseId(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, commonModel.JsonResultError("参数错误："+err.Error()))
		return
	}

	res, err := this.tabUserRoleService.GetTabUserRoleById(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, commonModel.JsonResultError("查询用户角色关联失败："+err.Error()))
		return
	}
	c.JSON(http.StatusOK, commonModel.JsonResultSuccess[any](res))
}

// GetTabUserRoleList 查询用户角色关联列表
// @Tags 用户角色关联模块
// @Summary 查询用户角色关联列表-Summary
// @Description 查询用户角色关联列表-Description
// @Security BearerAuth
// @Accept param
// @Produce param
// @Param tabUserRoleDto body model.TabUserRoleDto true "用户角色关联实体Dto"
// @Success 200 {object} commonModel.JsonResult "查询用户角色关联列表响应数据"
// @Failure 401 {object} commonModel.JsonResult "未授权"
// @Failure 500 {object} commonModel.JsonResult "服务器异常"
// @Router /tab/user/role/list [get]
func (this *TabUserRoleApi) GetTabUserRoleList(c *gin.Context) {
	var tabUserRoleDto dto.TabUserRoleDto
	c.ShouldBindQuery(&tabUserRoleDto)

	res, err := this.tabUserRoleService.GetTabUserRoleList(&tabUserRoleDto)

	if err != nil {
		c.JSON(http.StatusInternalServerError, commonModel.JsonResultError("查询用户角色关联列表失败："+err.Error()))
		return
	}
	c.JSON(http.StatusOK, commonModel.JsonResultSuccess[any](res))
}

// GetTabUserRolePageList 分页查询用户角色关联列表
// @Tags 用户角色关联模块
// @Summary 分页查询用户角色关联列表-Summary
// @Description 分页查询用户角色关联列表-Description
// @Security BearerAuth
// @Accept param
// @Produce param
// @Param tabUserRoleDto body model.TabUserRoleDto true "用户角色关联实体Dto"
// @Success 200 {object} commonModel.JsonResult "分页查询用户角色关联列表响应数据"
// @Failure 401 {object} commonModel.JsonResult "未授权"
// @Failure 500 {object} commonModel.JsonResult "服务器异常"
// @Router /tab/user/role/pageList [get]
func (this *TabUserRoleApi) GetTabUserRolePageList(c *gin.Context) {
	var tabUserRoleDto dto.TabUserRoleDto
	c.ShouldBindQuery(&tabUserRoleDto)

	res, err := this.tabUserRoleService.GetTabUserRolePageList(&tabUserRoleDto)

	if err != nil {
		c.JSON(http.StatusInternalServerError, commonModel.JsonResultError("查询用户角色关联列表失败："+err.Error()))
		return
	}
	c.JSON(http.StatusOK, commonModel.JsonResultSuccess[any](res))
}

// ExportTabUserRole 导出用户角色关联列表
// @Tags 用户角色关联模块
// @Summary 导出用户角色关联列表-Summary
// @Description 导出用户角色关联列表-Description
// @Security BearerAuth
// @Accept param
// @Produce param
// @Param tabUserRoleDto body model.TabUserRoleDto true "用户角色关联实体Dto"
// @Success 200 {object} commonModel.JsonResult "导出用户角色关联列表响应数据"
// @Failure 401 {object} commonModel.JsonResult "未授权"
// @Failure 500 {object} commonModel.JsonResult "服务器异常"
// @Router /tab/user/role/export [get]
func (this *TabUserRoleApi) ExportTabUserRole(c *gin.Context) {
	var tabUserRoleDto dto.TabUserRoleDto
	c.ShouldBindQuery(&tabUserRoleDto)

	res, err := this.tabUserRoleService.ExportTabUserRole(&tabUserRoleDto)

	if err != nil {
		c.JSON(http.StatusInternalServerError, commonModel.JsonResultError("导出用户角色关联列表失败："+err.Error()))
		return
	}
	c.JSON(http.StatusOK, commonModel.JsonResultSuccess[any](res))
}
