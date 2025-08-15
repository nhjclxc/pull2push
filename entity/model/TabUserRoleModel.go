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
//// TabUserRole 用户角色关联 结构体
//// @author
//// @date 2025-08-13T10:41:42.054
//type TabUserRole struct {
//
//
//            Id int64 `gorm:"column:id;primaryKey;auto_increment;not null;type:bigint" json:"id" form:"id"`// 主键ID
//
//            UserId int64 `gorm:"column:user_id;not null;type:bigint" json:"userId" form:"userId"`// 用户ID
//
//            RoleId int64 `gorm:"column:role_id;not null;type:bigint" json:"roleId" form:"roleId"`// 角色ID
//
//            CreatedAt time.Time `gorm:"column:created_at;type:datetime" json:"createdAt" form:"createdAt"`// 关联创建时间
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
//func (this *TabUserRole) TableName() string {
//    return "tab_user_role"
//}
//
//
//// 可用钩子函数包括：
//// BeforeCreate / AfterCreate
//// BeforeUpdate / AfterUpdate
//// BeforeDelete / AfterDelete
//func (this *TabUserRole) BeforeCreate(tx *gorm.DB) (err error) {
//    this.CreateTime = time.Now()
//    this.UpdateTime = time.Now()
//    return
//}
//
//func (this *TabUserRole) BeforeUpdate(tx *gorm.DB) (err error) {
//    this.UpdateTime = time.Now()
//    return
//}
//
//// MapToTabUserRole map映射转化为当前结构体
//func MapToTabUserRole(inputMap map[string]any) (*TabUserRole) {
//    //go get github.com/mitchellh/mapstructure
//
//    var tabUserRole TabUserRole
//    err := mapstructure.Decode(inputMap, &tabUserRole)
//    if err != nil {
//        fmt.Printf("MapToStruct Decode error: %v", err)
//        return nil
//    }
//    return &tabUserRole
//}
//
//// TabUserRoleToMap 当前结构体转化为map映射
//func (this *TabUserRole) TabUserRoleToMap() (map[string]any) {
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
//// InsertTabUserRole 新增用户角色关联
//func (this *TabUserRole) InsertTabUserRole(DB *gorm.DB) (int, error) {
//    fmt.Printf("InsertTabUserRole：%#v \n", this)
//
//    // 先查询是否有相同 name 的数据存在
//    temp := &TabUserRole{}
//    // todo update name
//    tx := DB.Where("name = ?", this.?).First(temp)
//    fmt.Printf("InsertTabUserRole.Where：%#v \n", temp)
//    if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
//        return 0, errors.New("InsertTabUserRole.Where, 存在相同name: " + temp.?)
//    }
//
//    // 执行 Insert
//    err := DB.Create(&this).Error
//
//    if err != nil {
//        return 0, errors.New("InsertTabUserRole.DB.Create, 新增失败: " + err.Error())
//    }
//    return 1, nil
//}
//
//// BatchInsertTabUserRoles 批量新增用户角色关联
//func (this *TabUserRole) BatchInsertTabUserRoles(DB *gorm.DB, tables []TabUserRole) (int, error) {
//
//    result := DB.Create(&tables)
//
//    if result.Error != nil {
//        return 0, errors.New("BatchInsertTabUserRoles.DB.Create, 新增失败: " + result.Error.Error())
//    }
//    return int(result.RowsAffected), nil
//}
//
//// UpdateTabUserRoleById 根据主键修改用户角色关联的所有字段
//func (this *TabUserRole) UpdateTabUserRoleById(DB *gorm.DB) (int, error) {
//    fmt.Printf("UpdateTabUserRoleById：%#v \n", this)
//
//    // 1、查询该id是否存在
//    if this.Id == 0 {
//        return 0, errors.New("Id 不能为空！！！: ")
//    }
//
//    // 2、再看看name是否重复
//    temp := &TabUserRole{}
//    // todo update name
//    tx := DB.Where("name = ?", this.?).First(temp)
//    fmt.Printf("UpdateTabUserRoleById.Where：%#v \n", temp)
//    if !errors.Is(tx.Error, gorm.ErrRecordNotFound) && temp.Id != this.Id {
//        return 0, errors.New("UpdateTabUserRoleById.Where, 存在相同name: " + temp.?)
//    }
//
//    // 3、执行修改
//    //保存整个结构体（全字段更新）
//    saveErr := DB.Save(this).Error
//    if saveErr != nil {
//        return 0, errors.New("UpdateTabUserRoleById.Save, 修改失败: " + saveErr.Error())
//    }
//    return 1, nil
//}
//
//// UpdateTabUserRoleSelective 修改用户角色关联不为默认值的字段
//func (this *TabUserRole) UpdateTabUserRoleSelective(DB *gorm.DB) (int, error) {
//    fmt.Printf("UpdateTabUserRoleSelective：%#v \n", this)
//
//    // db.Model().Updates()：只更新指定字段
//    err := DB.Model(this).
//        Where("id = ?", this.Id).
//        Updates(this).
//        Error
//    if err != nil {
//        return 0, errors.New("UpdateTabUserRoleSelective.Updates, 选择性修改失败: " + err.Error())
//    }
//
//    return 1, nil
//}
//
//// DeleteTabUserRole 删除用户角色关联
//func (this *TabUserRole) DeleteTabUserRole(DB *gorm.DB, idList []int64) (int, error) {
//    fmt.Printf("DeleteTabUserRole：%#v \n", idList)
//
//    // 当存在DeletedAt gorm.DeletedAt字段时为软删除，否则为物理删除
//    result := DB.Delete(&this, "id in ?", idList)
//    // result := DB.Model(&this).Where("id IN ?", tableIdList).Update("state", 0)
//    if result.Error != nil {
//        return 0, errors.New("DeleteTabUserRole.Delete, 删除失败: " + result.Error.Error())
//    }
//
//    //// 以下使用的是物理删除
//    //result := DB.Unscoped().Delete(this, "id in ?", idList)
//    //if result.Error != nil {
//    //	return 0, errors.New("DeleteTabUserRole.Delete, 删除失败: " + result.Error.Error())
//    //}
//
//    return int(result.RowsAffected), nil
//}
//
//// BatchDeleteTabUserRoles 根据主键批量删除用户角色关联
//func (this *TabUserRole) BatchDeleteTabUserRoles(DB *gorm.DB, idList []int64) error {
//    return DB.Where("id IN ?", idList).Delete(&this).Error
//}
//
//// FindTabUserRoleById 获取用户角色关联详细信息
//func (this *TabUserRole) FindTabUserRoleById(DB *gorm.DB, id int64) (error) {
//    fmt.Printf("DeleteTabUserRole：%#v \n", id)
//    return DB.First(this, "id = ?", id).Error
//}
//
//// FindTabUserRolesByIdList 根据主键批量查询用户角色关联详细信息
//func FindTabUserRolesByIdList(DB *gorm.DB, idList []int64) ([]TabUserRole, error) {
//    var result []TabUserRole
//    err := DB.Where("id IN ?", idList).Find(&result).Error
//    return result, err
//}
//
//// FindTabUserRoleList 查询用户角色关联列表
//func (this *TabUserRole) FindTabUserRoleList(DB *gorm.DB, satrtTime time.Time, endTime time.Time) ([]TabUserRole, error) {
//    fmt.Printf("GetTabUserRoleList：%#v \n", this)
//
//    var tables []TabUserRole
//    query := DB.Model(this)
//
//        // 构造查询条件
//        if this.Id != 0 { query = query.Where("id = ?", this.Id ) }
//        if this.UserId != 0 { query = query.Where("user_id = ?", this.UserId ) }
//        if this.RoleId != 0 { query = query.Where("role_id = ?", this.RoleId ) }
//        if !this.CreatedAt.IsZero() {
//            query = query.Where("created_at = ?", this.CreatedAt)
//            // query = query.Where("DATE(created_at) = ?", this.$column.goField.Format("2006-01-02"))
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
//    // if tabUserRole.PageNum > 0 && tabUserRole.PageSize > 0 {
//    //     offset := (tabUserRole.PageNum - 1) * tabUserRole.PageSize
//    //     query = query.Offset(offset).Limit(tabUserRole.PageSize)
//    // }
//
//    err := query.Find(&tables).Error
//    return tables, err
//}
//
//// FindTabUserRolePageList 分页查询用户角色关联列表
//func (this *TabUserRole) FindTabUserRolePageList(DB *gorm.DB, satrtTime time.Time, endTime time.Time, pageNum int, pageSize int) ([]TabUserRole, int64, error) {
//    fmt.Printf("GetTabUserRolePageList：%#v \n", this)
//
//    var (
//        tabUserRoles []TabUserRole
//        total     int64
//    )
//
//    query := DB.Model(&TabUserRole{})
//
//// 构造查询条件
//        if this.Id != 0 { query = query.Where("id = ?", this.Id ) }
//        if this.UserId != 0 { query = query.Where("user_id = ?", this.UserId ) }
//        if this.RoleId != 0 { query = query.Where("role_id = ?", this.RoleId ) }
//        if !this.CreatedAt.IsZero() {
//            query = query.Where("created_at = ?", this.CreatedAt)
//            // query = query.Where("DATE(created_at) = ?", this.$column.goField.Format("2006-01-02"))
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
//        Find(&tabUserRoles).Error
//
//    if err != nil {
//        return nil, 0, err
//    }
//
//    return tabUserRoles, total, nil
//}
//
