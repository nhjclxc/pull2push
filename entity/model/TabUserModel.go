package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
	"time"
)

// TabUser 用户 结构体
// @author
// @date 2025-08-13T10:41:41.793
type TabUser struct {
	Id int64 `gorm:"column:id;primaryKey;auto_increment;not null;type:bigint" json:"id" form:"id"` // 用户ID

	Username string `gorm:"column:username;not null;type:varchar(50)" json:"username" form:"username"` // 用户名

	Password string `gorm:"column:password;not null;type:varchar(255)" json:"password" form:"password"` // 密码（加密存储）

	Email string `gorm:"column:email;type:varchar(100)" json:"email" form:"email"` // 邮箱

	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetime" json:"deletedAt"` // 删除标记, 删除时间 GORM 默认启用了“软删除（Soft Delete）”只要存在这个字段，GORM 默认启用软删除。

	DeletedBy uint64 `gorm:"column:deleted_by;type:bigint" json:"deletedBy"` // 删除人id

	CreatedAt time.Time `gorm:"column:created_at;type:datetime" json:"createdAt"` // 创建时间

	CreatedBy uint64 `gorm:"column:created_by;type:bigint" json:"createdBy"` // 创建人id

	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime" json:"updatedAt"` // 更新时间

	UpdatedBy uint64 `gorm:"column:updated_by;type:bigint" json:"updatedBy"` // 更新人id

	// time_format:"2006-01-02 15:04:05"
}

// TableName 返回当前实体类的表名
func (this *TabUser) TableName() string {
	return "tab_user"
}

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

// 可用钩子函数包括：
// BeforeCreate / AfterCreate
// BeforeUpdate / AfterUpdate
// BeforeDelete / AfterDelete
func (this *TabUser) BeforeCreate(tx *gorm.DB) (err error) {
	this.CreatedAt = time.Now()
	this.UpdatedAt = time.Now()
	return
}

func (this *TabUser) BeforeUpdate(tx *gorm.DB) (err error) {
	this.UpdatedAt = time.Now()
	return
}

// MapToTabUser map映射转化为当前结构体
func MapToTabUser(inputMap map[string]any) *TabUser {
	//go get github.com/mitchellh/mapstructure

	var tabUser TabUser
	err := mapstructure.Decode(inputMap, &tabUser)
	if err != nil {
		fmt.Printf("MapToStruct Decode error: %v", err)
		return nil
	}
	return &tabUser
}

// TabUserToMap 当前结构体转化为map映射
func (this *TabUser) TabUserToMap() map[string]any {
	var m map[string]any
	bytes, err := json.Marshal(this)
	if err != nil {
		fmt.Printf("StructToMap marshal error: %v", err)
		return nil
	}

	err = json.Unmarshal(bytes, &m)
	if err != nil {
		fmt.Printf("StructToMap unmarshal error: %v", err)
		return nil
	}
	return m
}

// 由于有时需要开启事务，因此 DB *gorm.DB 选择从外部传入

// InsertTabUser 新增用户
func (this *TabUser) InsertTabUser(DB *gorm.DB) (int, error) {
	fmt.Printf("InsertTabUser：%#v \n", this)

	// 先查询是否有相同 name 的数据存在
	temp := &TabUser{}
	tx := DB.Where("username = ?", this.Username).First(temp)
	fmt.Printf("InsertTabUser.Where：%#v \n", temp)
	if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return 0, errors.New("InsertTabUser.Where, 存在相同username: " + temp.Username)
	}

	// 执行 Insert
	err := DB.Create(&this).Error

	if err != nil {
		return 0, errors.New("InsertTabUser.DB.Create, 新增失败: " + err.Error())
	}
	return 1, nil
}

// BatchInsertTabUsers 批量新增用户
func (this *TabUser) BatchInsertTabUsers(DB *gorm.DB, tables []TabUser) (int, error) {

	result := DB.Create(&tables)

	if result.Error != nil {
		return 0, errors.New("BatchInsertTabUsers.DB.Create, 新增失败: " + result.Error.Error())
	}
	return int(result.RowsAffected), nil
}

// UpdateTabUserById 根据主键修改用户的所有字段
func (this *TabUser) UpdateTabUserById(DB *gorm.DB) (int, error) {
	fmt.Printf("UpdateTabUserById：%#v \n", this)

	// 1、查询该id是否存在
	if this.Id == 0 {
		return 0, errors.New("Id 不能为空！！！: ")
	}

	// 2、再看看name是否重复
	temp := &TabUser{}
	tx := DB.Where("username = ?", this.Username).First(temp)
	fmt.Printf("UpdateTabUserById.Where：%#v \n", temp)
	if !errors.Is(tx.Error, gorm.ErrRecordNotFound) && temp.Id != this.Id {
		return 0, errors.New("UpdateTabUserById.Where, 存在相同 username: " + temp.Username)
	}

	// 3、执行修改
	//保存整个结构体（全字段更新）
	saveErr := DB.Save(this).Error
	if saveErr != nil {
		return 0, errors.New("UpdateTabUserById.Save, 修改失败: " + saveErr.Error())
	}
	return 1, nil
}

