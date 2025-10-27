package def

// 这个是DTO层，会暴露给前端 主要是接口定义

type User struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Dogs     []*Dog `json:"dogs"`
}

type Dog struct {
	DogID   string `json:"dog_id"`
	DogName string `json:"dog_name"`
}

type LoginReq struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
type LoginResp struct {
	User *User `json:"user"`
}
