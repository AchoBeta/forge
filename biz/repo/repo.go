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
	CreateUser(ctx context.Context, user *entity.User) error
}

// 计数器服务，用于对指定key原子递增加上某个特定的值
type CounterRepo interface {
	Incr(ctx context.Context, key string, delta int) (int, error)
	Count(ctx context.Context, key string) (int, error)
}