// UpdateTabUserSelective 修改用户不为默认值的字段
func (this *TabUser) UpdateTabUserSelective(DB *gorm.DB) (int, error) {
	fmt.Printf("UpdateTabUserSelective：%#v \n", this)

	// db.Model().Updates()：只更新指定字段
	err := DB.Model(this).
		Where("id = ?", this.Id).
		Updates(this).
		Error
	if err != nil {
		return 0, errors.New("UpdateTabUserSelective.Updates, 选择性修改失败: " + err.Error())
	}

	return 1, nil
}

// DeleteTabUser 删除用户
func (this *TabUser) DeleteTabUser(DB *gorm.DB, idList []int64) (int, error) {
	fmt.Printf("DeleteTabUser：%#v \n", idList)

	// 当存在DeletedAt gorm.DeletedAt字段时为软删除，否则为物理删除
	result := DB.Delete(&this, "id in ?", idList)
	// result := DB.Model(&this).Where("id IN ?", tableIdList).Update("state", 0)
	if result.Error != nil {
		return 0, errors.New("DeleteTabUser.Delete, 删除失败: " + result.Error.Error())
	}

	//// 以下使用的是物理删除
	//result := DB.Unscoped().Delete(this, "id in ?", idList)
	//if result.Error != nil {
	//	return 0, errors.New("DeleteTabUser.Delete, 删除失败: " + result.Error.Error())
	//}

	return int(result.RowsAffected), nil
}

// BatchDeleteTabUsers 根据主键批量删除用户
func (this *TabUser) BatchDeleteTabUsers(DB *gorm.DB, idList []int64) error {
	return DB.Where("id IN ?", idList).Delete(&this).Error
}

// FindTabUserById 获取用户详细信息
func (this *TabUser) FindTabUserById(DB *gorm.DB, id int64) error {
	fmt.Printf("DeleteTabUser：%#v \n", id)
	return DB.First(this, "id = ?", id).Error
}

// FindTabUsersByIdList 根据主键批量查询用户详细信息
func FindTabUsersByIdList(DB *gorm.DB, idList []int64) ([]TabUser, error) {
	var result []TabUser
	err := DB.Where("id IN ?", idList).Find(&result).Error
	return result, err
}

// FindTabUserList 查询用户列表
func (this *TabUser) FindTabUserList(DB *gorm.DB, satrtTime time.Time, endTime time.Time) ([]TabUser, error) {
	fmt.Printf("GetTabUserList：%#v \n", this)

	var tables []TabUser
	query := DB.Model(this)

	// 构造查询条件
	if this.Id != 0 {
		query = query.Where("id = ?", this.Id)
	}
	if this.Username != "" {
		query = query.Where("username LIKE ?", "%"+this.Username+"%")
	}
	if this.Password != "" {
		query = query.Where("password = ?", this.Password)
	}
	if this.Email != "" {
		query = query.Where("email = ?", this.Email)
	}
	if !this.CreatedAt.IsZero() {
		query = query.Where("created_at = ?", this.CreatedAt)
		// query = query.Where("DATE(created_at) = ?", this.$column.goField.Format("2006-01-02"))
	}
	if !this.UpdatedAt.IsZero() {
		query = query.Where("updated_at = ?", this.UpdatedAt)
		// query = query.Where("DATE(updated_at) = ?", this.$column.goField.Format("2006-01-02"))
	}

	if !satrtTime.IsZero() {
		query = query.Where("create_at >= ?", satrtTime)
	}
	if !endTime.IsZero() {
		query = query.Where("create_at <= ?", endTime)
	}

	// // 添加分页逻辑
	// if tabUser.PageNum > 0 && tabUser.PageSize > 0 {
	//     offset := (tabUser.PageNum - 1) * tabUser.PageSize
	//     query = query.Offset(offset).Limit(tabUser.PageSize)
	// }

	err := query.Find(&tables).Error
	return tables, err
}

// FindTabUserPageList 分页查询用户列表
func (this *TabUser) FindTabUserPageList(DB *gorm.DB, satrtTime time.Time, endTime time.Time, pageNum int, pageSize int) ([]TabUser, int64, error) {
	fmt.Printf("GetTabUserPageList：%#v \n", this)

	var (
		tabUsers []TabUser
		total    int64
	)

	query := DB.Model(&TabUser{})

	// 构造查询条件
	if this.Id != 0 {
		query = query.Where("id = ?", this.Id)
	}
	if this.Username != "" {
		query = query.Where("username LIKE ?", "%"+this.Username+"%")
	}
	if this.Password != "" {
		query = query.Where("password = ?", this.Password)
	}
	if this.Email != "" {
		query = query.Where("email = ?", this.Email)
	}
	if !this.CreatedAt.IsZero() {
		query = query.Where("created_at = ?", this.CreatedAt)
		// query = query.Where("DATE(created_at) = ?", this.$column.goField.Format("2006-01-02"))
	}
	if !this.UpdatedAt.IsZero() {
		query = query.Where("updated_at = ?", this.UpdatedAt)
		// query = query.Where("DATE(updated_at) = ?", this.$column.goField.Format("2006-01-02"))
	}

	if !satrtTime.IsZero() {
		query = query.Where("create_at >= ?", satrtTime)
	}
	if !endTime.IsZero() {
		query = query.Where("create_at <= ?", endTime)
	}

	// 查询总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页参数默认值
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// 分页数据
	err := query.
		Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Order("create_at desc").
		Find(&tabUsers).Error

	if err != nil {
		return nil, 0, err
	}

	return tabUsers, total, nil
}
