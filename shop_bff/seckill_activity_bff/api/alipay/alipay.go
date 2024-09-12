package alipay

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	orderPb "github.com/wujuan-glit/shop/order"
	"go.uber.org/zap"
	"log"
	"math/rand"
	"net/http"
	"order_bff/api"
	"order_bff/global"
	"order_bff/model"
	"os"
	"strconv"
	"time"
)

// GenerateAlipayUrl 生成支付宝链接   买家账号:vhgkfv0395@sandbox.com
//
//	func GenerateAlipayUrl(c *gin.Context) {
//		order_sn := c.PostForm("order_sn")
//
//		if len(order_sn) == 0 {
//			err := errors.New("参数不能空")
//			api.ReturnErrorJson(c, err)
//			return
//		}
//
//		client, err := alipay.New(global.ServerConfig.Alipay.Appid, global.ServerConfig.Alipay.PrivateKey, false)
//		if err != nil {
//			c.JSON(http.StatusOK, gin.H{
//				"code":    0,
//				"message": "支付宝实例化失败",
//				"data":    nil,
//			})
//			return
//		}
//		//查询订单信息
//		detail, err := global.OrderClient.OrderDetail(c, &orderPb.OrderReq{
//			OrderSn: order_sn,
//		})
//
//		if err != nil {
//			api.ReturnErrorJson(c, err)
//			return
//		}
//		//strconv.FormatFloat(float64(detail.OrderInfo.Total), 'f', 1, 64)  第三个参数是小数点后几位
//		total := strconv.FormatFloat(float64(detail.OrderInfo.Total), 'f', 2, 64)
//		var p = alipay.TradeWapPay{}
//		p.NotifyURL = "http://29aaa0b6.r20.cpolar.top/o/v1/alipay/notify" //异步回调地址  用户处理更新订单信息
//
//		p.ReturnURL = "http://127.0.0.1:8888/o/v1/alipay/return" //同步回调地址 返回给用户订单详情
//
//		p.Subject = "生鲜支付" + detail.OrderInfo.OrderSn
//
//		p.OutTradeNo = detail.OrderInfo.OrderSn
//
//		p.TotalAmount = total
//
//		p.ProductCode = "QUICK_WAP_WAY"
//
//		url, err := client.TradeWapPay(p)
//
//		if err != nil {
//			api.ReturnErrorJson(c, err)
//			return
//		}
//
//		// 这个 payURL 即是用于打开支付宝支付页面的 URL，可将输出的内容复制，到浏览器中访问该 URL 即可打开支付页面。
//		var payURL = url.String()
//
//		c.JSON(http.StatusOK, gin.H{
//			"code":    0,
//			"message": "支付链接",
//			"data": map[string]interface{}{
//				"pay_url": payURL,
//			},
//		})
//
// }
//
// // Return 同步回调地址
// func ReturnUrl(c *gin.Context) {
//
//		orderSn := c.Query("out_trade_no") //自己系统的订单号
//
//		if len(orderSn) == 0 {
//			c.JSON(http.StatusOK, gin.H{
//				"code": 0,
//				"msg":  "订单号不能为空",
//				"data": nil,
//			})
//			return
//		}
//		detail, err := global.OrderClient.OrderDetail(c, &orderPb.OrderReq{
//			OrderSn: orderSn,
//		})
//		if err != nil {
//			return
//		}
//
//		c.JSON(http.StatusOK, gin.H{
//			"code":    0,
//			"message": "订单详情",
//			"data":    detail,
//		})
//	}
//
// // Notice 异步回调地址  主要做订单状态的修改
// func NotifyUrl(c *gin.Context) {
//
//		tradeStatus := c.PostForm("trade_status")
//
//		tradeNo := c.PostForm("trade_no")
//
//		outTradeNo := c.PostForm("out_trade_no")
//		log.Println("c", c)
//		if tradeStatus == "TRADE_SUCCESS" {
//			//交易成功  Status:  2  表示支付成功  PayType: 1  表示支付宝支付
//			_, err := global.OrderClient.UpdateOrder(c, &orderPb.UpdateOrderInfo{
//				PayType: 1,
//				Status:  2,
//				TradeNo: tradeNo,
//				OrderSn: outTradeNo,
//			})
//			if err != nil {
//				log.Fatalf("err", err)
//			}
//		}
//		c.String(200, "success")
//	}
//
// 生成支付链接
// InitAlipay 初始化支付宝客户端
var err error

