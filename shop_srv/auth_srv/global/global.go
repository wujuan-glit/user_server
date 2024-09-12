package global

import (
	"gorm.io/gorm"
	"user_srv/config"
)

var (
	ServerID         string
	UserServerConfig config.UserServer
	Db               *gorm.DB
)
