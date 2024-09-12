package global

import (
	"gorm.io/gorm"
	"userop_srv/config"
)

var (
	ServerID         string
	UserServerConfig config.UserServer
	Db               *gorm.DB
)
