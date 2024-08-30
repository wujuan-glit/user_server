package config

type UserServer struct {
	Name     string   `json:"name"`
	Port     int      `json:"port"`
	Ip       string   `json:"ip"`
	Mysql    Mysql    `json:"mysql"`
	Consul   Consul   `json:"consul"`
	Nacos    Nacos    `json:"nacos"`
	Rocketmq Rocketmq `json:"rocketmq"`
}

type Mysql struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Db       string `json:"db"`
}
type Consul struct {
	Host             string   `json:"host"`
	Port             int      `json:"port"`
	Tags             []string `json:"tags"`
	GoodsService     string   `json:"goods_service"`
	InventoryService string   `json:"inventory_service"`
}
type Nacos struct {
	Namespace      string ` mapstructure:"namespace" json:"namespace"`
	ServerPort     uint64 ` mapstructure:"server_port" json:"server_port"`
	ServerIp       string ` mapstructure:"server_ip" json:"server_ip"`
	ServerGrpcPort uint64 ` mapstructure:"server_grpc_port" json:"server_grpc_port"`
	DataId         string ` mapstructure:"data_id" json:"data_id"`
	Group          string ` mapstructure:"group" json:"group"`
}

type Rocketmq struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	RebackGroup string `json:"reback_group"`
	RebackTopic string `json:"reback_topic"`
	DelayTopic  string `json:"delay_topic"`
	DelayLevel  int    `json:"delay_level"`
	DelayGroup  string `json:"delay_group"`
}
