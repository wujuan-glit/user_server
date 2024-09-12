package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"

	"time"
)

type BaseModel struct {
	ID        int32 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt sql.NullTime   //允许为空
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type User struct {
	BaseModel
	Password string `gorm:"type:varchar(150);not null;comment:密码"`
	NickName string `gorm:"type:varchar(50);not null;comment:昵称"`
	Mobile   string `gorm:"type:varchar(15);not null;comment:手机号"`
	Role     int32  `gorm:"column:role;type:tinyint(1);not null;comment:角色"`
	Birthday uint64 `gorm:"column:birthday;type:int(10);not null;comment:生日"`
	Gender   string `gorm:"type:varchar(2);not null;comment:性别"`
	Pid      int32  `gorm:"int(11);comment:上级id"` //用来做二级营销
}

func main() {

	db, err := gorm.Open(mysql.Open("root:123456@tcp(127.0.0.1:3306)/user?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		log.Println("数据库连接失败")
		return
	}

	r := gin.Default()
	r.GET("/list", func(c *gin.Context) {
		var userID []*User
		var userPID []*User

		//pid := c.Query("pid")
		id := c.Query("id")

		//pidInt, _ := strconv.Atoi(pid)
		idInt, _ := strconv.Atoi(id)

		var subQuery *gorm.DB

		//根据用户id查询出所有的id
		subQuery = db.Model(&User{}).Where("id = ?", idInt).Find(&userID)
		var ids []int32
		var uid int32
		for _, info := range userID {
			uid = info.Pid
			ids = append(ids, info.ID)

		}
		subQuery = db.Model(&User{}).Where("pid in (?) ", ids).Find(&userPID)

		var pids []int32
		for _, user := range userPID {
			pids = append(pids, user.ID)
		}

		if subQuery.RowsAffected == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "数据不存在"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"msg": "success",
			"data": map[string]interface{}{
				"up_ids":   uid,
				"down_ids": pids,
			},
		})
	})

	r.Run("127.0.0.1:8888")

}

// GetUserIdByPid 根据用户id查询出所有的pid
func GetUserIdByPid() {

}

// GetUserPidById 根据用户pid查询出用户id
func GetUserPidById() {

}
