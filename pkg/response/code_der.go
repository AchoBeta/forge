package response

type MsgCode struct {
	Code int
	Msg  string
	Err  error
}

func (c *MsgCode) WithErr(err error) *MsgCode {
	c.Err = err
	return c
}

var (
	/* 成功 */
	SUCCESS = MsgCode{Code: 200, Msg: "成功"}

	/* 默认失败 */
	COMMON_FAIL = MsgCode{Code: -4396, Msg: "失败"}

	/* 请求错误 <0 */
	TOKEN_IS_EXPIRED = MsgCode{Code: -2, Msg: "token已过期"}

	/* 内部错误 600 ~ 999 */
	INTERNAL_ERROR             = MsgCode{Code: 601, Msg: "内部错误, check log"}
	INTERNAL_FILE_UPLOAD_ERROR = MsgCode{Code: 602, Msg: "文件上传失败"}
	/* 参数错误：1000 ~ 1999 */
	PARAM_NOT_VALID    = MsgCode{Code: 1001, Msg: "参数无效"}
	PARAM_IS_BLANK     = MsgCode{Code: 1002, Msg: "参数为空"}
	PARAM_TYPE_ERROR   = MsgCode{Code: 1003, Msg: "参数类型错误"}
	PARAM_NOT_COMPLETE = MsgCode{Code: 1004, Msg: "参数缺失"}
	INVALID_PARAMS     = MsgCode{Code: 1005, Msg: "请求体无效"}

	PARAM_FILE_SIZE_TOO_BIG = MsgCode{Code: 1010, Msg: "文件过大"}

	/* 用户错误 2000 ~ 2999 */
	USER_NOT_LOGIN             = MsgCode{Code: 2001, Msg: "用户未登录"}
	USER_PASSWORD_DIFFERENT    = MsgCode{Code: 2002, Msg: "用户两次密码输入不一致"}
	USER_ACCOUNT_NOT_EXIST     = MsgCode{Code: 2003, Msg: "账号不存在"}
	USER_CREDENTIALS_ERROR     = MsgCode{Code: 2004, Msg: "密码错误"}
	USER_ACCOUNT_ALREADY_EXIST = MsgCode{Code: 2008, Msg: "账号已存在"}
	CAPTCHA_ERROR              = MsgCode{Code: 2100, Msg: "验证码错误"}
	INSUFFICENT_PERMISSIONS    = MsgCode{Code: 2200, Msg: "权限不足"}

	/* 思维导图错误 3000 ~ 3999 */
	MINDMAP_NOT_FOUND        = MsgCode{Code: 3001, Msg: "思维导图不存在"}
	MINDMAP_ALREADY_EXISTS   = MsgCode{Code: 3002, Msg: "思维导图已存在"}
	MINDMAP_PERMISSION_DENIED = MsgCode{Code: 3003, Msg: "思维导图权限不足"}
)
