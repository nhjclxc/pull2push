package response

import (
	"fmt"
	"github.com/jinzhu/copier"
)

// 角色对象 TabRoleVo
// @author
// @date 2025-08-13T10:41:41.967
type TabRoleVo struct {
	model.TabRole

	Foo string `form:"foo"` // foo
	Bar string `form:"bar"` // bar
	// ...
}

// ModelToVo model 转化为 modelVo
func (this *TabRoleVo) ModelToVo(tabRole *model.TabRole) error {
	// go get github.com/jinzhu/copier

	tabRole = &model.TabRole{} // copier.Copy 不会自动为其分配空间，所以初始化指针指向的结构体
	err := copier.Copy(&this, &tabRole)
	if err != nil {
		fmt.Printf("ModelToVo Copy error: %v", err)
		return err
	}
	return nil
}
