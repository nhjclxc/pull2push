package model

//
//import (
//    "encoding/json"
//    "errors"
//    "fmt"
//    "github.com/mitchellh/mapstructure"
//    "gorm.io/gorm"
//    "time"
//)
//
//// TabRole 角色 结构体
//// @author
//// @date 2025-08-13T10:41:41.967
//type TabRole struct {
//
//
//            Id int64 `gorm:"column:id;primaryKey;auto_increment;not null;type:bigint" json:"id" form:"id"`// 角色ID
//
//            RoleName string `gorm:"column:role_name;not null;type:varchar(50)" json:"roleName" form:"roleName"`// 角色名称
//
//            Description string `gorm:"column:description;type:varchar(255)" json:"description" form:"description"`// 角色描述
//
//            CreatedAt time.Time `gorm:"column:created_at;type:datetime" json:"createdAt" form:"createdAt"`// 创建时间
//
//            UpdatedAt time.Time `gorm:"column:updated_at;type:datetime" json:"updatedAt" form:"updatedAt"`// 更新时间
//
//    // todo update The following predefined fields
//
//    Version uint `gorm:"column:version;default:1" json:"version"` // 乐观锁（版本控制）
//
//    Remark string `gorm:"column:remark;" description:"备注"` // 备注
//
//    DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetime" json:"deletedAt"` // 删除标记, 删除时间 GORM 默认启用了“软删除（Soft Delete）”只要存在这个字段，GORM 默认启用软删除。
//
//    DeletedBy uint64 `gorm:"column:deleted_by;type:bigint" json:"deletedBy"` // 删除人id
//
//    CreatedAt time.Time `gorm:"column:created_at;type:datetime" json:"createdAt"` // 创建时间
//
//    CreatedBy uint64 `gorm:"column:created_by;type:bigint" json:"createdBy"` // 创建人id
//
//    UpdatedAt time.Time `gorm:"column:updated_at;type:datetime" json:"updatedAt"` // 更新时间
//
//    UpdatedBy uint64 `gorm:"column:updated_by;type:bigint" json:"updatedBy"` // 更新人id
//
//    // time_format:"2006-01-02 15:04:05"
//}
//
//// TableName 返回当前实体类的表名
//func (this *TabRole) TableName() string {
//    return "tab_role"
//}
//
//
//// 可用钩子函数包括：
//// BeforeCreate / AfterCreate
//// BeforeUpdate / AfterUpdate
//// BeforeDelete / AfterDelete
//func (this *TabRole) BeforeCreate(tx *gorm.DB) (err error) {
//    this.CreateTime = time.Now()
//    this.UpdateTime = time.Now()
//    return
//}
//
//func (this *TabRole) BeforeUpdate(tx *gorm.DB) (err error) {
//    this.UpdateTime = time.Now()
//    return
//}
//
//// MapToTabRole map映射转化为当前结构体
//func MapToTabRole(inputMap map[string]any) (*TabRole) {
//    //go get github.com/mitchellh/mapstructure
//
//    var tabRole TabRole
//    err := mapstructure.Decode(inputMap, &tabRole)
//    if err != nil {
//        fmt.Printf("MapToStruct Decode error: %v", err)
//        return nil
//    }
//    return &tabRole
//}
//
//// TabRoleToMap 当前结构体转化为map映射
//func (this *TabRole) TabRoleToMap() (map[string]any) {
//    var m map[string]any
//    bytes, err := json.Marshal(this)
//    if err != nil {
//        fmt.Printf("StructToMap marshal error: %v", err)
//        return nil
//    }
//
//    err = json.Unmarshal(bytes, &m)
//    if err != nil {
//        fmt.Printf("StructToMap unmarshal error: %v", err)
//        return nil
//    }
//    return m
//}
//
//
//
//// 由于有时需要开启事务，因此 DB *gorm.DB 选择从外部传入
//
//// InsertTabRole 新增角色
//func (this *TabRole) InsertTabRole(DB *gorm.DB) (int, error) {
//    fmt.Printf("InsertTabRole：%#v \n", this)
//
//    // 先查询是否有相同 name 的数据存在
//    temp := &TabRole{}
//    // todo update name
//    tx := DB.Where("name = ?", this.?).First(temp)
//    fmt.Printf("InsertTabRole.Where：%#v \n", temp)
//    if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
//        return 0, errors.New("InsertTabRole.Where, 存在相同name: " + temp.?)
//    }
//
//    // 执行 Insert
//    err := DB.Create(&this).Error
//
//    if err != nil {
//        return 0, errors.New("InsertTabRole.DB.Create, 新增失败: " + err.Error())
//    }
//    return 1, nil
//}
//
//// BatchInsertTabRoles 批量新增角色
//func (this *TabRole) BatchInsertTabRoles(DB *gorm.DB, tables []TabRole) (int, error) {
//
//    result := DB.Create(&tables)
//
//    if result.Error != nil {
//        return 0, errors.New("BatchInsertTabRoles.DB.Create, 新增失败: " + result.Error.Error())
//    }
//    return int(result.RowsAffected), nil
//}
//
//// UpdateTabRoleById 根据主键修改角色的所有字段
//func (this *TabRole) UpdateTabRoleById(DB *gorm.DB) (int, error) {
//    fmt.Printf("UpdateTabRoleById：%#v \n", this)
//
//    // 1、查询该id是否存在
//    if this.Id == 0 {
//        return 0, errors.New("Id 不能为空！！！: ")
//    }
//
//    // 2、再看看name是否重复
//    temp := &TabRole{}
//    // todo update name
//    tx := DB.Where("name = ?", this.?).First(temp)
//    fmt.Printf("UpdateTabRoleById.Where：%#v \n", temp)
//    if !errors.Is(tx.Error, gorm.ErrRecordNotFound) && temp.Id != this.Id {
//        return 0, errors.New("UpdateTabRoleById.Where, 存在相同name: " + temp.?)
//    }
//
//    // 3、执行修改
//    //保存整个结构体（全字段更新）
//    saveErr := DB.Save(this).Error
//    if saveErr != nil {
//        return 0, errors.New("UpdateTabRoleById.Save, 修改失败: " + saveErr.Error())
//    }
//    return 1, nil
//}
//
//// UpdateTabRoleSelective 修改角色不为默认值的字段
//func (this *TabRole) UpdateTabRoleSelective(DB *gorm.DB) (int, error) {
//    fmt.Printf("UpdateTabRoleSelective：%#v \n", this)
//
//    // db.Model().Updates()：只更新指定字段
//    err := DB.Model(this).
//        Where("id = ?", this.Id).
//        Updates(this).
//        Error
//    if err != nil {
//        return 0, errors.New("UpdateTabRoleSelective.Updates, 选择性修改失败: " + err.Error())
//    }
//
//    return 1, nil
//}
//
//// DeleteTabRole 删除角色
//func (this *TabRole) DeleteTabRole(DB *gorm.DB, idList []int64) (int, error) {
//    fmt.Printf("DeleteTabRole：%#v \n", idList)
//
//    // 当存在DeletedAt gorm.DeletedAt字段时为软删除，否则为物理删除
//    result := DB.Delete(&this, "id in ?", idList)
//    // result := DB.Model(&this).Where("id IN ?", tableIdList).Update("state", 0)
//    if result.Error != nil {
//        return 0, errors.New("DeleteTabRole.Delete, 删除失败: " + result.Error.Error())
//    }
//
//    //// 以下使用的是物理删除
//    //result := DB.Unscoped().Delete(this, "id in ?", idList)
//    //if result.Error != nil {
//    //	return 0, errors.New("DeleteTabRole.Delete, 删除失败: " + result.Error.Error())
//    //}
//
//    return int(result.RowsAffected), nil
//}
//
//// BatchDeleteTabRoles 根据主键批量删除角色
//func (this *TabRole) BatchDeleteTabRoles(DB *gorm.DB, idList []int64) error {
//    return DB.Where("id IN ?", idList).Delete(&this).Error
//}
//
//// FindTabRoleById 获取角色详细信息
//func (this *TabRole) FindTabRoleById(DB *gorm.DB, id int64) (error) {
//    fmt.Printf("DeleteTabRole：%#v \n", id)
//    return DB.First(this, "id = ?", id).Error
//}
//
//// FindTabRolesByIdList 根据主键批量查询角色详细信息
//func FindTabRolesByIdList(DB *gorm.DB, idList []int64) ([]TabRole, error) {
//    var result []TabRole
//    err := DB.Where("id IN ?", idList).Find(&result).Error
//    return result, err
//}
//
//// FindTabRoleList 查询角色列表
//func (this *TabRole) FindTabRoleList(DB *gorm.DB, satrtTime time.Time, endTime time.Time) ([]TabRole, error) {
//    fmt.Printf("GetTabRoleList：%#v \n", this)
//
//    var tables []TabRole
//    query := DB.Model(this)
//
//        // 构造查询条件
//        if this.Id != 0 { query = query.Where("id = ?", this.Id ) }
//        if this.RoleName != "" { query = query.Where("role_name LIKE ?", "%" + this.RoleName + "%") }
//        if this.Description != "" { query = query.Where("description LIKE ?", "%" + this.Description + "%") }
//        if !this.CreatedAt.IsZero() {
//            query = query.Where("created_at = ?", this.CreatedAt)
//            // query = query.Where("DATE(created_at) = ?", this.$column.goField.Format("2006-01-02"))
//        }
//        if !this.UpdatedAt.IsZero() {
//            query = query.Where("updated_at = ?", this.UpdatedAt)
//            // query = query.Where("DATE(updated_at) = ?", this.$column.goField.Format("2006-01-02"))
//        }
//
//    if !satrtTime.IsZero() {
//        query = query.Where("create_time >= ?", satrtTime)
//    }
//    if !endTime.IsZero() {
//        query = query.Where("create_time <= ?", endTime)
//    }
//
//    // // 添加分页逻辑
//    // if tabRole.PageNum > 0 && tabRole.PageSize > 0 {
//    //     offset := (tabRole.PageNum - 1) * tabRole.PageSize
//    //     query = query.Offset(offset).Limit(tabRole.PageSize)
//    // }
//
//    err := query.Find(&tables).Error
//    return tables, err
//}
//
//// FindTabRolePageList 分页查询角色列表
//func (this *TabRole) FindTabRolePageList(DB *gorm.DB, satrtTime time.Time, endTime time.Time, pageNum int, pageSize int) ([]TabRole, int64, error) {
//    fmt.Printf("GetTabRolePageList：%#v \n", this)
//
//    var (
//        tabRoles []TabRole
//        total     int64
//    )
//
//    query := DB.Model(&TabRole{})
//
//// 构造查询条件
//        if this.Id != 0 { query = query.Where("id = ?", this.Id ) }
//        if this.RoleName != "" { query = query.Where("role_name LIKE ?", "%" + this.RoleName + "%") }
//        if this.Description != "" { query = query.Where("description LIKE ?", "%" + this.Description + "%") }
//        if !this.CreatedAt.IsZero() {
//            query = query.Where("created_at = ?", this.CreatedAt)
//            // query = query.Where("DATE(created_at) = ?", this.$column.goField.Format("2006-01-02"))
//        }
//        if !this.UpdatedAt.IsZero() {
//            query = query.Where("updated_at = ?", this.UpdatedAt)
//            // query = query.Where("DATE(updated_at) = ?", this.$column.goField.Format("2006-01-02"))
//        }
//
//    if !satrtTime.IsZero() {
//        query = query.Where("create_time >= ?", satrtTime)
//    }
//    if !endTime.IsZero() {
//        query = query.Where("create_time <= ?", endTime)
//    }
//
//    // 查询总数
//    if err := query.Count(&total).Error; err != nil {
//        return nil, 0, err
//    }
//
//    // 分页参数默认值
//    if pageNum <= 0 {
//        pageNum = 1
//    }
//    if pageSize <= 0 {
//        pageSize = 10
//    }
//
//    // 分页数据
//    // todo update create_time
//    err := query.
//        Limit(pageSize).Offset((pageNum - 1) * pageSize).
//        Order("create_time desc").
//        Find(&tabRoles).Error
//
//    if err != nil {
//        return nil, 0, err
//    }
//
//    return tabRoles, total, nil
//}
//
