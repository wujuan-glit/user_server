package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	goodsPb "github.com/china-li-shuo/shop_proto/goods"
	inventoryPb "github.com/wujuan-glit/shop/inventory"
	orderPb "github.com/wujuan-glit/shop/order"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"math/rand"
	"order_srv/global"
	"order_srv/model"
	"strconv"
	"time"
)

type OrderListener struct {
	Code        codes.Code
	Detail      string
	ID          int32
	OrderAmount float32
	Ctx         context.Context
}

// ExecuteLocalTransaction 执行本地事务 做下订单，订单商品和删除购物车  这三个事件没有成功则rollback舍弃消息  否则提交消息
func (o *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	var orderInfo model.OrderInfo
	err := json.Unmarshal(msg.Body, &orderInfo)
	if err != nil {
		zap.S().Errorf("订单信息发序列化失败")
		return 0
	}

	/*
	   新建订单
	       1. 从购物车中获取到选中的商品
	       2. 商品的价格自己查询 - 访问商品服务 (跨微服务)
	       3. 库存的扣减 - 访问库存服务 (跨微服务)
	       4. 订单的基本信息表 - 订单的商品信息表
	       5. 从购物车中删除已购买的记录
	*/

	//定义一个切片保存购物车商品数据
	var shopCart []*model.ShoppingCart
	//定义一个切片存放购物车下的商品id
	var goodsIds []int32

	//定义商品数据字典  里面key  是商品id  value 是商品数量
	goodsNumsMap := make(map[int32]int32)

	//查询购物车中用户选中的商品
	tx := global.Db.Model(&model.ShoppingCart{}).Where("user_id = ? and checked = true", orderInfo.UserID).Find(&shopCart)

	if tx.Error != nil || tx.RowsAffected == 0 {
		o.Code = codes.InvalidArgument
		o.Detail = "没有选中的商品"
		return primitive.RollbackMessageState
	}
	for _, cart := range shopCart {
		goodsIds = append(goodsIds, cart.GoodsID)
		goodsNumsMap[cart.GoodsID] = cart.Nums

	}

	//去商品微服务批量查询商品信息
	goods, err := global.GoodsClient.BatchGetGoods(context.Background(), &goodsPb.BatchGoodsIdInfo{Id: goodsIds})
	if err != nil {
		o.Code = codes.Internal
		o.Detail = "批量查询商品失败"
		return primitive.RollbackMessageState
	}

	//定义一个切片 保存扣减商品和数量信息
	var inventoryList []*inventoryPb.GoodsInvInfo

	//订单的总金额 = 所有商品的金额加一起 (商品的价格（goods）*购物车商品的数量(shopCarts.nums))
	var sum float32
	//循环商品信息  进行总价计算
	//定义一个切片 保存订单购买的所有商品
	var orderList []*model.OrderGoods

	for _, info := range goods.Data {

		sum += (info.ShopPrice * float32(goodsNumsMap[info.Id]))

		orderList = append(orderList, &model.OrderGoods{
			GoodsId:    info.Id,
			GoodsName:  info.Name,
			GoodsPrice: info.ShopPrice,
			GoodsImage: info.GoodsFrontImage,
			Nums:       goodsNumsMap[info.Id],
		})

		inventoryList = append(inventoryList, &inventoryPb.GoodsInvInfo{
			GoodsId: info.Id,
			Num:     goodsNumsMap[info.Id],
		})
	}

	//进行库存扣减
	_, err = global.InventoryClient.Sell(context.Background(), &inventoryPb.SellInfo{
		GoodsInfo: inventoryList,
		OrderSn:   orderInfo.OrderSn,
	})
	if err != nil {
		o.Code = codes.Internal
		o.Detail = fmt.Sprintf("扣减库存失败:%v", err.Error())
		return primitive.CommitMessageState
	}
	//开启事务
	begin := global.Db.Begin()

	//创建订单基本信息表
	order := model.OrderInfo{
		UserID:       orderInfo.UserID,
		OrderSn:      orderInfo.OrderSn,
		Status:       1,
		OrderMount:   sum,
		Address:      orderInfo.Address,
		SignerName:   orderInfo.SignerName,
		SingerMobile: orderInfo.SingerMobile,
		Post:         orderInfo.Post,
	}
	create := begin.Model(&orderInfo).Create(&order)

	if create.RowsAffected == 0 || create.Error != nil {
		o.Code = codes.Internal
		o.Detail = "创建订单失败"
		begin.Rollback()
		return primitive.CommitMessageState
	}

	//创建订单商品表
	for _, orderGoods := range orderList {
		orderGoods.OrderId = order.ID
	}
	//批量添加 第二个参数 允许同时添加多少个  这样写可以省去循环
	batches := begin.CreateInBatches(orderList, 100)

	if batches.Error != nil || batches.RowsAffected == 0 {
		o.Code = codes.Internal
		o.Detail = "创建订单商品失败"
		begin.Rollback()
		return primitive.CommitMessageState
	}

	//进行购物车删除

	db := global.Db.Model(&model.ShoppingCart{}).Where("user_id =? and checked= true ", orderInfo.UserID).Delete(&model.ShoppingCart{})

	if db.Error != nil || db.RowsAffected == 0 {
		o.Code = codes.Internal
		o.Detail = "购物车删除失败"

		begin.Rollback()

		return primitive.CommitMessageState
	}

	//发送延时消息
	//创建延迟队列订单信息
	p, err := rocketmq.NewProducer(
		producer.WithGroupName(GenerateOrderSn(orderInfo.UserID)),
		producer.WithNameServer([]string{fmt.Sprintf("%s:%d", global.UserServerConfig.Rocketmq.Host, global.UserServerConfig.Rocketmq.Port)}))
	if err != nil {
		tx.Rollback()
		zap.S().Error("生成producer失败")
		return primitive.CommitMessageState
	}

	//不要在一个进程中使用多个producer， 但是不要随便调用shutdown因为会影响其他的producer
	if err = p.Start(); err != nil {
		tx.Rollback()
		zap.S().Error("启动producer失败")
		return primitive.CommitMessageState
	}

	delayMsg := primitive.NewMessage("delay_order", msg.Body)
	//设置延迟队列时间 level 3 代表延迟10秒 5 1分钟  具体可以点进去代码看详情
	delayMsg.WithDelayTimeLevel(5)

	_, err = p.SendSync(context.Background(), delayMsg)

	if err != nil {
		zap.S().Errorf("发送延时消息失败: %v\n", err)
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "发送延时消息失败"
		return primitive.CommitMessageState
	}

	//提交事务
	begin.Commit()
	o.Code = codes.OK

	return primitive.RollbackMessageState
}

