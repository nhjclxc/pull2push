package service

import "pull2push/entity/model"

/*

-- 用户表
CREATE TABLE tab_user (
id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '用户ID',
username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
password VARCHAR(255) NOT NULL COMMENT '密码（加密存储）',
email VARCHAR(100) DEFAULT NULL COMMENT '邮箱',
created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
*/
// TabUserService 用户 Service 层
type TabUserService struct {
}

// InsertTabUser 新增用户
func (this *TabUserService) InsertTabUser(tabUser *model.TabUser) (res any, err error) {

	return tabUser.InsertTabUser(core.GLOBAL_DB)
}

// UpdateTabUser 修改用户
func (this *TabUserService) UpdateTabUser(tabUser *model.TabUser) (res any, err error) {

	return tabUser.UpdateTabUserById(core.GLOBAL_DB)
}

// DeleteTabUser 删除用户
func (this *TabUserService) DeleteTabUser(idList []int64) (res any, err error) {

	return (&model.TabUser{}).DeleteTabUser(core.GLOBAL_DB, idList)
}

// GetTabUserById 获取用户业务详细信息
func (this *TabUserService) GetTabUserById(id int64) (res any, err error) {

	tabUser := model.TabUser{}
	err = (&tabUser).FindTabUserById(core.GLOBAL_DB, id)
	if err != nil {
		return nil, err
	}

	return tabUser, nil
}

// GetTabUserList 查询用户业务列表
func (this *TabUserService) GetTabUserList(tabUserDto *dto.TabUserDto) (res any, err error) {

	tabUser, err := tabUserDto.DtoToModel()
	tabUserList, err := tabUser.FindTabUserList(core.GLOBAL_DB, tabUserDto.SatrtTime, tabUserDto.EndTime)
	if err != nil {
		return nil, err
	}

	return tabUserList, nil
}

// GetTabUserPageList 分页查询用户业务列表
func (this *TabUserService) GetTabUserPageList(tabUserDto *dto.TabUserDto) (res any, err error) {

	tabUser, err := tabUserDto.DtoToModel()
	tabUserList, total, err := tabUser.FindTabUserPageList(core.GLOBAL_DB, tabUserDto.SatrtTime, tabUserDto.EndTime, tabUserDto.PageNum, tabUserDto.PageSize)
	if err != nil {
		return nil, err
	}

	return commonUtils.BuildPageData[model.TabUser](tabUserList, total, tabUserDto.PageNum, tabUserDto.PageSize), nil
}

// ExportTabUser 导出用户业务列表
func (this *TabUserService) ExportTabUser(tabUserDto *dto.TabUserDto) (res any, err error) {

	tabUser, err := tabUserDto.DtoToModel()
	tabUser.FindTabUserPageList(core.GLOBAL_DB, tabUserDto.SatrtTime, tabUserDto.EndTime, 1, 10000)
	// 实现导出 ...

	return nil, nil
}
