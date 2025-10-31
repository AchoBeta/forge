package def

// 这个是DTO层，会暴露给前端 主要是接口定义

type User struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	Avatar   string `json:"avatar,omitempty"`

	Phone string `json:"phone,omitempty"`
	Email string `json:"email,omitempty"`

	Dogs []*Dog `json:"dogs"`
}

type Dog struct {
	DogID   string `json:"dog_id"`
	DogName string `json:"dog_name"`
}

// ---------登录相关----------
type LoginReq struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
type LoginResp struct {
	User *User `json:"user"`
}

//---------注册相关------------
// 注册：用户名 + 手机号/邮箱 + 验证码 + 设置密码
type RegisterReq struct {
	UserName    string `json:"user_name,omitempty"`
	Account     string `json:"account"`
	AccountType string `json:"account_type"` // 手机号或邮箱
	Code        string `json:"code"`
	Password    string `json:"password"`
}

type RegisterResp struct {
	// 只需返回成功信息
}

//---------重置密码-----------
type ResetPasswordReq struct {
	Account         string `json:"account"`
	AccountType     string `json:"account_type"` // 手机号或邮箱
	Code            string `json:"code"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

type ResetPasswordResp struct {
	Success bool `json:"success"`
}

//---------第三方--------- 暂时先不做