// CheckLocalTransaction 回查
func (o *OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	//查询订单状态
	var orderInfo model.OrderInfo
	err := json.Unmarshal(msg.Body, &orderInfo)
	if err != nil {
		zap.S().Errorf("反序列化失败" + err.Error())
		return 0
	}
	//查询订单状态
	tx := global.Db.Model(&model.OrderInfo{}).Where(&model.OrderInfo{OrderSn: orderInfo.OrderSn}).Limit(1).Find(&orderInfo)

	if tx.Error != nil || tx.RowsAffected == 0 {
		return primitive.CommitMessageState
	}

	return primitive.RollbackMessageState
}

type OrderService struct {
	orderPb.UnimplementedOrderServer
}

// CreateOrder 创建订单 RocketMq

func (OrderService) CreateOrder(ctx context.Context, req *orderPb.OrderReq) (*orderPb.OrderInfoResp, error) {

	//创建事物的生产者  并且传入配置信息
	orderListener := OrderListener{Ctx: context.Background()}
	intn := rand.Intn(100)
	name := GenerateOrderSn(int32(intn))

	p, err := rocketmq.NewTransactionProducer(
		&orderListener, //消息的监听者，必须写上ExecuteLocalTransaction方法这个方法主要是写本地事务逻辑，成功与否来提交半事务消息 CheckLocalTransaction方法是事务回查机制 如果ExecuteLocalTransaction方法一直没有提交半事务消息或者回滚半事务消息，就会走到这个回查方法了
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"10.3.189.2:9876"})), //连接rocketmq
		producer.WithRetry(1),        //生产消息重试次数
		producer.WithGroupName(name), // Use a unique name here
	)
	if err != nil {
		zap.S().Info("生成producer失败:", err)
		return nil, err
	}

	if err = p.Start(); err != nil {
		zap.S().Info("启动producer失败:", err.Error())
		return nil, err
	}
	orderInfo := model.OrderInfo{
		UserID:       req.UserId,
		OrderSn:      GenerateOrderSn(req.UserId),
		Address:      req.Address,
		SignerName:   req.Name,
		SingerMobile: req.Mobile,
		Post:         req.Post,
	}
	//应该在消息中具体指明一个订单的具体的商品的扣减情况
	jsonString, _ := json.Marshal(orderInfo)

	_, err = p.SendMessageInTransaction(context.Background(),
		primitive.NewMessage("order", jsonString))
	if err != nil {
		zap.S().Info("发送失败:", err)
		return nil, status.Errorf(codes.Internal, "发送消息失败")
	}
	if orderListener.Code != codes.OK {
		return nil, status.Error(orderListener.Code, orderListener.Detail)
	}

	return &orderPb.OrderInfoResp{
		Id:      orderListener.ID,
		UserId:  req.UserId,
		OrderSn: orderInfo.OrderSn,
		Post:    req.Post,
		Address: req.Address,
		Name:    req.Name,
		Mobile:  req.Mobile,
	}, nil

}

