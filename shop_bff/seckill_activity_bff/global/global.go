package global

import (
	"fmt"
	goodsPb "github.com/china-li-shuo/shop_proto/goods"
	ut "github.com/go-playground/universal-translator"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/redis/go-redis/v9"
	uuid "github.com/satori/go.uuid"
	"github.com/smartwalle/alipay/v3"
	inventoryPb "github.com/wujuan-glit/shop/inventory"
	orderPb "github.com/wujuan-glit/shop/order"

	"order_bff/config"
)

var (
	GoodsClient  goodsPb.GoodsClient
	NacosClient  config_client.IConfigClient
	ServerConfig config.GoodsServer
	// 定义一个全局翻译器T
	Trans           ut.Translator
	ServerId        string = fmt.Sprintf("%s", uuid.NewV4())
	InventoryClient inventoryPb.InventoryClient
	OrderClient     orderPb.OrderClient
	RedisClient     *redis.Client
	AliClient       *alipay.Client
)
