package repo

import (
	"context"
	"forge/biz/entity"
	"time"
)

// 得益于repo的概念，service的代码只需要调用该方法即可，
// 不用考虑具体实现
// repo应该做到尽量一个接口就能解决一个问题，不要讲接口拆的很细，
// 如果接口很细的话也可以自己聚合一下
// 屏蔽掉什么锁，什么事务概念，这不是领域考虑的事

// 用户服务，围绕用户实体的存储展开
type UserRepo interface {
	// 只读操作
	// 创建用户
	CreateUser(ctx context.Context, user *entity.User) error

	//更新用户  统一更新接口 包括密码
	UpdateUser(ctx context.Context, updateInfo *UserUpdateInfo) error

	// 只读操作
	// GetUser 根据查询条件获取用户，支持多种查询方式
	GetUser(ctx context.Context, query UserQuery) (*entity.User, error)

	/*  根据第三方登录方式查询 后续可能有更多第三方登录方式
	GetByThirdParty(ctx context.Context, platform string, id string) (*entity.User, error)
	*/

	/*
	   // 换绑手机号/邮箱
	   BindPhone(ctx context.Context, userID string, newPhone string) error
	   UnbindPhone(ctx context.Context, userID string) error
	   BindEmail(ctx context.Context, userID string, newEmail string) error
	   UnbindEmail(ctx context.Context, userID string) error

	   // 第三方绑定/解绑
	   BindWechat(ctx context.Context, userID string, openID string, unionID string) error
	   UnbindWechat(ctx context.Context, userID string) error
	   BindGithub(ctx context.Context, userID string, githubID string, githubLogin string) error
	   UnbindGithub(ctx context.Context, userID string) error
	*/
}

// 计数器服务，用于对指定key原子递增加上某个特定的值
type CounterRepo interface {
	Incr(ctx context.Context, key string, delta int) (int, error)
	Count(ctx context.Context, key string) (int, error)
}

// UserUpdateInfo 用户更新信息
type UserUpdateInfo struct {
	UserID string // 用户ID

	// 基础信息
	Name   *string // 用户名
	Avatar *string // 头像URL

	// 联系方式
	Phone *string // 手机号
	Email *string // 邮箱

	Password *string // 密码

	// 状态信息
	Status        *int  // 用户状态 1:正常 0:禁用
	PhoneVerified *bool // 手机号是否已验证
	EmailVerified *bool // 邮箱是否已验证

	// 时间信息
	LastLoginAt *time.Time // 最后登录时间

	// 第三方登录（暂不开放，后续扩展）
	/*
	   WechatOpenID  *string
	   WechatUnionID *string
	   GithubID      *string
	   GithubLogin   *string
	*/
}

// UserQuery 用户查询条件
type UserQuery struct {
	UserID string // 根据用户ID查询
	Name   string // 根据用户名查询
	Phone  string // 根据手机号查询
	Email  string // 根据邮箱查询
	// Platform string // 第三方平台
	// ThirdID  string // 第三方ID
}

// NewUserQueryByID 创建根据用户ID查询的条件
func NewUserQueryByID(userID string) UserQuery {
	return UserQuery{UserID: userID}
}

// NewUserQueryByName 创建根据用户名查询的条件
func NewUserQueryByName(name string) UserQuery {
	return UserQuery{Name: name}
}

// NewUserQueryByPhone 创建根据手机号查询的条件
func NewUserQueryByPhone(phone string) UserQuery {
	return UserQuery{Phone: phone}
}

// NewUserQueryByEmail 创建根据邮箱查询的条件
func NewUserQueryByEmail(email string) UserQuery {
	return UserQuery{Email: email}
}