//CreateOrder 创建订单  TCC版
//创建订单TCC版
/*func (OrderService) CreateOrder(ctx context.Context, req *orderPb.OrderReq) (*orderPb.OrderInfoResp, error) {

	  // 新建订单
	    //   1. 从购物车中获取到选中的商品
	     //  2. 商品的价格自己查询 - 访问商品服务 (跨微服务)
	      // 3. 库存的扣减 - 访问库存服务 (跨微服务)
	      // 4. 订单的基本信息表 - 订单的商品信息表
	      // 5. 从购物车中删除已购买的记录


	//定义一个切片保存购物车商品数据
	var shopCart []*model.ShoppingCart
	//定义一个切片存放购物车下的商品id
	var goodsIds []int32

	//定义商品数据字典  里面key  是商品id  value 是商品数量
	goodsNumsMap := make(map[int32]int32)

	//查询购物车中用户选中的商品
	tx := global.Db.Model(&model.ShoppingCart{}).Where("user_id = ? and checked = true", req.UserId).Find(&shopCart)

	if tx.Error != nil {
		return nil, status.Errorf(codes.Internal, "查询购物车商品失败")
	}
	if tx.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车为空")
	}
	for _, cart := range shopCart {
		goodsIds = append(goodsIds, cart.GoodsID)
		goodsNumsMap[cart.GoodsID] = cart.Nums
	}

	//去商品微服务批量查询商品信息
	goods, err := global.GoodsClient.BatchGetGoods(context.Background(), &goodsPb.BatchGoodsIdInfo{Id: goodsIds})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "没有此商品信息")
	}

	//定义一个切片 保存扣减商品和数量信息
	var inventoryList []*inventoryPb.GoodsInvInfo

	//订单的总金额 = 所有商品的金额加一起 (商品的价格（goods）*购物车商品的数量(shopCarts.nums))
	var sum float32
	//循环商品信息  进行总价计算
	//定义一个切片 保存订单购买的所有商品
	var orderList []*model.OrderGoods

	for _, info := range goods.Data {

		sum += (info.ShopPrice * float32(goodsNumsMap[info.Id]))

		orderList = append(orderList, &model.OrderGoods{
			BaseModel:  model.BaseModel{},
			GoodsId:    info.Id,
			GoodsName:  info.Name,
			GoodsPrice: info.ShopPrice,
			GoodsImage: info.GoodsFrontImage,
			Nums:       goodsNumsMap[info.Id],
		})

		inventoryList = append(inventoryList, &inventoryPb.GoodsInvInfo{
			GoodsId: info.Id,
			Num:     goodsNumsMap[info.Id],
		})
	}

	//尝试扣减库存  是增加冻结库存
	_, err = global.InventoryClient.TrySell(ctx, &inventoryPb.SellInfo{GoodsInfo: inventoryList})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	//开启事务
	begin := global.Db.Begin()

	//创建订单基本信息表
	orderInfo := model.OrderInfo{
		UserID:       req.UserId,
		OrderSn:      GenerateOrderSn(req.UserId),
		Status:       1,
		OrderMount:   sum,
		Address:      req.Address,
		SignerName:   req.Name,
		SingerMobile: req.Mobile,
		Post:         req.Post,
	}

	create := begin.Model(&model.OrderInfo{}).Create(&orderInfo)

	if create.Error != nil {
		begin.Rollback()
		//取消冻结扣减库存
		_, err = global.InventoryClient.CancelSell(ctx, &inventoryPb.SellInfo{GoodsInfo: inventoryList})
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		return nil, status.Errorf(codes.Internal, create.Error.Error())
	}
	if create.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "创建订单失败")
	}

	//创建订单商品表
	for _, orderGoods := range orderList {
		orderGoods.OrderId = orderInfo.ID
	}
	//批量添加 第二个参数 允许同时添加多少个  这样写可以省去循环
	batches := begin.CreateInBatches(orderList, 100)

	if batches.Error != nil || batches.RowsAffected == 0 {
		begin.Rollback()
		//取消冻结扣减库存
		_, err = global.InventoryClient.CancelSell(ctx, &inventoryPb.SellInfo{GoodsInfo: inventoryList})
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		return nil, status.Errorf(codes.NotFound, "创建订单商品表失败")
	}

	//进行购物车删除

	db := global.Db.Model(&model.ShoppingCart{}).Where("user_id =? and checked= true ", req.UserId).Delete(&model.ShoppingCart{})

	if db.Error != nil || db.RowsAffected == 0 {

		begin.Rollback()
		//取消冻结扣减库存
		_, err = global.InventoryClient.CancelSell(ctx, &inventoryPb.SellInfo{GoodsInfo: inventoryList})
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "购物车删除失败")
	}

	begin.Commit()

	//确认扣减库存
	_, err = global.InventoryClient.ConfirmSell(ctx, &inventoryPb.SellInfo{
		GoodsInfo: inventoryList,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &orderPb.OrderInfoResp{
		Id:      orderInfo.ID,
		UserId:  orderInfo.UserID,
		Status:  int32(orderInfo.Status),
		Post:    orderInfo.Post,
		Total:   orderInfo.OrderMount,
		Address: orderInfo.Address,
		Name:    orderInfo.SignerName,
		Mobile:  orderInfo.SingerMobile,
		OrderSn: orderInfo.OrderSn,
	}, nil
}*/

// CartItemList 购物车列列表
func (OrderService) CartItemList(ctx context.Context, req *orderPb.UserInfo) (*orderPb.CartItemListResp, error) {
	var shoppingCarts []*model.ShoppingCart
	tx := global.Db.Model(&model.ShoppingCart{}).Where(&model.ShoppingCart{UserID: req.Id}).Find(&shoppingCarts)

	if tx.Error != nil {
		return nil, status.Errorf(codes.Internal, "查询失败")
	}
	var count int64
	db := global.Db.Model(&model.ShoppingCart{}).Where(&model.ShoppingCart{UserID: req.Id}).Count(&count)
	if db.Error != nil {
		return nil, status.Errorf(codes.Internal, "获取条数失败")
	}

	var list []*orderPb.ShopCartInfoResp

	for _, cart := range shoppingCarts {
		list = append(list, &orderPb.ShopCartInfoResp{
			Id:      cart.ID,
			UserId:  cart.UserID,
			GoodsId: cart.GoodsID,
			Nums:    cart.Nums,
			Checked: cart.Checked,
		})
	}
	return &orderPb.CartItemListResp{
		Total: int32(count),
		Info:  list,
	}, nil
}

