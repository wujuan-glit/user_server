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

// ShoppingCart 商品购物车
type ShoppingCart struct {
	BaseModel
	UserID  int32 `gorm:"type:int;index"`
	GoodsID int32 `gorm:"type:int;index"`
	Nums    int32 `gorm:"type:int"`
	Checked bool
}

// OrderInfo 订单信息表
type OrderInfo struct {
	BaseModel
	UserID       int32   `gorm:"type:int;index"`
	OrderSn      string  `gorm:"type:varchar(100);index;comment:订单号"` //自己公司生成的，必须唯一
	PayType      int8    `gorm:"type:tinyint(1);comment:1(支付宝),2(微信)"`
	Status       int8    `gorm:"type:tinyint(1);comment:1(待支付),2(成功),3(超时关闭),4(交易创建),5(交易结束)"`
	TradeNo      string  `gorm:"type:varchar(100);comment:交易号"`   //支付平台的唯一号，支付流水号 微信流水号 后期可做对账系统 退款
	OrderMount   float32 `gorm:"type:decimal(10,2);comment:订单金额"` //订单金额 float精度会丢失 decimal 精度不丢失 mysql 5.7 版本开始支持 int(金额*100存角)
	PayTime      string  `gorm:"comment:支付时间"`
	Address      string  `gorm:"type:varchar(100);comment:收货地址"` //地址快照
	SignerName   string  `gorm:"type:varchar(20);comment:收货人名称"`
	SingerMobile string  `gorm:"type:varchar(11);comment:收货人电话"`
	Post         string  `gorm:"type:varchar(20);comment:订单备注"`
}

// OrderGoods 订单商品
type OrderGoods struct {
	BaseModel
	OrderId    int32   `gorm:"type:int;index;comment:订单id"`
	GoodsId    int32   `gorm:"type:int;index;comment:商品id"`
	GoodsName  string  `gorm:"type:varchar(100);index;comment:商品名称"`
	GoodsImage string  `gorm:"type:varchar(200)"`
	GoodsPrice float32 `gorm:"type:decimal(10,2);comment:商品价格"`
	Nums       int32   `gorm:"type:int(10);comment:商品数量"`
}

// GoodsStock 商品库存记录表
type GoodsStock struct {
	BaseModel
	GoodsId  int32  `gorm:"type:int(11);comment:商品id"`
	GoodsNum int32  `gorm:"type:int(11);comment:商品数量"`
	OrderSn  string `gorm:"type:varchar(100);comment:订单编号"`
}
