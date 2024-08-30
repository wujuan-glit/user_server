package config

type GoodsServer struct {
	Name   string `json:"name"`
	Port   int    `json:"port"`
	Ip     string `json:"ip"`
	PubIp  string `json:"pub_ip"`
	Jwt    Jwt    `json:"jwt"`
	Consul Consul `json:"consul"`
	Nacos  Nacos  `mapstructure:"nacos" json:"nacos"`
	Alipay Alipay `mapstructure:"alipay" json:"alipay"`
	Redis  Redis  `json:"redis"`
}

type Jwt struct {
	Key           string `mapstructure:"key" json:"key"`
	AccessExpire  int64  `mapstructure:"access_expire" json:"access_expire"`
	RefreshExpire int64  `mapstructure:"refresh_expire" json:"refresh_expire"`
}
type Consul struct {
	Host             string   `json:"host"`
	Port             int      `json:"port"`
	Name             string   `json:"name"`
	GoodsService     string   `json:"goods_service"`
	OrderService     string   `json:"order_service"`
	InventoryService string   `json:"inventory_service"`
	Tags             []string `json:"tags"`
}
type Nacos struct {
	Namespace      string `mapstructure:"namespace" json:"namespace"`
	ServerPort     uint64 `mapstructure:"server_port" json:"server_port"`
	ServerIp       string `mapstructure:"server_ip" json:"server_ip"`
	ServerGrpcPort uint64 `mapstructure:"server_grpc_port" json:"server_grpc_port"`
	DataId         string `mapstructure:"data_id" json:"data_id"`
	Group          string `mapstructure:"group" json:"group"`
}

type Alipay struct {
	Appid      string `json:"appid"`
	PrivateKey string `json:"private_key"`
}

type Redis struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}
