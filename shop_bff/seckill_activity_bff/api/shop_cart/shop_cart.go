package shop_cart

import (
	"errors"
	goodsPb "github.com/china-li-shuo/shop_proto/goods"
	"github.com/gin-gonic/gin"
	inventoryPb "github.com/wujuan-glit/shop/inventory"
	orderPb "github.com/wujuan-glit/shop/order"
	"net/http"
	"order_bff/api"
	"order_bff/forms"
	"order_bff/global"
	"order_bff/model"
)

// CreateCart 添加购物车
func CreateCart(c *gin.Context) {

	//接收用户id
	user_id, _ := c.Get("user_id")

	//验证表单信息
	var cartForm forms.ShopCartItemForm
	if err := c.ShouldBind(&cartForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 999,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	//需要查询商品表是否有此商品信息
	_, err := global.GoodsClient.GetGoodsDetail(c, &goodsPb.GoodInfoRequest{Id: cartForm.GoodsId})
	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
		return
	}

	////查询库存信息  返回库存和商品id
	//invDetail, err := global.InventoryClient.InvDetail(c, &inventoryPb.GoodsInvInfo{
	//	GoodsId: cartForm.GoodsId,
	//	Num:     cartForm.Num,
	//})
	//if err != nil {
	//	api.HandleGrpcErrorToHttp(c, err)
	//	return
	//}
	//if invDetail.Num < cartForm.Num {
	//	err = errors.New("库存不足")
	//	api.ReturnErrorJson(c, err)
	//	return
	//
	//}
	//进行添加购物车
	item, err := global.OrderClient.CreateCartItem(c, &orderPb.CartItemReq{
		UserId:  user_id.(int32),
		GoodsId: cartForm.GoodsId,
		Nums:    cartForm.Num,
		Checked: true,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "购物车添加成功",
		"data":    item,
	})
}

// CartList 购物车列表
func CartList(c *gin.Context) {

	user_id, _ := c.Get("user")

	//获取到购物车列表  但是里面没有商品详情
	list, err := global.OrderClient.CartItemList(c, &orderPb.UserInfo{Id: user_id.(int32)})
	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
		return
	}
	//返回结果
	//{
	//	"info": [
	//		{
	//			"id": 1,
	//			"userId": 1,
	//			"goodsId": 421,
	//			"nums": 15,
	//			"checked": false
	//		},
	//		{
	//			"id": 5,
	//			"userId": 1,
	//			"goodsId": 435,
	//			"nums": 6,
	//			"checked": true
	//		}
	//	],
	//	"total": 2
	//}
	var ids []int32
	var goodsList []*model.CartGoodsInfo
	for _, info := range list.Info {

		ids = append(ids, info.GoodsId)

		goodsList = append(goodsList, &model.CartGoodsInfo{
			Id:      info.Id,
			UserId:  info.UserId,
			Nums:    info.Nums,
			Checked: info.Checked,
		})

	}

	if len(ids) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "购物车为空",
			"total":   0,
		})
		return
	}
	goods, err := global.GoodsClient.BatchGetGoods(c, &goodsPb.BatchGoodsIdInfo{Id: ids})
	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
		return
	}

	//因为购物车列表里面是只有商品id  所以需要获取商品详情
	for i, info := range goods.Data {
		goodsList[i].GoodsId = info.Id
		goodsList[i].GoodName = info.Name
		goodsList[i].GoodImage = info.GoodsFrontImage
		goodsList[i].GoodPrice = info.ShopPrice
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "购物车列表获取成功",
		"data":    goodsList,
		"total":   list.Total,
	})
}

// UpdateCart 更新购物车
func UpdateCart(c *gin.Context) {
	//id := c.PostForm("id")

	var shopCartItemUpdateForm forms.ShopCartItemUpdateForm

	if err := c.ShouldBind(&shopCartItemUpdateForm); err != nil {
		api.ReturnErrorJson(c, err)
		return
	}

	var checkedBool bool = true
	if shopCartItemUpdateForm.Checked == false {
		checkedBool = false
	}
	user_id, _ := c.Get("user")

	//查询库存信息  返回库存和商品id
	invDetail, err := global.InventoryClient.InvDetail(c, &inventoryPb.GoodsInvInfo{
		GoodsId: shopCartItemUpdateForm.GoodsId,
		Num:     shopCartItemUpdateForm.Num,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
		return
	}

	if invDetail.Num < int32(shopCartItemUpdateForm.Num) {
		err = errors.New("库存不足")
		api.ReturnErrorJson(c, err)
		return

	}

	_, err = global.OrderClient.UpdateCartItem(c, &orderPb.CartItemReq{
		Id:      shopCartItemUpdateForm.CartId,
		UserId:  user_id.(int32),
		Nums:    shopCartItemUpdateForm.Num,
		Checked: checkedBool,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "购物车更新成功",
	})

}

// DeleteCartList 删除购物车
func DeleteCart(c *gin.Context) {

}
