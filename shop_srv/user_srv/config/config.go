package config

type UserServer struct {
	Name   string `json:"name"`
	Port   int    `json:"port"`
	Ip     string `json:"ip"`
	Mysql  Mysql  `json:"mysql"`
	Consul Consul `json:"consul"`
	Nacos  Nacos  `json:"nacos"`
}

type Mysql struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Db       string `json:"db"`
}
type Consul struct {
	Host string   `json:"host"`
	Port int      `json:"port"`
	Tags []string `json:"tags"`
}
type Nacos struct {
	Namespace      string ` mapstructure:"namespace" json:"namespace"`
	ServerPort     uint64 ` mapstructure:"server_port" json:"server_port"`
	ServerIp       string ` mapstructure:"server_ip" json:"server_ip"`
	ServerGrpcPort uint64 ` mapstructure:"server_grpc_port" json:"server_grpc_port"`
	DataId         string ` mapstructure:"data_id" json:"data_id"`
	Group          string ` mapstructure:"group" json:"group"`
}
