package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

// TabRoleApi 角色 api 层
type TabRoleApi struct {
	tabRoleService service.TabRoleService
}

// InsertTabRole 新增角色
// @Tags 角色模块
// @Summary 新增角色-Summary
// @Description 新增角色-Description
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param tabRole body model.TabRole true "修改角色实体类"
// @Success 200 {object} commonModel.JsonResult "新增角色响应数据"
// @Failure 401 {object} commonModel.JsonResult "未授权"
// @Failure 500 {object} commonModel.JsonResult "服务器异常"
// @Router /tab/role [post]
func (this *TabRoleApi) InsertTabRole(c *gin.Context) {
	var tabRole model.TabRole
	c.ShouldBindJSON(&tabRole)

	res, err := this.tabRoleService.InsertTabRole(&tabRole)

	if err != nil {
		c.JSON(http.StatusInternalServerError, commonModel.JsonResultError("新增角色失败："+err.Error()))
		return
	}
	c.JSON(http.StatusOK, commonModel.JsonResultSuccess[any](res))
}

// UpdateTabRole 修改角色
// @Tags 角色模块
// @Summary 修改角色-Summary
// @Description 修改角色-Description
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param tabRole body model.TabRole true "修改角色实体类"
// @Success 200 {object} commonModel.JsonResult "修改角色响应数据"
// @Failure 401 {object} commonModel.JsonResult "未授权"
// @Failure 500 {object} commonModel.JsonResult "服务器异常"
// @Router /tab/role [put]
func (this *TabRoleApi) UpdateTabRole(c *gin.Context) {
	var tabRole model.TabRole
	c.ShouldBindJSON(&tabRole)

	res, err := this.tabRoleService.UpdateTabRole(&tabRole)

	if err != nil {
		c.JSON(http.StatusInternalServerError, commonModel.JsonResultError("修改角色失败："+err.Error()))
		return
	}
	c.JSON(http.StatusOK, commonModel.JsonResultSuccess[any](res))
}

// DeleteTabRole 删除角色
// @Tags 角色模块
// @Summary 删除角色-Summary
// @Description 删除角色-Description
// @Security BearerAuth
// @Accept path
// @Produce path
// @Param idList path string true "角色主键List"
// @Success 200 {object} commonModel.JsonResult "删除角色响应数据"
// @Failure 401 {object} commonModel.JsonResult "未授权"
// @Failure 500 {object} commonModel.JsonResult "服务器异常"
// @Router /tab/role/:idList [delete]
func (this *TabRoleApi) DeleteTabRole(c *gin.Context) {
	idListStr := c.Param("idList") // 例如: "1,2,3"
	idList, err := commonUtils.ParseIds(idListStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, commonModel.JsonResultError("参数错误："+err.Error()))
		return
	}

	res, err := this.tabRoleService.DeleteTabRole(idList)

	if err != nil {
		c.JSON(http.StatusInternalServerError, commonModel.JsonResultError("删除角色失败："+err.Error()))
		return
	}
	c.JSON(http.StatusOK, commonModel.JsonResultSuccess[any](res))
}

// GetTabRoleById 获取角色详细信息
// @Tags 角色模块
// @Summary 获取角色详细信息-Summary
// @Description 获取角色详细信息-Description
// @Security BearerAuth
// @Accept path
// @Produce path
// @Param id path int64 true "角色主键List"
// @Success 200 {object} commonModel.JsonResult "获取角色详细信息"
// @Failure 401 {object} commonModel.JsonResult "未授权"
// @Failure 500 {object} commonModel.JsonResult "服务器异常"
// @Router /tab/role/:id [get]
func (this *TabRoleApi) GetTabRoleById(c *gin.Context) {
	idStr := c.Param("id") // 例如: "1"
	id, err := commonUtils.ParseId(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, commonModel.JsonResultError("参数错误："+err.Error()))
		return
	}

	res, err := this.tabRoleService.GetTabRoleById(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, commonModel.JsonResultError("查询角色失败："+err.Error()))
		return
	}
	c.JSON(http.StatusOK, commonModel.JsonResultSuccess[any](res))
}

// GetTabRoleList 查询角色列表
// @Tags 角色模块
// @Summary 查询角色列表-Summary
// @Description 查询角色列表-Description
// @Security BearerAuth
// @Accept param
// @Produce param
// @Param tabRoleDto body model.TabRoleDto true "角色实体Dto"
// @Success 200 {object} commonModel.JsonResult "查询角色列表响应数据"
// @Failure 401 {object} commonModel.JsonResult "未授权"
// @Failure 500 {object} commonModel.JsonResult "服务器异常"
// @Router /tab/role/list [get]
func (this *TabRoleApi) GetTabRoleList(c *gin.Context) {
	var tabRoleDto dto.TabRoleDto
	c.ShouldBindQuery(&tabRoleDto)

	res, err := this.tabRoleService.GetTabRoleList(&tabRoleDto)

	if err != nil {
		c.JSON(http.StatusInternalServerError, commonModel.JsonResultError("查询角色列表失败："+err.Error()))
		return
	}
	c.JSON(http.StatusOK, commonModel.JsonResultSuccess[any](res))
}

// GetTabRolePageList 分页查询角色列表
// @Tags 角色模块
// @Summary 分页查询角色列表-Summary
// @Description 分页查询角色列表-Description
// @Security BearerAuth
// @Accept param
// @Produce param
// @Param tabRoleDto body model.TabRoleDto true "角色实体Dto"
// @Success 200 {object} commonModel.JsonResult "分页查询角色列表响应数据"
// @Failure 401 {object} commonModel.JsonResult "未授权"
// @Failure 500 {object} commonModel.JsonResult "服务器异常"
// @Router /tab/role/pageList [get]
func (this *TabRoleApi) GetTabRolePageList(c *gin.Context) {
	var tabRoleDto dto.TabRoleDto
	c.ShouldBindQuery(&tabRoleDto)

	res, err := this.tabRoleService.GetTabRolePageList(&tabRoleDto)

	if err != nil {
		c.JSON(http.StatusInternalServerError, commonModel.JsonResultError("查询角色列表失败："+err.Error()))
		return
	}
	c.JSON(http.StatusOK, commonModel.JsonResultSuccess[any](res))
}

// ExportTabRole 导出角色列表
// @Tags 角色模块
// @Summary 导出角色列表-Summary
// @Description 导出角色列表-Description
// @Security BearerAuth
// @Accept param
// @Produce param
// @Param tabRoleDto body model.TabRoleDto true "角色实体Dto"
// @Success 200 {object} commonModel.JsonResult "导出角色列表响应数据"
// @Failure 401 {object} commonModel.JsonResult "未授权"
// @Failure 500 {object} commonModel.JsonResult "服务器异常"
// @Router /tab/role/export [get]
func (this *TabRoleApi) ExportTabRole(c *gin.Context) {
	var tabRoleDto dto.TabRoleDto
	c.ShouldBindQuery(&tabRoleDto)

	res, err := this.tabRoleService.ExportTabRole(&tabRoleDto)

	if err != nil {
		c.JSON(http.StatusInternalServerError, commonModel.JsonResultError("导出角色列表失败："+err.Error()))
		return
	}
	c.JSON(http.StatusOK, commonModel.JsonResultSuccess[any](res))
}
