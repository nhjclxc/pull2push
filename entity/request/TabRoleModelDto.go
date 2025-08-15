package request

import (
	"fmt"
	"github.com/jinzhu/copier"
	"time"
)

// 角色对象 TabRoleDto
// @author
// @date 2025-08-13T10:41:41.967
type TabRoleDto struct {
	model.TabRole

	Keyword string `form:"keyword"` // 模糊搜索字段

	PageNum  int `form:"pageNum"`  // 页码
	PageSize int `form:"pageSize"` // 页大小

	SatrtTime time.Time `form:"satrtTime" time_format:"2006-01-02 15:04:05"` // 开始时间
	EndTime   time.Time `form:"endTime" time_format:"2006-01-02 15:04:05"`   // 结束时间
}

// DtoToModel modelDto 转化为 model
func (this *TabRoleDto) DtoToModel() (tabRole *model.TabRole, err error) {
	// go get github.com/jinzhu/copier

	tabRole = &model.TabRole{} // copier.Copy 不会自动为其分配空间，所以初始化指针指向的结构体
	err = copier.Copy(&tabRole, &this)
	return
}

// ModelToDto model 转化为 modelDto
func (this *TabRoleDto) ModelToDto(tabRole *model.TabRole) error {
	// go get github.com/jinzhu/copier

	err := copier.Copy(&this, &tabRole)
	if err != nil {
		fmt.Printf("DtoTo Copy error: %v", err)
		return err
	}
	return nil
}