// CreateCartItem 添加购物车  多加一个条数限制
func (OrderService) CreateCartItem(ctx context.Context, req *orderPb.CartItemReq) (*orderPb.ShopCartInfoResp, error) {
	//进行购物车表查询 当前用户是否存在此商品
	var cart model.ShoppingCart

	tx := global.Db.Model(&model.ShoppingCart{}).Where("goods_id = ? and user_id = ?", req.GoodsId, req.UserId).Find(&cart)
	//查询失败
	if tx.Error != nil {
		return nil, status.Errorf(codes.Internal, "查询失败")
	}
	if tx.RowsAffected > 0 {
		//调用库存表查询库存信息
		detail, err := global.InventoryClient.InvDetail(ctx, &inventoryPb.GoodsInvInfo{GoodsId: cart.GoodsID})
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}

		//进行判断要加入的数量是否超过当前库存
		if cart.Nums+req.Nums > detail.Num {
			return nil, status.Errorf(codes.Internal, "库存已达上限")
		}
		//进行累加  商品存在  数量进行累加
		db := global.Db.Model(&model.ShoppingCart{}).Where("goods_id = ? and user_id =?", req.GoodsId, req.UserId).Updates(&model.ShoppingCart{Nums: req.Nums + cart.Nums})

		if db.Error != nil || db.RowsAffected == 0 {
			return nil, status.Errorf(codes.Internal, "更新失败"+db.Error.Error())
		}

		return &orderPb.ShopCartInfoResp{
			Id:      cart.ID,
			UserId:  cart.UserID,
			GoodsId: cart.GoodsID,
			Nums:    req.Nums + cart.Nums,
			Checked: cart.Checked,
		}, nil
	} else {
		cart = model.ShoppingCart{
			UserID:  req.UserId,
			GoodsID: req.GoodsId,
			Nums:    req.Nums,
			Checked: req.Checked,
		}
		//进行添加
		db := global.Db.Model(&model.ShoppingCart{}).Create(&cart)
		if db.Error != nil || db.RowsAffected == 0 {
			return nil, status.Errorf(codes.Internal, "购物车添加失败")
		}
		return &orderPb.ShopCartInfoResp{
			Id:      cart.ID,
			UserId:  cart.UserID,
			GoodsId: cart.GoodsID,
			Nums:    cart.Nums,
			Checked: cart.Checked,
		}, nil
	}

}

// UpdateCartItem 修改购物车
func (OrderService) UpdateCartItem(ctx context.Context, req *orderPb.CartItemReq) (*orderPb.OrderEmpty, error) {
	// 查询购物车信息
	var cart model.ShoppingCart
	tx := global.Db.Model(&model.ShoppingCart{}).Where("id = ? AND user_id = ?", req.Id, req.UserId).First(&cart)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "购物车不存在")
		}
		return nil, status.Errorf(codes.Internal, "获取购物车失败: %v", tx.Error)
	}

	// 更新购物车项的选中状态和数量
	cart.Nums = req.Nums
	cart.Checked = req.Checked
	db := global.Db.Save(&cart)
	if db.Error != nil {
		return nil, status.Errorf(codes.Internal, "修改购物车信息失败: %v", db.Error)
	}

	//var cart model.ShoppingCart
	//
	////查询购物车信息
	//tx := global.Db.Model(&model.ShoppingCart{}).Where(&model.ShoppingCart{
	//	BaseModel: model.BaseModel{ID: req.Id},
	//	UserID:    req.UserId,
	//}).Limit(1).Find(&cart)
	//
	//if tx.Error != nil || tx.RowsAffected == 0 {
	//	return nil, status.Errorf(codes.Internal, "获取购物车失败")
	//}
	//cart = model.ShoppingCart{
	//	Nums:    req.Nums,
	//	Checked: req.Checked,
	//}
	//
	////进行修改 主要是修改选中状态和数量
	//db := global.Db.Model(&model.ShoppingCart{}).Where(&model.ShoppingCart{
	//	BaseModel: model.BaseModel{ID: req.Id},
	//}).Select("nums", "checked").Updates(&cart)
	//
	//if db.Error != nil || db.RowsAffected == 0 {
	//	return nil, status.Errorf(codes.Internal, "修改购物车信息失败")
	//}

	return &orderPb.OrderEmpty{}, nil
}

// 删除购物车
func (OrderService) DeleteCart(ctx context.Context, req *orderPb.DeleteCartReq) (*orderPb.OrderEmpty, error) {

	// 查询购物车信息
	var cart model.ShoppingCart
	tx := global.Db.Model(&model.ShoppingCart{}).Where("id = ? AND user_id = ?", req.Id, req.UserId).First(&cart)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "购物车不存在")
		}
		return nil, status.Errorf(codes.Internal, "获取购物车失败: %v", tx.Error)
	}
	//删除购物车信息
	db := global.Db.Model(&model.ShoppingCart{}).Where("user_id = ? and checked = true", req.UserId).Delete(&model.ShoppingCart{})

	if db.Error != nil {
		return nil, status.Errorf(codes.Internal, db.Error.Error())
	}

	if db.RowsAffected == 0 {
		return nil, status.Errorf(codes.Internal, "删除失败")
	}

	return &orderPb.OrderEmpty{}, nil
}

