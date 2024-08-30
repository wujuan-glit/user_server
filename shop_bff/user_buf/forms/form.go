package forms

// 绑定为json
type RegisterUserForm struct {
	Mobile   string `form:"mobile" json:"mobile" xml:"mobile"  binding:"required,mobile"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
	NickName string `form:"nickname" json:"nickname" xml:"nickname" binding:"required"`
}

// 更新
type UpdateUserForm struct {
	//Id       string `form:"id" json:"id" xml:"id" binding:"required"`
	NickName string `form:"nickname" json:"nickname" xml:"nickname" binding:"required"`
	Gender   string `form:"gender" json:"gender" xml:"gender" binding:"required"`
	Birthday string `form:"birthday" json:"birthday" xml:"birthday" binding:"required"`
}

// 登录
type LoginUserForm struct {
	CaptchaId string `form:"id" json:"id" xml:"id" binding:"required"`
	Captcha   string `form:"captcha" json:"captcha" xml:"captcha" binding:"required"`
	Mobile    string `form:"mobile" json:"mobile" xml:"mobile"  binding:"required,mobile"`
	Password  string `form:"password" json:"password" xml:"password" binding:"required"`
}
