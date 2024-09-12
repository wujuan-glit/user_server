package model

type CartGoodsInfo struct {
	Id        int32   `json:"id"`
	UserId    int32   `json:"user_id"`
	GoodsId   int32   `json:"goods_id"`
	GoodName  string  `json:"good_name"`
	GoodImage string  `json:"good_image"`
	GoodPrice float32 `json:"good_price"`
	Nums      int32   `json:"nums"`
	Checked   bool    `json:"checked"`
}
