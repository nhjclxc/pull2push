package response

import (
	"fmt"
	"github.com/jinzhu/copier"
)

// 用户对象 TabUserVo
// @author
// @date 2025-08-13T10:41:41.793
type TabUserVo struct {
	model.TabUser

	Foo string `form:"foo"` // foo
	Bar string `form:"bar"` // bar
	// ...
}

// ModelToVo model 转化为 modelVo
func (this *TabUserVo) ModelToVo(tabUser *model.TabUser) error {
	// go get github.com/jinzhu/copier

	tabUser = &model.TabUser{} // copier.Copy 不会自动为其分配空间，所以初始化指针指向的结构体
	err := copier.Copy(&this, &tabUser)
	if err != nil {
		fmt.Printf("ModelToVo Copy error: %v", err)
		return err
	}
	return nil
}
