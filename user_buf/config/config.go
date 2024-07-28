package config

type UserServer struct {
	Name   string
	Port   int
	Ip     string
	Mysql  Mysql  `yaml:"mysql"`
	Jwt    Jwt    `yaml:"jwt"`
	Consul Consul `yaml:"consul"`
}

type Mysql struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Db       string `yaml:"db"`
}
type Jwt struct {
	Key           string `mapstructure:"key"`
	AccessExpire  int64  `mapstructure:"access_expire"`
	RefreshExpire int64  `mapstructure:"refresh_expire"`
}
type Consul struct {
	Host string
	Port int
	Name string
	Tags []string
}