// CreateOrder 创建订单  普通版
//
//	func (OrderService) CreateOrder(ctx context.Context, req *orderPb.OrderReq) (*orderPb.OrderInfoResp, error) {
//		/*
//		   新建订单
//		       1. 从购物车中获取到选中的商品
//		       2. 商品的价格自己查询 - 访问商品服务 (跨微服务)
//		       3. 库存的扣减 - 访问库存服务 (跨微服务)
//		       4. 订单的基本信息表 - 订单的商品信息表
//		       5. 从购物车中删除已购买的记录
//		*/
//
//		//定义一个切片保存购物车商品数据
//		var shopCart []*model.ShoppingCart
//		//定义一个切片存放购物车下的商品id
//		var goodsIds []int32
//
//		//定义商品数据字典  里面key  是商品id  value 是商品数量
//		goodsNumsMap := make(map[int32]int32)
//
//		//查询购物车中用户选中的商品
//		tx := global.Db.Model(&model.ShoppingCart{}).Where("user_id = ? and checked = true", req.UserId).Find(&shopCart)
//
//		if tx.Error != nil {
//			return nil, status.Errorf(codes.Internal, "查询购物车商品失败")
//		}
//		if tx.RowsAffected == 0 {
//			return nil, status.Errorf(codes.NotFound, "购物车为空")
//		}
//		for _, cart := range shopCart {
//			goodsIds = append(goodsIds, cart.GoodsID)
//			goodsNumsMap[cart.GoodsID] = cart.Nums
//		}
//
//		//去商品微服务批量查询商品信息
//		goods, err := global.GoodsClient.BatchGetGoods(context.Background(), &goodsPb.BatchGoodsIdInfo{Id: goodsIds})
//		if err != nil {
//			return nil, status.Errorf(codes.NotFound, "没有此商品信息")
//		}
//
//		//定义一个切片 保存扣减商品和数量信息
//		var inventoryList []*inventoryPb.GoodsInvInfo
//
//		//订单的总金额 = 所有商品的金额加一起 (商品的价格（goods）*购物车商品的数量(shopCarts.nums))
//		var sum float32
//		//循环商品信息  进行总价计算
//		//定义一个切片 保存订单购买的所有商品
//		var orderList []*model.OrderGoods
//
//		for _, info := range goods.Data {
//
//			sum += (info.ShopPrice * float32(goodsNumsMap[info.Id]))
//
//			orderList = append(orderList, &model.OrderGoods{
//				BaseModel:  model.BaseModel{},
//				GoodsId:    info.Id,
//				GoodsName:  info.Name,
//				GoodsPrice: info.ShopPrice,
//				GoodsImage: info.GoodsFrontImage,
//				Nums:       goodsNumsMap[info.Id],
//			})
//
//			inventoryList = append(inventoryList, &inventoryPb.GoodsInvInfo{
//				GoodsId: info.Id,
//				Num:     goodsNumsMap[info.Id],
//			})
//		}
//
//		//进行库存扣减
//		_, err = global.InventoryClient.Sell(ctx, &inventoryPb.SellInfo{
//			GoodsInfo: inventoryList,
//		})
//		if err != nil {
//			return nil, status.Errorf(codes.NotFound, err.Error())
//		}
//		//开启事务
//		begin := global.Db.Begin()
//
//		//创建订单基本信息表
//		orderInfo := model.OrderInfo{
//			UserID:       req.UserId,
//			OrderSn:      GenerateOrderSn(req.UserId),
//			Status:       1,
//			OrderMount:   sum,
//			Address:      req.Address,
//			SignerName:   req.Name,
//			SingerMobile: req.Mobile,
//			Post:         req.Post,
//		}
//
//		create := begin.Model(&model.OrderInfo{}).Create(&orderInfo)
//
//		if create.RowsAffected == 0 || create.Error != nil {
//			begin.Rollback()
//			return nil, status.Errorf(codes.NotFound, "创建订单失败")
//		}
//
//		//创建订单商品表
//		for _, orderGoods := range orderList {
//			orderGoods.OrderId = orderInfo.ID
//		}
//		//批量添加 第二个参数 允许同时添加多少个  这样写可以省去循环
//		batches := begin.CreateInBatches(orderList, 100)
//
//		if batches.Error != nil || batches.RowsAffected == 0 {
//			begin.Rollback()
//			return nil, status.Errorf(codes.NotFound, "创建订单商品表失败")
//		}
//
//		//进行购物车删除
//
//		db := global.Db.Model(&model.ShoppingCart{}).Where("user_id =? and checked= true ", req.UserId).Delete(&model.ShoppingCart{})
//
//		if db.Error != nil || db.RowsAffected == 0 {
//
//			begin.Rollback()
//
//			return nil, status.Errorf(codes.Internal, "购物车删除失败")
//		}
//
//		begin.Commit()
//
//		return &orderPb.OrderInfoResp{
//			Id:      orderInfo.ID,
//			UserId:  orderInfo.UserID,
//			Status:  int32(orderInfo.Status),
//			Post:    orderInfo.Post,
//			Total:   orderInfo.OrderMount,
//			Address: orderInfo.Address,
//			Name:    orderInfo.SignerName,
//			Mobile:  orderInfo.SingerMobile,
//			OrderSn: orderInfo.OrderSn,
//		}, nil
//	}
//
// CreateOrder 创建订单  TCC版
// 创建订单TCC版
//func (OrderService) CreateOrder(ctx context.Context, req *orderPb.OrderReq) (*orderPb.OrderInfoResp, error) {
//	/*
//	   新建订单
//	       1. 从购物车中获取到选中的商品
//	       2. 商品的价格自己查询 - 访问商品服务 (跨微服务)
//	       3. 库存的扣减 - 访问库存服务 (跨微服务)
//	       4. 订单的基本信息表 - 订单的商品信息表
//	       5. 从购物车中删除已购买的记录
//	*/
//
//	//定义一个切片保存购物车商品数据
//	var shopCart []*model.ShoppingCart
//	//定义一个切片存放购物车下的商品id
//	var goodsIds []int32
//
//	//定义商品数据字典  里面key  是商品id  value 是商品数量
//	goodsNumsMap := make(map[int32]int32)
//
//	//查询购物车中用户选中的商品
//	tx := global.Db.Model(&model.ShoppingCart{}).Where("user_id = ? and checked = true", req.UserId).Find(&shopCart)
//
//	if tx.Error != nil {
//		return nil, status.Errorf(codes.Internal, "查询购物车商品失败")
//	}
//	if tx.RowsAffected == 0 {
//		return nil, status.Errorf(codes.NotFound, "购物车为空")
//	}
//	for _, cart := range shopCart {
//		goodsIds = append(goodsIds, cart.GoodsID)
//		goodsNumsMap[cart.GoodsID] = cart.Nums
//	}
//
//	//去商品微服务批量查询商品信息
//	goods, err := global.GoodsClient.BatchGetGoods(context.Background(), &goodsPb.BatchGoodsIdInfo{Id: goodsIds})
//	if err != nil {
//		return nil, status.Errorf(codes.NotFound, "没有此商品信息")
//	}
//
//	//定义一个切片 保存扣减商品和数量信息
//	var inventoryList []*inventoryPb.GoodsInvInfo
//
//	//订单的总金额 = 所有商品的金额加一起 (商品的价格（goods）*购物车商品的数量(shopCarts.nums))
//	var sum float32
//	//循环商品信息  进行总价计算
//	//定义一个切片 保存订单购买的所有商品
//	var orderList []*model.OrderGoods
//
//	for _, info := range goods.Data {
//
//		sum += (info.ShopPrice * float32(goodsNumsMap[info.Id]))
//
//		orderList = append(orderList, &model.OrderGoods{
//			BaseModel:  model.BaseModel{},
//			GoodsId:    info.Id,
//			GoodsName:  info.Name,
//			GoodsPrice: info.ShopPrice,
//			GoodsImage: info.GoodsFrontImage,
//			Nums:       goodsNumsMap[info.Id],
//		})
//
//		inventoryList = append(inventoryList, &inventoryPb.GoodsInvInfo{
//			GoodsId: info.Id,
//			Num:     goodsNumsMap[info.Id],
//		})
//	}
//
//	//尝试扣减库存  是增加冻结库存
//	_, err = global.InventoryClient.TrySell(ctx, &inventoryPb.SellInfo{GoodsInfo: inventoryList})
//	if err != nil {
//		return nil, status.Errorf(codes.Internal, err.Error())
//	}
//
//	//开启事务
//	begin := global.Db.Begin()
//
//	//创建订单基本信息表
//	orderInfo := model.OrderInfo{
//		UserID:       req.UserId,
//		OrderSn:      GenerateOrderSn(req.UserId),
//		Status:       1,
//		OrderMount:   sum,
//		Address:      req.Address,
//		SignerName:   req.Name,
//		SingerMobile: req.Mobile,
//		Post:         req.Post,
//	}
//
//	create := begin.Model(&model.OrderInfo{}).Create(&orderInfo)
//
//	if create.Error != nil {
//		begin.Rollback()
//		//取消冻结扣减库存
//		_, err = global.InventoryClient.CancelSell(ctx, &inventoryPb.SellInfo{GoodsInfo: inventoryList})
//		if err != nil {
//			return nil, status.Errorf(codes.Internal, err.Error())
//		}
//		return nil, status.Errorf(codes.Internal, create.Error.Error())
//	}
//	if create.RowsAffected == 0 {
//		return nil, status.Errorf(codes.NotFound, "创建订单失败")
//	}
//
//	//创建订单商品表
//	for _, orderGoods := range orderList {
//		orderGoods.OrderId = orderInfo.ID
//	}
//	//批量添加 第二个参数 允许同时添加多少个  这样写可以省去循环
//	batches := begin.CreateInBatches(orderList, 100)
//
//	if batches.Error != nil || batches.RowsAffected == 0 {
//		begin.Rollback()
//		//取消冻结扣减库存
//		_, err = global.InventoryClient.CancelSell(ctx, &inventoryPb.SellInfo{GoodsInfo: inventoryList})
//		if err != nil {
//			return nil, status.Errorf(codes.Internal, err.Error())
//		}
//		return nil, status.Errorf(codes.NotFound, "创建订单商品表失败")
//	}
//
//	//进行购物车删除
//
//	db := global.Db.Model(&model.ShoppingCart{}).Where("user_id =? and checked= true ", req.UserId).Delete(&model.ShoppingCart{})
//
//	if db.Error != nil || db.RowsAffected == 0 {
//
//		begin.Rollback()
//		//取消冻结扣减库存
//		_, err = global.InventoryClient.CancelSell(ctx, &inventoryPb.SellInfo{GoodsInfo: inventoryList})
//		if err != nil {
//			return nil, status.Errorf(codes.Internal, err.Error())
//		}
//		return nil, status.Errorf(codes.Internal, "购物车删除失败")
//	}
//
//	begin.Commit()
//
//	//确认扣减库存
//	_, err = global.InventoryClient.ConfirmSell(ctx, &inventoryPb.SellInfo{
//		GoodsInfo: inventoryList,
//	})
//	if err != nil {
//		return nil, status.Errorf(codes.Internal, err.Error())
//	}
//	return &orderPb.OrderInfoResp{
//		Id:      orderInfo.ID,
//		UserId:  orderInfo.UserID,
//		Status:  int32(orderInfo.Status),
//		Post:    orderInfo.Post,
//		Total:   orderInfo.OrderMount,
//		Address: orderInfo.Address,
//		Name:    orderInfo.SignerName,
//		Mobile:  orderInfo.SingerMobile,
//		OrderSn: orderInfo.OrderSn,
//	}, nil
//}

