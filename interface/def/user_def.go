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
	Account     string `json:"account"`      // 账号（手机号或邮箱）
	AccountType string `json:"account_type"` // 账号类型：phone（手机号）或 email（邮箱）
	Password    string `json:"password"`     // 密码
}

type LoginResp struct {
	Token    string `json:"token,omitempty"`     // JWT token
	UserID   string `json:"user_id,omitempty"`   // 用户ID
	UserName string `json:"user_name,omitempty"` // 用户名
	Avatar   string `json:"avatar,omitempty"`    // 头像
	Phone    string `json:"phone,omitempty"`     // 手机号
	Email    string `json:"email,omitempty"`     // 邮箱
	Success  bool   `json:"success"`             // 登录是否成功
}

//---------注册相关------------
// 注册：用户名 + 手机号/邮箱 + 验证码 + 设置密码
type RegisterReq struct {
	UserName    string `json:"user_name"`
	Account     string `json:"account"`
	AccountType string `json:"account_type"` // 手机号或邮箱
	Code        string `json:"code"`
	Password    string `json:"password"`
}

type RegisterResp struct {
	Success bool `json:"success"` // 注册是否成功
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

//---------更新头像-----------
type UpdateAvatarReq struct {
	AvatarURL string `json:"avatar_url" binding:"required"` // 头像URL
}

type UpdateAvatarResp struct {
	Success bool `json:"success"` // 更新是否成功
}

//---------发送验证码-----------
type SendVerificationCodeReq struct {
	Account     string `json:"account"`      // 账号（手机号或邮箱）  目前只支持邮箱 邮件收取验证码
	AccountType string `json:"account_type"` // 账号类型：phone（手机号）或 email（邮箱）
}

type SendVerificationCodeResp struct {
	Success bool `json:"success"` // 发送是否成功
}

//---------第三方--------- 暂时先不做
