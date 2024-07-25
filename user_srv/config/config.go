package config

type UserServer struct {
	Name   string
	Port   int
	Ip     string
	Mysql  Mysql  `yaml:"mysql"`
	Consul Consul `yaml:"consul"`
}

type Mysql struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Db       string `yaml:"db"`
}
type Consul struct {
	Host string
	Port int
	Tags []string
}