func InitAlipay() {
	global.AliClient, err = alipay.New(global.ServerConfig.Alipay.Appid, global.ServerConfig.Alipay.PrivateKey, false)
	if err != nil {
		zap.S().Info("实例化支付宝失败", err)
		return
	}
	err = global.AliClient.LoadAliPayPublicKey(global.ServerConfig.Alipay.PublicKey)
	if err != nil {
		zap.S().Info("载入支付宝公钥失败", err)
		return
	}
}
func GenerateAlipayUrl(c *gin.Context) {

	order_sn := c.PostForm("order_sn")

	if len(order_sn) == 0 {
		api.ReturnErrorJson(c, err)
		return
	}

	//根据订单编号进行查询订单详情
	detail, err := global.OrderClient.OrderDetail(c, &orderPb.OrderReq{
		OrderSn: order_sn,
	})

	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
		return
	}
	var p = alipay.TradeWapPay{}
	total := strconv.FormatFloat(float64(detail.OrderInfo.Total), 'f', 2, 64)
	p.NotifyURL = "http://6a6d884f.r20.cpolar.top/o/v1/alipay/notify" //异步回调地址  是我们用来更改订单信息的  需要内网穿透
	p.ReturnURL = "http://192.168.1.114:8891/o/v1/alipay/return"      //同步回调地址  用来展示给用户的界面 一般是订单详情
	p.Subject = "生鲜" + detail.OrderInfo.OrderSn
	p.OutTradeNo = detail.OrderInfo.OrderSn
	p.TotalAmount = total
	p.ProductCode = "QUICK_WAP_WAY"

	url, err := global.AliClient.TradeWapPay(p)
	if err != nil {
		fmt.Println(err)
	}

	// 这个 payURL 即是用于打开支付宝支付页面的 URL，可将输出的内容复制，到浏览器中访问该 URL 即可打开支付页面。
	var payURL = url.String()
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "生成支付链接",
		"data": map[string]string{
			"pay_url": payURL,
		},
	})
}

// 同步回调
func ReturnUrl(c *gin.Context) {
	//通过订单编号查询订单详情
	out_trade_no := c.Query("out_trade_no")
	if len(out_trade_no) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "订单号不能为空",
		})

		return
	}

	detail, err := global.OrderClient.OrderDetail(c, &orderPb.OrderReq{
		OrderSn: out_trade_no,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "订单详情",
		"data":    detail,
	})
}

//NotifyUrl 使用redis队列实现异步回调

//func NotifyUrl(c *gin.Context) {
//	trade_no := c.PostForm("trade_no")
//	trade_status := c.PostForm("trade_status")
//	out_trade_no := c.PostForm("out_trade_no")
//
//	//除了交易成功 还有几种状态  如交易关闭  等待交易  交易完成
//
//	// 定义订单信息变量
//	var orderInfo model.OrderInfo
//
//	// 根据交易状态设置订单信息
//	//1(待支付),2(成功),3(超时关闭),4(交易失败),5(交易结束)
//	switch trade_status {
//	case "TRADE_SUCCESS":
//		//交易成功
//		orderInfo.Status = 2
//
//	case "TRADE_PAY_FAIL":
//		//交易失败
//		orderInfo.Status = 4
//
//	case "TRADE_CLOSED":
//		//交易关闭
//		orderInfo.Status = 3
//
//	case "WAIT_BUYER_PAY":
//		//等待买家付款
//		orderInfo.Status = 1
//
//	case "TRADE_FINISHED":
//		//交易结束 不可退款
//		orderInfo.Status = 5
//
//	default:
//		// 可以根据需要处理其他交易状态
//		zap.S().Info("未处理交易状态: " + trade_status)
//		c.String(http.StatusBadRequest, "未处理交易状态")
//		return
//	}
//	//1代表支付类型是支付宝
//	orderInfo = model.OrderInfo{
//		PayType: 1,
//		Status:  orderInfo.Status,
//		TradeNo: trade_no,
//		OrderSn: out_trade_no,
//	}
//	// 序列化订单信息
//	marshal, err := json.Marshal(orderInfo)
//	if err != nil {
//		zap.S().Info("序列化失败: " + err.Error())
//		c.String(http.StatusInternalServerError, "序列化失败")
//		return
//	}
//
//	// 将修改订单放到异步队列里面，使用唯一的队列键
//
//	push := global.RedisClient.LPush(context.Background(), "order_info", marshal)
//
//	if push.Err() != nil {
//
//		zap.S().Info("放入redis队列失败: " + push.Err().Error())
//
//		c.String(http.StatusInternalServerError, "放入redis队列失败")
//
//		return
//	}
//
//	c.String(http.StatusOK, "success")
//}

