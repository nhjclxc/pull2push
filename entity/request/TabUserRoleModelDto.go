package request

import (
	"fmt"
	"github.com/jinzhu/copier"
	"time"
)

// 用户角色关联对象 TabUserRoleDto
// @author
// @date 2025-08-13T10:41:42.054
type TabUserRoleDto struct {
	model.TabUserRole

	Keyword string `form:"keyword"` // 模糊搜索字段

	PageNum  int `form:"pageNum"`  // 页码
	PageSize int `form:"pageSize"` // 页大小

	SatrtTime time.Time `form:"satrtTime" time_format:"2006-01-02 15:04:05"` // 开始时间
	EndTime   time.Time `form:"endTime" time_format:"2006-01-02 15:04:05"`   // 结束时间
}

// DtoToModel modelDto 转化为 model
func (this *TabUserRoleDto) DtoToModel() (tabUserRole *model.TabUserRole, err error) {
	// go get github.com/jinzhu/copier

	tabUserRole = &model.TabUserRole{} // copier.Copy 不会自动为其分配空间，所以初始化指针指向的结构体
	err = copier.Copy(&tabUserRole, &this)
	return
}

// ModelToDto model 转化为 modelDto
func (this *TabUserRoleDto) ModelToDto(tabUserRole *model.TabUserRole) error {
	// go get github.com/jinzhu/copier

	err := copier.Copy(&this, &tabUserRole)
	if err != nil {
		fmt.Printf("DtoTo Copy error: %v", err)
		return err
	}
	return nil
}
