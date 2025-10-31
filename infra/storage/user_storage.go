package storage

import (
	"context"
	"forge/biz/entity"
	"forge/biz/repo"
	"forge/infra/database"

	"gorm.io/gorm"
)

type userPersistence struct {
	db *gorm.DB
}

var up *userPersistence

func InitUserStorage() {
	up = &userPersistence{
		db: database.ForgeDB(),
	}
}

func GetUserPersistence() repo.UserRepo {
	return up
}

func (u *userPersistence) CreateUser(ctx context.Context, user *entity.User) error {
	userPO := CastUserDO2PO(user)
	err := u.db.WithContext(ctx).Create(&userPO).Error
	if err != nil {
		//todo 这里如何让上游更好地感知到错误类型，甚至前端感知到错误类型呢？
		return err
	}
	return nil
}

// 其他仓储
