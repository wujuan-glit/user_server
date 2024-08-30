package global

import (
	goodsPb "github.com/china-li-shuo/shop_proto/goods"
	inventoryPb "github.com/wujuan-glit/shop/inventory"
	"gorm.io/gorm"
	"order_srv/config"
)

var (
	ServerID         string
	UserServerConfig config.UserServer
	Db               *gorm.DB
	GoodsClient      goodsPb.GoodsClient
	InventoryClient  inventoryPb.InventoryClient
)
