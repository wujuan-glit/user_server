package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
	"user/proto"
	"user/user_srv/global"
	"user/user_srv/model"
)

type UserServer struct {
	proto.UnimplementedUserServer
}

// CreateUser 创建用户
func (u *UserServer) Register(ctx context.Context, in *proto.RegisterReq) (*proto.UserInfoResp, error) {
	//接收参数
	user := model.User{}
	//先去查询要注册的用户有没有存在
	global.Db.Model(&model.User{}).Where("mobile=?", in.Mobile).First(&user)

	//如果存在了，直接提示用户已存在
	//不存在在执行添加
	if user.ID > 0 {
		return nil, status.Error(codes.AlreadyExists, "用户已存在")
	}

	//用户结构体赋值

	//密码进行加密 对称加密 非对称加密
	// Using custom options
	options := &password.Options{10, 10000, 50, sha512.New}
	salt, encodedPwd := password.Encode(in.Password, options)

	user.Mobile = in.Mobile
	user.NickName = in.Nickname
	user.Password = fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	user.CreatedAt = time.Now()
	//调用gorm.Create()创建用户
	res := global.Db.Model(&model.User{}).Create(&user)

	if res.Error != nil {
		return nil, status.Error(codes.Internal, "添加用户失败")
	}

	response := ModelToResponse(user)

	return &response, nil
}

// CheckPassWord 检查用户密码是否正确
func (u *UserServer) CheckPassword(ctx context.Context, in *proto.CheckPasswordReq) (*proto.CheckPasswordResp, error) {
	options := &password.Options{10, 10000, 50, sha512.New}

	//字符串分割
	split := strings.Split(in.EncryptedPassword, "$")
	check := password.Verify(in.Password, split[2], split[3], options)

	return &proto.CheckPasswordResp{
		Success: check,
	}, nil
}

// 页码
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// ModelToResponse 把数据的数据转化成proto想要的响应格式

func ModelToResponse(user model.User) proto.UserInfoResp { // grpc message nil

	userInfoResp := proto.UserInfoResp{
		Id:       uint32(user.ID),
		Password: user.Password,
		Gender:   user.Gender,
		Role:     string(user.Role),
		Mobile:   user.Mobile,
		Birthday: strconv.FormatUint(user.Birthday, 10),
	}

	return userInfoResp
}

// GetUserList 获取用户列表
func (u *UserServer) GetUserList(ctx context.Context, in *proto.PageInfo) (*proto.UserListResp, error) {
	zap.S().Info("列表获取")
	//链式操作

	userList := []model.User{}

	tx := global.Db.Model(&model.User{}).Scopes(Paginate(int(in.Pn), int(in.PSize))).Find(&userList)

	if tx.Error != nil {
		return nil, status.Error(codes.Internal, "数据库异常")
	}

	if tx.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户列表为空")
	}

	//查询总条数
	var count int64

	res := global.Db.Model(&model.User{}).Count(&count)

	if res.Error != nil {
		return nil, status.Error(codes.Internal, "数据库异常")
	}

	userListResponse := []*proto.UserInfoResp{}

	for _, user := range userList {
		response := ModelToResponse(user)
		userListResponse = append(userListResponse, &response)
	}

	return &proto.UserListResp{
		Total:    uint32(count),
		UserInfo: userListResponse,
	}, nil
}

// GetUserByMobile 通过手机号码获取用户信息
func (u *UserServer) GetUserByMobile(ctx context.Context, in *proto.MobileReq) (*proto.UserInfoResp, error) {
	//定义一个结构体用来保存用户的数据
	user := model.User{}

	tx := global.Db.Model(&model.User{}).Where("mobile=?", in.Mobile).First(&user)

	if tx.Error != nil {
		return nil, status.Error(codes.Internal, "用户不存在")
	}

	response := ModelToResponse(user)

	return &response, nil
}

// UpdateUser 更新用户信息
func (u *UserServer) UpdateUserInfo(ctx context.Context, in *proto.UpdateUserInfoReq) (*proto.UpdateUserInfoResp, error) {

	user := model.User{}
	tx := global.Db.Model(&model.User{}).Where("id=?", in.Id).First(&user)

	if tx.Error != nil {
		return nil, status.Error(codes.Internal, "用户不存在")
	}

	user.NickName = in.NickName
	user.Gender = in.Gender
	user.Birthday = in.BirthDay

	res := global.Db.Model(&model.User{
		BaseModel: model.BaseModel{
			ID: user.ID,
		},
	}).Save(&user)

	if res.Error != nil {
		return nil, status.Error(codes.Internal, "更新用户信息失败")
	}

	return &proto.UpdateUserInfoResp{}, nil
}