//func (OrderService) CreateOrder(ctx context.Context, req *orderPb.OrderReq) (*orderPb.OrderInfoResp, error) {
//
//	//1.获取购物车中用户选中的商品
//	var shoppingCart []*model.ShoppingCart
//
//	tx := global.Db.Model(&model.ShoppingCart{}).Where(&model.ShoppingCart{
//		UserID:  req.UserId,
//		Checked: true,
//	}).Find(&shoppingCart)
//
//	if tx.Error != nil {
//		return nil, status.Errorf(codes.Internal, tx.Error.Error())
//	}
//
//	if tx.RowsAffected == 0 {
//		return nil, status.Errorf(codes.NotFound, "购物车为空")
//	}
//
//	var goodsIds []int32
//
//	var goodsInvInfo []*inventoryPb.GoodsInvInfo
//
//	var goodsIdInfo = make(map[int32]int32)
//	//把商品id放入切片中  用来做批量查询
//	for _, cart := range shoppingCart {
//
//		goodsIds = append(goodsIds, cart.GoodsID)
//
//		goodsInvInfo = append(goodsInvInfo, &inventoryPb.GoodsInvInfo{
//			GoodsId: cart.GoodsID,
//			Num:     cart.Nums,
//		})
//
//		goodsIdInfo[cart.GoodsID] = cart.Nums
//
//	}
//	//2.批量从商品表里面进行查询商品详情  调用商品微服务
//	goods, err := global.GoodsClient.BatchGetGoods(ctx, &goodsPb.BatchGoodsIdInfo{Id: goodsIds})
//	if err != nil {
//		return nil, status.Errorf(codes.Internal, err.Error())
//	}
//
//	//进行总价计算
//	var sum float32
//	var orderGoods []*model.OrderGoods
//
//	for _, info := range goods.Data {
//		sum += info.ShopPrice * float32(goodsIdInfo[info.Id])
//
//		orderGoods = append(orderGoods, &model.OrderGoods{
//			GoodsId:    info.Id,
//			GoodsName:  info.Name,
//			GoodsPrice: info.ShopPrice,
//			Nums:       goodsIdInfo[info.Id],
//		})
//
//	}
//	//进行扣减库存
//	_, err = global.InventoryClient.Sell(ctx, &inventoryPb.SellInfo{
//		GoodsInfo: goodsInvInfo,
//	})
//	if err != nil {
//		return nil, status.Errorf(codes.Internal, err.Error())
//	}
//
//	//进行订单表的创建
//	var orderInfo = model.OrderInfo{
//		UserID:       req.UserId,
//		OrderSn:      GenerateOrderSn(req.UserId),
//		PayType:      1,
//		Status:       1,
//		OrderMount:   sum,
//		Address:      req.Address,
//		SignerName:   req.Name,
//		SingerMobile: req.Mobile,
//		Post:         req.Post,
//	}
//
//	//创建订单商品表  批量添加
//	for _, info := range orderGoods {
//		info.OrderId = orderInfo.ID
//	}
//
//	batches := global.Db.CreateInBatches(&orderGoods, 100)
//
//	if batches.Error != nil {
//		return nil, status.Errorf(codes.Internal, batches.Error.Error())
//	}
//
//	return &orderPb.OrderInfoResp{}, nil
//}

