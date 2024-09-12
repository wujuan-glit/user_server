package mdoel

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        int32 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt sql.NullTime   //允许为空
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Permission 权限表
type Permission struct {
	Name        string `gorm:"type:varchar(50);comment:权限名称"`
	Description string `gorm:"type:varchar(100);comment:权限的描述"`
}

// Role 角色表
type Role struct {
	BaseModel
	Name        string `gorm:"type:varchar(50);comment:角色名称"`
	Description string `gorm:"type:varchar(100);comment:角色的描述"`
}

// UserRoleMapping 用户角色关联表
type UserRoleMapping struct {
	BaseModel
	UserId int32 `gorm:"type:int(11);comment:用户id"`
	RoleId int32 `gorm:"type:int(11);comment:角色id"`
}

// RolePermissionMapping 角色权限关联表
type RolePermissionMapping struct {
	BaseModel
	UserId int32 `gorm:"type:int(11);comment:用户id"`
	RoleId int32 `gorm:"type:int(11);comment:权限id"`
}
