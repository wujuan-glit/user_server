package forms

type ShopCartItemForm struct {
	GoodsId int32 `form:"goods_id" json:"goods_id" xml:"goods_id"  binding:"required"`
	Num     int32 `form:"num" json:"num" xml:"num"  binding:"required"`
}

type ShopCartItemUpdateForm struct {
	Num     int32 `form:"num" json:"num" binding:"required,min=1"`
	Checked bool  ` form:"checked" json:"checked"`
	GoodsId int32 `form:"goods_id" json:"goods_id"`
	CartId  int32 ` form:"cart_id" json:"cart_id"`
}