// 计算页数
func GetPage(page int, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		if limit == 0 {
			limit = 10
		}
		if limit > 10 {
			limit = 100
		}
		//计算偏移量
		offset := (page - 1) * limit

		return db.Offset(offset).Limit(limit)
	}

}

// 订单列表
func (OrderService) OrderList(ctx context.Context, req *orderPb.OrderFilterReq) (*orderPb.OrderListResp, error) {
	var list []*model.OrderInfo
	//查询用户的所有订单
	tx := global.Db.Model(&model.OrderInfo{}).Where("user_id = ?", req.UserId)
	if req.PayType != 0 {
		//根据订单支付方式去查询
		tx = tx.Where(&model.OrderInfo{PayType: int8(req.PayType)})
	}

	if req.Status != 0 {
		//根据订单状态去查询
		tx = tx.Where("status = ?", req.Status)
	}
	//符合条件的个数
	var count int64
	tx.Count(&count)
	tx = tx.Scopes(GetPage(int(req.Pages), int(req.PagePerNums))).Find(&list)

	var orderInfoList []*orderPb.OrderInfoResp

	for _, info := range list {
		orderInfoList = append(orderInfoList, &orderPb.OrderInfoResp{
			Id:      info.ID,
			UserId:  info.UserID,
			OrderSn: info.OrderSn,
			PayType: int32(info.PayType),
			Status:  int32(info.Status),
			Post:    info.Post,
			Total:   info.OrderMount,
			Address: info.Address,
			Name:    info.SignerName,
			Mobile:  info.SingerMobile,
		})
	}
	return &orderPb.OrderListResp{
		Total: int32(count),
		Data:  orderInfoList,
	}, nil
}

