package types

type User struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}

type LoginReq struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
type LoginResp struct {
	User User `json:"user"`
}
