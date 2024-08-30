package initiliza

import (
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	goodsPb "github.com/china-li-shuo/shop_proto/goods"
	inventoryPb "github.com/wujuan-glit/shop/inventory"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"order_srv/handler"

	"github.com/apache/rocketmq-client-go/v2"
	_ "github.com/mbobakov/grpc-consul-resolver"
	orderPb "github.com/wujuan-glit/shop/order"
	"order_srv/global"
	"os"
	"os/signal"
	"syscall"
)

// InitConn 初始化grpc连接
func InitConn() {

	g := grpc.NewServer()

	s := handler.OrderService{}

	orderPb.RegisterOrderServer(g, &s)

	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", global.UserServerConfig.Ip, global.UserServerConfig.Port))

	if err != nil {
		zap.S().Info("网络监听失败", err.Error())
		return
	}

	grpc_health_v1.RegisterHealthServer(g, health.NewServer())

	// 启动 gRPC 服务器并开始监听。
	go func() {
		err = g.Serve(listen)
		if err != nil {
			zap.S().Info("订单服务监听失败" + err.Error())
			return
		}
	}()

	DelayOrderConsumer()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	//接收终止信号

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	//注销服务
	Deregister()

	zap.S().Info("服务注销成功")

}

// InitGoodsServerConn  连接商品
func InitGoodsServerConn() {
	goodsConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", global.UserServerConfig.Consul.Host, global.UserServerConfig.Consul.Port, global.UserServerConfig.Consul.GoodsService),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))

	if err != nil {
		zap.S().Panic("err", err)
		return
	}
	//连接库存微服务
	global.GoodsClient = goodsPb.NewGoodsClient(goodsConn)
	inventoryConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", global.UserServerConfig.Consul.Host, global.UserServerConfig.Consul.Port, global.UserServerConfig.Consul.InventoryService),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))

	if err != nil {
		zap.S().Panic("err", err)
		return
	}

	global.InventoryClient = inventoryPb.NewInventoryClient(inventoryConn)

}

// DelayOrderConsumer 消费者
func DelayOrderConsumer() {
	// 创建消费者
	c, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName("testGroup1"),
		//consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{fmt.Sprintf("%s:%d", global.UserServerConfig.Rocketmq.Host, global.UserServerConfig.Rocketmq.Port)})),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{"10.3.189.2:9876"})),
		consumer.WithRetry(2),
	)

	//如果订阅的主题不存在，进行每2秒重试一次，直至订阅成功为止

	if err != nil {
		zap.S().Info("初始化消费者失败", err)
		return
	}

	// 订阅主题`delay_order`，并设置回调函数
	err = c.Subscribe("delay_order", consumer.MessageSelector{}, handler.OrderTimeout)

	if err != nil {
		zap.S().Info("delay_order_consumer error", err.Error())
		return
	}
	err = c.Start()
	if err != nil {
		zap.S().Info("delay_order_consumer error", err.Error())
		return
	}
}