// OrderDetail 订单详情
func (OrderService) OrderDetail(ctx context.Context, req *orderPb.OrderReq) (*orderPb.OrderInfoDetailResp, error) {
	var orderInfo model.OrderInfo
	//根据订单id获取订单信息
	//查询订单表
	tx := global.Db.Model(&model.OrderInfo{}).Where("order_sn = ?", req.OrderSn).Limit(1).Find(&orderInfo)

	if tx.Error != nil {
		return nil, status.Errorf(codes.Internal, tx.Error.Error())
	}

	if tx.RowsAffected == 0 {
		return nil, status.Errorf(codes.Internal, "订单不存在")
	}
	var orderInfoResp = orderPb.OrderInfoResp{
		Id:      orderInfo.ID,
		UserId:  orderInfo.UserID,
		OrderSn: orderInfo.OrderSn,
		PayType: int32(orderInfo.PayType),
		Status:  int32(orderInfo.Status),
		Post:    orderInfo.Post,
		Total:   orderInfo.OrderMount,
		Address: orderInfo.Address,
		Name:    orderInfo.SignerName,
		Mobile:  orderInfo.SingerMobile,
	}

	var orderGoods []*model.OrderGoods
	//查询订单商品表
	find := global.Db.Model(&model.OrderGoods{}).Where("order_id = ?", orderInfo.ID).Find(&orderGoods)

	if find.Error != nil {
		return nil, status.Errorf(codes.Internal, tx.Error.Error())
	}

	if find.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单商品不存在")
	}
	var orderItem []*orderPb.OrderItemResponse

	for _, info := range orderGoods {
		orderItem = append(orderItem, &orderPb.OrderItemResponse{
			Id:         info.ID,
			OrderId:    info.OrderId,
			GoodsId:    info.GoodsId,
			GoodsName:  info.GoodsName,
			GoodsImage: info.GoodsImage,
			GoodsPrice: info.GoodsPrice,
			Nums:       info.Nums,
		})
	}

	return &orderPb.OrderInfoDetailResp{
		OrderInfo: &orderInfoResp,
		Goods:     orderItem,
	}, nil
}

// UpdateOrder 更新订单
func (OrderService) UpdateOrder(ctx context.Context, req *orderPb.UpdateOrderInfo) (*orderPb.OrderEmpty, error) {

	var orderInfo model.OrderInfo
	//查询是否有订单信息
	tx := global.Db.Model(&model.OrderInfo{}).Where("order_sn= ?", req.OrderSn).Limit(1).Find(&orderInfo)
	if tx.Error != nil {
		return nil, status.Errorf(codes.Internal, tx.Error.Error())
	}

	if tx.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}
	orderInfo = model.OrderInfo{
		PayType: int8(req.PayType),
		Status:  int8(req.Status),
		TradeNo: req.TradeNo,
		PayTime: req.PayTime,
	}

	updates := global.Db.Model(&model.OrderInfo{}).Where("order_sn= ?", req.OrderSn).Updates(&orderInfo)

	if updates.Error != nil {
		return nil, status.Errorf(codes.Internal, tx.Error.Error())
	}

	if updates.RowsAffected == 0 {
		return nil, status.Errorf(codes.Internal, "订单修改失败")
	}
	return &orderPb.OrderEmpty{}, nil
}

// GenerateOrderSn 生成订单编号
func GenerateOrderSn(userId int32) string {
	now := time.Now().Format("20060102150405")

	intn := rand.Intn(100)

	orderSn := now + strconv.Itoa(int(userId)) + strconv.Itoa(intn)
	return orderSn

}

type OrderTimeoutListener struct {
	Code        codes.Code
	Detail      string
	ID          int32
	OrderAmount float32
	Ctx         context.Context
}

// DelayProducer 延迟队列生产者
func DelayProducer() {

}

// OrderTimeout 订单超时
func OrderTimeout(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

	//把msgs的信息反解到订单结构体

	for _, msg := range msgs {
		var orderInfo model.OrderInfo
		//将延迟队列里的消息反解到orderInfo结构体
		err := json.Unmarshal(msg.Body, &orderInfo)
		if err != nil {
			zap.S().Info("反序列化失败")
			return 0, err
		}

		//查询订单号是否存在
		tx := global.Db.Model(&model.OrderInfo{}).Where(&model.OrderInfo{OrderSn: orderInfo.OrderSn}).Limit(1).Find(&orderInfo)

		if tx.Error != nil {
			zap.S().Info("订单不存在")
			return consumer.ConsumeRetryLater, nil
		}

		//查询订单状态
		if orderInfo.Status == 1 {
			//待支付
			orderInfo.Status = 3 //3 代表交易关闭
			//更改状态
			save := global.Db.Save(&orderInfo)

			if save.Error != nil {
				return consumer.ConsumeRetryLater, nil
			}
			if save.RowsAffected == 0 {
				return consumer.ConsumeRetryLater, nil
			}
			//回归库存
			_, err = global.InventoryClient.Reback(ctx, &inventoryPb.SellInfo{
				OrderSn: orderInfo.OrderSn,
			})
			if err != nil {
				return consumer.ConsumeRetryLater, nil
			}
		} else if orderInfo.Status == 2 {
			return consumer.ConsumeSuccess, nil
		}

	}

	return consumer.ConsumeSuccess, nil
}
