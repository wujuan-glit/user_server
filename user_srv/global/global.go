package global

import (
	"gorm.io/gorm"
	"user/user_srv/config"
)

var (
	ServerID         string
	UserServerConfig config.UserServer
	Db               *gorm.DB
)
