package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
	"user/proto"
	"user/user_buf/forms"
	"user/user_buf/global"
	"user/user_buf/model"
	"user/user_buf/server"
)

//我们bff层需要遵循restful风格
//路由尽量不要出现动词,出现的是名词
// get（获取数据） post（添加数据） put（修改数据） delete（删除数据）
// httpcode(http状态码) code（自定义错误码） message（提示消息） data（返回数据） 接口返回数据三要素

// RegisterUser  用户注册接口
func RegisterUser(c *gin.Context) {
	//接收参数
	var form forms.RegisterUserForm
	//校验参数 因为我们gin自带的验证器不好用，所以我们需要自定义封装验证器，符合国人的简洁中文提示

	if err := c.ShouldBind(&form); err != nil {
		ReturnErrorJson(c, err)
		return
	}

	//调用微服务rpc
	register, err := global.ServerConn.Register(context.Background(), &proto.RegisterReq{
		Mobile:   form.Mobile,
		Password: form.Password,
		Nickname: form.NickName,
	})
	if err != nil {
		HandleGrpcErrorToHttp(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    1000,
		"message": "注册成功",
		"data":    register,
	})
}

// 更新用户
func UpdateUserInfo(c *gin.Context) {
	var form forms.UpdateUserForm

	if err := c.ShouldBind(&form); err != nil {
		ReturnErrorJson(c, err)
		return
	}

	userId, _ := c.Get("user_id")

	birthday, _ := strconv.Atoi(form.Birthday)

	info, err := global.ServerConn.UpdateUserInfo(context.Background(), &proto.UpdateUserInfoReq{
		Id:       userId.(int32),
		NickName: form.NickName,
		Gender:   form.Gender,
		BirthDay: uint64(birthday),
	})

	if err != nil {
		HandleGrpcErrorToHttp(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    1000,
		"message": "更新成功",
		"data":    info,
	})

}

// 获取用户列表
func GetUserList(c *gin.Context) {
	page, limit, offset := GetPage(c)

	list, err := global.ServerConn.GetUserList(context.Background(), &proto.PageInfo{
		Pn:     uint32(page),
		PSize:  uint32(limit),
		Offset: uint32(offset),
	})
	if err != nil {
		HandleGrpcErrorToHttp(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "列表获取成功",
		"list":    list.UserInfo,
		"total":   list.Total,
	})
}

// 用户登录
func Login(c *gin.Context) {
	//接收参数
	var form forms.LoginUserForm
	//校验参数 因为我们gin自带的验证器不好用，所以我们需要自定义封装验证器，符合国人的简洁中文提示

	if err := c.ShouldBind(&form); err != nil {
		ReturnErrorJson(c, err)
		return
	}
	//校验图形验证码是否正确
	//verify := store.Verify(form.CaptchaId, form.Captcha, false)
	//if !verify {
	//	c.JSON(http.StatusOK, gin.H{
	//		"code": 0,
	//		"msg":  "图形验证码不正确",
	//		"data": nil,
	//	})
	//	return
	//}
	info, err := global.ServerConn.GetUserByMobile(context.Background(), &proto.MobileReq{Mobile: form.Mobile})
	if err != nil {
		HandleGrpcErrorToHttp(c, err)
		return
	}

	//判断密码是否正确
	check_is, err := global.ServerConn.CheckPassword(context.Background(), &proto.CheckPasswordReq{
		Password:          form.Password,
		EncryptedPassword: info.Password,
	})

	if err != nil {
		HandleGrpcErrorToHttp(c, err)
		return
	}

	if !check_is.Success {
		err = errors.New("密码错误")
		ReturnErrorJson(c, err)
		return

	}

	user := model.User{
		ID:        int32(info.Id),
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		DeletedAt: time.Time{},
		Password:  info.Password,
		NickName:  info.UserName,
		Mobile:    info.Mobile,
		Role:      info.Role,
		Birthday:  info.Birthday,
		Gender:    info.Gender,
	}
	marshal, err := json.Marshal(user)
	if err != nil {
		ReturnErrorJson(c, err)
		return
	}
	//生成token

	token, err := server.GenJwtToken(string(marshal))
	if err != nil {
		ReturnErrorJson(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":                      0,
		"msg":                       "登录成功",
		"access_token":              token.AccessToken,
		"refresh_token":             token.RefreshToken,
		"access_token_expire_time":  global.ServerConfig.Jwt.AccessExpire,
		"refresh_token_expire_time": global.ServerConfig.Jwt.RefreshExpire,
	})
}
func ReFresh(c *gin.Context) {
	refresh_token := c.PostForm("refresh_token")

	token, err := server.RefreshAccessToken(refresh_token)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":                      0,
		"msg":                       "登录成功",
		"access_token":              token,
		"refresh_token_expire_time": global.ServerConfig.Jwt.RefreshExpire,
	})
}