//NotifyUrl 使用rocketMQ实现异步回调

func NotifyUrl(c *gin.Context) {
	trade_no := c.PostForm("trade_no")
	trade_status := c.PostForm("trade_status")
	out_trade_no := c.PostForm("out_trade_no")

	//除了交易成功 还有几种状态  如交易关闭  等待交易  交易完成

	// 定义订单信息变量
	var orderInfo model.OrderInfo

	// 根据交易状态设置订单信息
	//1(待支付),2(成功),3(超时关闭),4(交易失败),5(交易结束)
	switch trade_status {
	case "TRADE_SUCCESS":
		//交易成功
		orderInfo.Status = 2

	case "TRADE_PAY_FAIL":
		//交易失败
		orderInfo.Status = 4

	case "TRADE_CLOSED":
		//交易关闭
		orderInfo.Status = 3

	case "WAIT_BUYER_PAY":
		//等待买家付款
		orderInfo.Status = 1

	case "TRADE_FINISHED":
		//交易结束 不可退款
		orderInfo.Status = 5

	default:
		// 可以根据需要处理其他交易状态
		zap.S().Info("未处理交易状态: " + trade_status)
		c.String(http.StatusBadRequest, "未处理交易状态")
		return
	}
	//1代表支付类型是支付宝
	order := model.OrderInfo{
		PayType: 1,
		Status:  orderInfo.Status,
		TradeNo: trade_no,
		OrderSn: out_trade_no,
	}
	// 序列化订单信息
	msg, err := json.Marshal(order)
	if err != nil {
		zap.S().Info("序列化失败: " + err.Error())
		c.String(http.StatusInternalServerError, "序列化失败")
		return
	}

	// 将修改订单放到异步队列里面，使用唯一的队列键

	RocketMqProducer(msg)

	c.String(http.StatusOK, "success")
}

// AlipayConsumer 支付
// AlipayConsumer 使用redis队列实现
func AlipayConsumer() {
	go func() {
		for {
			//使用阻塞弹出操作，等待直到有数据可取或超时 BLPOP  如果队列里面没有值  和
			val := global.RedisClient.LPop(context.Background(), "order_info").Val()

			if len(val) == 0 {

				zap.S().Info("队列内容格式错误")

				continue
			}

			var orderInfo model.OrderInfo

			err := json.Unmarshal([]byte(val), &orderInfo)

			if err != nil {

				zap.S().Info("反序列化失败: " + err.Error())

				continue
			}
			now := time.Now().Format(time.DateTime)
			// 更新订单状态
			_, err = global.OrderClient.UpdateOrder(context.Background(), &orderPb.UpdateOrderInfo{
				PayType: orderInfo.PayType,
				Status:  orderInfo.Status,
				TradeNo: orderInfo.TradeNo,
				OrderSn: orderInfo.OrderSn,
				PayTime: now,
			})

			if err != nil {

				zap.S().Info("订单更新失败: " + err.Error())

				continue
			}

			log.Println("修改订单状态成功")
		}
	}()

}

