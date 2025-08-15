package base

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// 项目统一返回结构
type JsonResult[T any] struct {
	Success bool `json:"success"` // 是否成功 true or false

	Code int `json:"code"` // 状态码

	Msg string `json:"msg"` // 返回内容

	Data T `json:"data"` // 数据对象
}

// Success 初始化一个新创建的 AjaxResult 对象
// @param type 状态类型
// @param msg 返回内容
// @param data 数据对象
func JsonResultSuccess[T any](data T) *JsonResult[T] {
	return &JsonResult[T]{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "操作成功",
		Data:    data,
	}
}

// Error 返回错误消息
// @param msg 返回内容
func JsonResultError(msg string) *JsonResult[string] {
	return &JsonResult[string]{
		Success: true,
		Code:    http.StatusOK,
		Msg:     msg,
		Data:    "",
	}
}

// ParseId 将 string 类型数据转化为 int64 类型数据
func ParseId(idStr string) (int64, error) {
	return strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
}

// ParseIds 将形如 "1,2,3" 的字符串解析为 []int64
func ParseIds(idListStr string) ([]int64, error) {
	if idListStr == "" {
		return nil, errors.New("idListStr 为空")
	}

	parts := strings.Split(idListStr, ",")
	ids := make([]int64, 0, len(parts))

	for _, part := range parts {
		id, err := ParseId(part)
		if err != nil {
			return nil, fmt.Errorf("无法解析 ID '%s'", part)
		}
		ids = append(ids, id)
	}

	return ids, nil
}

// BuildPageData 构造分页返回结构
func BuildPageData[T any](dataList []T, total int64, pageNum, pageSize int) map[string]any {
	if pageSize <= 0 {
		pageSize = 10 // 默认每页 10 条，避免除零
	}
	if pageNum <= 0 {
		pageNum = 1 // 默认第一页
	}

	pages := int((total + int64(pageSize) - 1) / int64(pageSize))
	size := len(dataList)

	startRow := 0
	endRow := 0
	if size > 0 {
		startRow = (pageNum-1)*pageSize + 1
		endRow = startRow + size - 1
	}

	prePage := 0
	if pageNum > 1 {
		prePage = pageNum - 1
	}

	nextPage := 0
	if pageNum < pages {
		nextPage = pageNum + 1
	}

	isFirstPage := pageNum == 1
	isLastPage := pageNum == pages || pages == 0
	hasPreviousPage := pageNum > 1
	hasNextPage := pageNum < pages

	navigatePages := 10
	navigatepageNums := []int{}
	startNav := pageNum - navigatePages/2
	if startNav < 1 {
		startNav = 1
	}
	endNav := startNav + navigatePages - 1
	if endNav > pages {
		endNav = pages
	}

	for i := startNav; i <= endNav; i++ {
		navigatepageNums = append(navigatepageNums, i)
	}

	navigateFirstPage := 0
	if len(navigatepageNums) > 0 {
		navigateFirstPage = navigatepageNums[0]
	}

	navigateLastPage := 0
	if len(navigatepageNums) > 0 {
		navigateLastPage = navigatepageNums[len(navigatepageNums)-1]
	}

	return map[string]any{
		"total":             total,
		"list":              dataList,
		"pageNum":           pageNum,
		"pageSize":          pageSize,
		"size":              size,
		"startRow":          startRow,
		"endRow":            endRow,
		"pages":             pages,
		"prePage":           prePage,
		"nextPage":          nextPage,
		"isFirstPage":       isFirstPage,
		"isLastPage":        isLastPage,
		"hasPreviousPage":   hasPreviousPage,
		"hasNextPage":       hasNextPage,
		"navigatePages":     navigatePages,
		"navigatepageNums":  navigatepageNums,
		"navigateFirstPage": navigateFirstPage,
		"navigateLastPage":  navigateLastPage,
	}
}
