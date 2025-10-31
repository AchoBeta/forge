package repo

import (
	"context"
	"forge/biz/entity"
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

	// 更新用户
	UpdateUser(ctx context.Context, updateInfo *UserUpdateInfo) error

	// 更改密码
	UpdatePassWord(ctx context.Context, userID string, NewPassword string) error

	// 只读操作
	// 根据id查询
	GetByUserID(ctx context.Context, userID string) (*entity.User, error)
	// 根据用户名查询
	GetByName(ctx context.Context, name string) (*entity.User, error)
	// 根据手机号查询
	GetByPhone(ctx context.Context, phone string) (*entity.User, error)
	// 根据邮箱查询
	GetByEmail(ctx context.Context, email string) (*entity.User, error)

	/*  根据第三方登录方式查询 后续可能有更多第三方登录方式
	FindByThirdParty(ctx context.Context, platform string, id string) (*entity.User, error)
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

type UserUpdateInfo struct {
	//后续可能能换头像，名称等
}