// RocketMqProducer  使用rocketmq实现
func RocketMqProducer(msgs []byte) {

	intn := rand.Intn(100)
	name := GenerateOrderSn(int32(intn))

	p, err := rocketmq.NewProducer(
		producer.WithGroupName(name),
		producer.WithNameServer([]string{fmt.Sprintf("%s:%d", global.ServerConfig.Rocketmq.Host, global.ServerConfig.Rocketmq.Port)}),
	)
	if err != nil {
		zap.S().Info("生产者失败", err)
		return
	}

	err = p.Start()
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}

	msg := &primitive.Message{
		Topic: "alipay",
		Body:  msgs,
	}
	_, err = p.SendSync(context.Background(), msg)

	if err != nil {
		zap.S().Info("支付发送消息失败", err)
		return
	}

}

// UpdateOrder 更改订单状态

func UpdateOrder(ctx context.Context, messages ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

	//把messages的信息反解到订单结构体

	for _, msg := range messages {
		var orderInfo model.OrderInfo
		//将延迟队列里的消息反解到orderInfo结构体
		err := json.Unmarshal(msg.Body, &orderInfo)
		if err != nil {
			zap.S().Info("反序列化失败")
			return 0, err
		}

		now := time.Now().Format(time.DateTime)
		// 更新订单状态
		_, err = global.OrderClient.UpdateOrder(context.Background(), &orderPb.UpdateOrderInfo{
			PayType: orderInfo.PayType,
			Status:  orderInfo.Status,
			TradeNo: orderInfo.TradeNo,
			OrderSn: orderInfo.OrderSn,
			PayTime: now,
		})

		if err != nil {
			zap.S().Info("订单更新失败: " + err.Error())
			return consumer.ConsumeRetryLater, nil
		}

		log.Println("修改订单状态成功")
	}

	return consumer.ConsumeSuccess, nil
}

// RocketMqConsumer 使用rocketMq消费者
func RocketMqConsumer() {
	// 创建消费者

	c, err := rocketmq.NewPushConsumer(

		consumer.WithGroupName("testGroup1"),

		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{"192.168.1.114:9876"})),

		consumer.WithRetry(2),
	)

	//如果订阅的主题不存在，进行每2秒重试一次，直至订阅成功为止

	if err != nil {

		zap.S().Info("初始化消费者失败", err)

		return
	}

	// 订阅主题`delay_order`，并设置回调函数
	err = c.Subscribe("alipay", consumer.MessageSelector{}, UpdateOrder)

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

// GenerateOrderSn 生成订单编号
func GenerateOrderSn(userId int32) string {
	now := time.Now().Format("20060102150405")

	intn := rand.Intn(100)

	orderSn := now + strconv.Itoa(int(userId)) + strconv.Itoa(intn)
	return orderSn

}

// Refund 支付宝退款
// 退款之前要查询订单信息  如果是交易状态是交易完成  就不能进行退款
func Refund(c *gin.Context) {

	// 从请求中获取退款参数，例如：outTradeNo（商户订单号）, refundAmount（退款金额）等
	outTradeNo := c.PostForm("trade_no")
	refundAmount := c.PostForm("refund_amount")
	refundReason := c.PostForm("refund_reason")
	tradeNo := c.PostForm("trade_no")

	if len(outTradeNo) == 0 || len(refundAmount) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "参数不能为空",
		})
		return
	}

	//查询订单状态

	detail, err := global.OrderClient.OrderDetail(c, &orderPb.OrderReq{
		OrderSn: outTradeNo,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
		return
	}
	//1(待支付),2(成功),3(超时关闭),4(交易创建),5(交易结束)
	if detail.OrderInfo.Status == 5 {
		//代表交易结束 不能够退款
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "交易已完成，无法退款",
		})
		return
	}
	//查询退款记录表
	//1 未退款 2 已退款 3 无法退款

	// 创建退款请求
	req := alipay.TradeRefund{
		OutTradeNo:   outTradeNo,
		TradeNo:      tradeNo,
		RefundAmount: refundAmount,
		RefundReason: refundReason,
	}

	// 发起退款请求
	resp, err := global.AliClient.TradeRefund(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "退款失败", "message": err.Error()})
		return
	}

	// 检查退款结果
	if resp.IsSuccess() {
		c.JSON(http.StatusOK, gin.H{"success": "退款成功"})
	} else {
		c.JSON(http.StatusOK, gin.H{"error": "退款失败", "message": resp.Msg + ": " + resp.SubMsg})
	}

}
