package model

import "time"

type Goods struct {
	Id              int32    `json:"id"`
	CategoryId      int32    `json:"category_id"`
	Name            string   `json:"name"`
	ClickNum        int32    `json:"click_num"`
	SoldNum         int32    `json:"sold_num"`
	FavNum          int32    `json:"fav_num"`
	MarketPrice     float64  `json:"market_price"`
	ShopPrice       float64  `json:"shop_price"`
	GoodsBrief      string   `json:"goods_brief"`
	ShipFree        bool     `json:"ship_free"`
	Images          []string `json:"images"`
	DescImages      []string `json:"desc_images"`
	GoodsFrontImage string   `json:"goods_front_image"`
	OnSale          bool     `json:"on_sale"`
	Category        struct {
		Id   int32  `json:"id"`
		Name string `json:"name"`
	} `json:"category"`
	Brand struct {
		Id   int32  `json:"id"`
		Name string `json:"name"`
		Logo string `json:"logo"`
	} `json:"brand"`
}

type Test struct {
	ID        int32     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt struct {
		Time  time.Time `json:"Time"`
		Valid bool      `json:"Valid"`
	} `json:"updated_at"`
	DeletedAt        interface{} `json:"deleted_at"`
	Name             string      `json:"name"`
	ParentCategoryID int32       `json:"ParentCategoryID"`
	SubCategory      []struct {
		ID        int32     `json:"ID"`
		CreatedAt time.Time `json:"CreatedAt"`
		UpdatedAt struct {
			Time  time.Time `json:"Time"`
			Valid bool      `json:"Valid"`
		} `json:"UpdatedAt"`
		DeletedAt        interface{} `json:"DeletedAt"`
		Name             string      `json:"Name"`
		ParentCategoryID int         `json:"ParentCategoryID"`
		SubCategory      []struct {
			ID        int32     `json:"ID"`
			CreatedAt time.Time `json:"CreatedAt"`
			UpdatedAt struct {
				Time  time.Time `json:"Time"`
				Valid bool      `json:"Valid"`
			} `json:"UpdatedAt"`
			DeletedAt        interface{} `json:"DeletedAt"`
			Name             string      `json:"Name"`
			ParentCategoryID int32       `json:"ParentCategoryID"`
			SubCategory      interface{} `json:"SubCategory"`
			Level            int32       `json:"Level"`
			IsTab            bool        `json:"IsTab"`
		} `json:"SubCategory"`
		Level int32 `json:"Level"`
		IsTab bool  `json:"IsTab"`
	} `json:"SubCategory"`
	Level int32 `json:"Level"`
	IsTab bool  `json:"IsTab"`
}
type Category struct {
	Id               int    `json:"id"`
	Name             string `json:"name"`
	ParentCategoryID int    `json:"parent_category_id"`
	SubCategory      []struct {
		Id               int    `json:"id"`
		Name             string `json:"name"`
		ParentCategoryID int    `json:"parent_category_id"`
		SubCategory      []struct {
			Id               int         `json:"id"`
			Name             string      `json:"name"`
			ParentCategoryID int         `json:"parent_category_id"`
			SubCategory      interface{} `json:"sub_category"`
			Level            int         `json:"level"`
			IsTab            bool        `json:"is_tab"`
		} `json:"SubCategory"`
		Level int  `json:"level"`
		IsTab bool `json:"is_tab"`
	} `json:"SubCategory"`
	Level int  `json:"level"`
	IsTab bool `json:"is_tab"`
}
