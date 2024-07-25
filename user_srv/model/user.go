package model

import (
	"database/sql"
	"errors"
	"gorm.io/gorm"
	"time"
	"user/user_srv/global"
)

type BaseModel struct {
	ID        int32 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt sql.NullTime   //允许为空
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type User struct {
	BaseModel
	Password string `gorm:"type:varchar(150);not null;comment:密码"`
	NickName string `gorm:"type:varchar(50);not null;comment:昵称"`
	Mobile   string `gorm:"type:varchar(15);not null;comment:手机号"`
	Role     int32  `gorm:"column:role;type:tinyint(1);not null;comment:角色"`
	Birthday uint64 `gorm:"column:birthday;type:int(10);not null;comment:生日"`
	Gender   string `gorm:"type:varchar(2);not null;comment:性别"`
}

// 根据用户手机号获取用用户信息
func (c *User) GetUserInfoByMobile() *gorm.DB {
	tx := global.Db.Table("users").Where("mobile = ?", c.Mobile).Limit(1).Find(&c)

	return tx
}

func (c *User) AddUser() error {
	tx := global.Db.Table("users").Create(&c)

	if tx.Error != nil || tx.RowsAffected < 0 {
		return errors.New("用户注册失败")
	}
	return nil
}
