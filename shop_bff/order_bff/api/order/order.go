package order

import (
	"github.com/gin-gonic/gin"
	orderPb "github.com/wujuan-glit/shop/order"
	"net/http"
	"order_bff/api"
	"order_bff/forms"
	"order_bff/global"
)

// CreateOrder 创建订单
func CreateOrder(c *gin.Context) {

	user_id, _ := c.Get("user_id")
	var createOrderForm forms.CreateOrderItem

	if err := c.ShouldBind(&createOrderForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 999,
			"msg":  "err:" + err.Error(),
			"data": createOrderForm,
		})
		return
	}

	order, err := global.OrderClient.CreateOrder(c, &orderPb.OrderReq{
		UserId:  user_id.(int32),
		Address: createOrderForm.Address,
		Name:    createOrderForm.Name,
		Mobile:  createOrderForm.Mobile,
		Post:    createOrderForm.Post,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "订单创建成功",
		"data":    order,
	})
}

// ListOrder 订单列表
func ListOrder(c *gin.Context) {
	user_id, _ := c.Get("user_id")

	var orderList forms.ListOrderItem

	if err := c.ShouldBind(&orderList); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 999,
			"msg":  "err:" + err.Error(),
		})

		return
	}

	list, err := global.OrderClient.OrderList(c, &orderPb.OrderFilterReq{
		UserId:      user_id.(int32),
		Pages:       orderList.Pages,
		PagePerNums: orderList.PagePerNums,
		PayType:     orderPb.OrderFilterReq_PayTypes(orderList.PayType),
		Status:      orderList.Status,
	})

	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "订单列表获取成功",
		"data":    list,
	})
}

// 更新订单
func UpdateOrder(c *gin.Context) {
	//user_id, _ := c.Get("user_id")

	order, err := global.OrderClient.UpdateOrder(c, &orderPb.UpdateOrderInfo{
		Id:      1,
		PayType: 2,
		Status:  2,
		TradeNo: "20240818204118171",
		OrderSn: "",
	})

	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "订单修改成功",
		"data":    order,
	})
}
