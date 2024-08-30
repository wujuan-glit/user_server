package alipay

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	orderPb "github.com/wujuan-glit/shop/order"
	"go.uber.org/zap"
	"net/http"
	"order_bff/api"
	"order_bff/global"
	"strconv"
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

func GenerateAlipayUrl(c *gin.Context) {
	//privateKey := "MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCNf0g6A6R9Kurd0Wbk0w8mpJUjRO14MJRmVXLIS0fBzTrup+tkpH627mANpE4y2D/00Lv/ZI8o0oRET0tZnaz4k/IEX/UIst9A92ju6ofsnW9qQG7CzeoARptVcaC74/1hMRiIzViD0RLu8jI3yvsa4PhKwbM340dRKj0DCq1duQlJY4ucF2tDVaLEyolm140dhtzCL10KFxWAcLQiKAXDT3mWRoi2DlPttmyBOJVxTa0nWoYQVr0tN0k480EKGAw3wYyzwgPhgqAsLsWsaMTwae+U1kxknq+zB+C5f3qjRNB9/5eXBh0gnjD9Gb+5c8uO7b48TvlyjsllN35y/ROHAgMBAAECggEAU2vNO1bWbW0WFzzTuuisMA4sVyTWFFwfwc1y5J9taNcEfZvGbgmFI3iabLCH4fYYjs9ZZxL0TA8BJ/zP4b/SMKOYtfeU0VITyYuT8/eVt2yCOVRPeM5JvWvjPJbHOr8JrXlyi4T1QJHM5c8oyDgFny0vdXOJo9N9Ql7ypY5v86a3Wo6Ue5KjK8e3Pmc9wPw54utiI7phiGJHBX2F0HwVhDWcaZFr3tgQgIcjorQjQmrktUKJ7moFhL1XCGTyigMxpYxRZAoR8qgCX9OSab3551/tKnviENPkAVPw8/VJxDRIhtLIsy3u3YhipEyUV4YsC77gKGH35dnlBGtZXeWSoQKBgQDZzXS4OAEfT1lafLaq3SUOSWeO00T/zOnV5LznmQiuIM2J8xA2YXN62PrXIPRq/c/tfoTtPLP5MUbJqqqNhpORm9CmNT4OIpVJVryN1TWXMooqFZvgDVo9Ht3OgpP5dsFd/jGblQO6bz6efvQi5ZnygxAjt5rvegVX+o/xtByJLwKBgQCmT/4LY7jN9TZDSPm5P1uuVv76TCyLretT4lPgy/LsnbrK+haTUAicrZZ3PmyYchlMDhQbtz0hGGW4TJkQR3A09LYIROceaLBFVHxd8T/L0dF7q/vcUnUdHRpGG158ssY6fgiVmlCARek5fj6t/LbasKbItCeuKu6+I5X66GnVKQKBgCV6OgRc9qx5jemJHjGGfhLYRK6J4gyWKQJ6Kps7dQfpcxSys254FFPmNDuCWyxx4i5+n8bmtB1EAmc/K7vQlWHvytZewP/TqZaGC0nojyEmPCoDr9+8zHNJ9WbMh0Pc0GcpD0YzPQH+lGrXc5DxqyzUqplKxalBeNvrrIstr99XAoGACv5GoKIa2SJYT+JG/4O8n62IdSsL1r/MSmMvgDB7AkD60+fsDhjAOPsQcxlhPEJugaR8l8ho9gMS1jfZ9kWCmT2DutAzJsNsw2huQBduTB62ZiJcJ5gbvazqy6+Lc1qt17f1AU6N+6yjWfWKVx3ZSGNc4u9loBGeblsT0t4CAOECgYEAxu2hPCrcVo+h9nyz/EDLaQ4acRYv0+RT5xl9Zzu2vCrLQ4vVke+k3ILa4hGFnZmaCRdi9CQ/KvHcKJWDk01ZwxHLVHSN898bsBo0Ml3af0iSJex+Acq6qaRrA/4a64KQ7NuHAC/3KQW2SRqnnGRwCsDwNBwJsYYWccqWeK6j+J4="
	client, err := alipay.New(global.ServerConfig.Alipay.Appid, global.ServerConfig.Alipay.PrivateKey, false)

	if err != nil {
		zap.S().Info("实例化支付宝失败", err)
		return
	}

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
	p.NotifyURL = "http://29aaa0b6.r20.cpolar.top/o/v1/alipay/notify" //异步回调地址  是我们用来更改订单信息的  需要内网穿透
	p.ReturnURL = "http://127.0.0.1:8888/o/v1/alipay/return"          //同步回调地址  用来展示给用户的界面 一般是订单详情
	p.Subject = "生鲜" + detail.OrderInfo.OrderSn
	p.OutTradeNo = detail.OrderInfo.OrderSn
	p.TotalAmount = total
	p.ProductCode = "QUICK_WAP_WAY"

	url, err := client.TradeWapPay(p)
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

//NotifyUrl 异步回调

func NotifyUrl(c *gin.Context) {
	trade_no := c.PostForm("trade_no")
	trade_status := c.PostForm("trade_status")
	out_trade_no := c.PostForm("out_trade_no")
	//除了交易成功 还有几种状态  如交易关闭  等待交易  交易完成
	if trade_status == "TRADE_SUCCESS" {
		_, err := global.OrderClient.UpdateOrder(c, &orderPb.UpdateOrderInfo{
			PayType: 1,
			Status:  2,
			TradeNo: trade_no,
			OrderSn: out_trade_no,
		})
		if err != nil {
			zap.S().Info("订单更新失败" + err.Error())
			return
		}
	}

	c.String(http.StatusOK, "success")
}
