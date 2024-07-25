package config

type UserServer struct {
	Name  string
	Port  int
	Mysql Mysql `yaml:"mysql"`
}

type Mysql struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Db       string `yaml:"db"`
}
