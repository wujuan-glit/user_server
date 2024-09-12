package initiliza

import (
	"fmt"
	goodsPb "github.com/china-li-shuo/shop_proto/goods"
	_ "github.com/mbobakov/grpc-consul-resolver"
	inventoryPb "github.com/wujuan-glit/shop/inventory"
	orderPb "github.com/wujuan-glit/shop/order"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"order_bff/global"
)

func InitServerConn() {

	//拨号连接商品微服务
	goodsConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", global.ServerConfig.Consul.Host, global.ServerConfig.Consul.Port, global.ServerConfig.Consul.GoodsService),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))

	if err != nil {
		zap.S().Panic("err", err)
		return
	}

	global.GoodsClient = goodsPb.NewGoodsClient(goodsConn)
	//拨号连接库存微服务
	inventoryConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", global.ServerConfig.Consul.Host, global.ServerConfig.Consul.Port, global.ServerConfig.Consul.InventoryService),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))

	if err != nil {
		zap.S().Panic("err", err)
		return
	}

	global.InventoryClient = inventoryPb.NewInventoryClient(inventoryConn)

	//拨号连接订单微服务
	orderConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", global.ServerConfig.Consul.Host, global.ServerConfig.Consul.Port, global.ServerConfig.Consul.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))

	if err != nil {
		zap.S().Panic("err", err)
		return
	}

	global.OrderClient = orderPb.NewOrderClient(orderConn)

}
