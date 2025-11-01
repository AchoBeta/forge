package storage

import (
	"context"
	"fmt"
	"forge/biz/entity"
	"forge/biz/repo"
	"forge/infra/database"
	"forge/infra/storage/po"

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

// CreateUser 创建用户
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

// UpdateUser 更新用户信息 - 统一的更新接口
func (u *userPersistence) UpdateUser(ctx context.Context, updateInfo *repo.UserUpdateInfo) error {
	if updateInfo == nil || updateInfo.UserID == "" {
		return fmt.Errorf("invalid update info: userID is required") // 需要id定位用户
	}

	updates := map[string]any{}

	// 基础信息
	if updateInfo.UserName != nil {
		updates["username"] = *updateInfo.UserName
	}
	if updateInfo.Avatar != nil {
		updates["avatar"] = *updateInfo.Avatar
	}

	// 联系方式
	if updateInfo.Phone != nil {
		updates["phone"] = *updateInfo.Phone
	}
	if updateInfo.Email != nil {
		updates["email"] = *updateInfo.Email
	}

	// 密码
	if updateInfo.Password != nil {
		updates["password"] = *updateInfo.Password
	}

	// 状态信息
	if updateInfo.Status != nil {
		updates["status"] = *updateInfo.Status
	}
	if updateInfo.PhoneVerified != nil {
		updates["phone_verified"] = *updateInfo.PhoneVerified
	}
	if updateInfo.EmailVerified != nil {
		updates["email_verified"] = *updateInfo.EmailVerified
	}

	// 时间信息
	if updateInfo.LastLoginAt != nil {
		updates["last_login_at"] = *updateInfo.LastLoginAt
	}

	if len(updates) == 0 {
		return nil
	}

	return u.db.WithContext(ctx).Model(&po.UserPO{}).Where("user_id = ?", updateInfo.UserID).Updates(updates).Error
}

// GetUser 用户查询接口，根据查询条件获取用户
func (u *userPersistence) GetUser(ctx context.Context, query repo.UserQuery) (*entity.User, error) {
	var userPO po.UserPO
	var err error

	// 根据查询条件构建查询
	db := u.db.WithContext(ctx)

	switch {
	case query.UserID != "":
		err = db.Where("user_id = ?", query.UserID).First(&userPO).Error
	case query.UserName != "":
		err = db.Where("username = ?", query.UserName).First(&userPO).Error
	case query.Phone != "":
		err = db.Where("phone = ?", query.Phone).First(&userPO).Error
	case query.Email != "":
		err = db.Where("email = ?", query.Email).First(&userPO).Error
		/// case query.Platform != "" && query.ThirdID != "":
		// 第三方登录查询逻辑
		return nil, nil
	default:
		return nil, fmt.Errorf("invalid user query: no query field provided")
	}

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return CastUserPO2DO(&userPO), nil
}
