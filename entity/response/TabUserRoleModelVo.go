package response

import (
	"fmt"
	"github.com/jinzhu/copier"
)

// 用户角色关联对象 TabUserRoleVo
// @author
// @date 2025-08-13T10:41:42.054
type TabUserRoleVo struct {
	model.TabUserRole

	Foo string `form:"foo"` // foo
	Bar string `form:"bar"` // bar
	// ...
}

// ModelToVo model 转化为 modelVo
func (this *TabUserRoleVo) ModelToVo(tabUserRole *model.TabUserRole) error {
	// go get github.com/jinzhu/copier

	tabUserRole = &model.TabUserRole{} // copier.Copy 不会自动为其分配空间，所以初始化指针指向的结构体
	err := copier.Copy(&this, &tabUserRole)
	if err != nil {
		fmt.Printf("ModelToVo Copy error: %v", err)
		return err
	}
	return nil
}
