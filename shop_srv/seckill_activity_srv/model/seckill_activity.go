package model

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

// LeavingMessages 留言信息
type LeavingMessages struct {
	BaseModel
	UserId      int32  `gorm:"type:int(11)"`
	MessageType int32  `gorm:"type:tinyint(1);comment:留言类型:1(留言),2(投诉),3(询问),4(售后),5(求购)"`
	Subject     string `gorm:"type:varchar(100);comment:留言标题"`
	Message     string `gorm:"type:varchar(100);comment:留言内容"`
	File        string `gorm:"type:varchar(200);comment:留言图片"`
}

// 自定义表名
func (LeavingMessages) TableName() string {
	return "leavingMessages"
}

type Address struct {
	BaseModel
	UserId       int32  `gorm:"type:int(11)"`
	Province     string `gorm:"type:varchar(20);comment:省"`
	City         string `gorm:"type:varchar(20);comment:市"`
	District     string `gorm:"type:varchar(20);comment:市/区"`
	Address      string `gorm:"type:varchar(20);comment:详细地址"`
	SignerName   string `gorm:"type:varchar(20);comment:收件人的姓名"`
	SignerMobile string `gorm:"type:varchar(15);comment:收件人手机"`
}

// UserFav 用户点赞表
type UserFav struct {
	BaseModel
	UserID  int32 `gorm:"type:int;index:idx_user_goods,unique"`
	GoodsID int32 `gorm:"type:int;index:idx_user_goods,unique"`
}

func (UserFav) TableName() string {
	return "userfav"
}
