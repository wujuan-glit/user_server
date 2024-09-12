package forms

// CreateOrderItem 创建订单表单
type CreateOrderItem struct {
	Address string `form:"address" json:"address" xml:"address"  binding:"required"`
	Name    string `form:"name" json:"name" xml:"name"  binding:"required"`
	Mobile  string `form:"mobile" json:"mobile" xml:"mobile"  binding:"required"`
	Post    string `form:"post" json:"post" xml:"post"  binding:"required"`
}

// ListOrderItem 订单列表
type ListOrderItem struct {
	OrderSn     string `form:"order_sn" json:"order_sn" xml:"order_sn"  binding:""`
	Pages       int32  `form:"pages" json:"pages" xml:"pages"  binding:""`
	PagePerNums int32  `form:"page_per_nums" json:"page_per_nums" xml:"page_per_nums"  binding:""`
	PayType     int32  `form:"pay_type" json:"pay_type" xml:"pay_type"  binding:""`
	Status      int32  `form:"status" json:"status" xml:"status"  binding:""`
}
