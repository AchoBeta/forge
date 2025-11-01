package storage

import (
	"forge/biz/entity"
	"forge/infra/storage/po"
)

// CastUserDO2PO
//
//	@Description: 实体与存储互转
//	@param user
//	@return *po.UserPO
func CastUserDO2PO(user *entity.User) *po.UserPO {
	if user == nil {
		return nil
	}
	return &po.UserPO{
		UserID:        user.UserID,
		UserName:      user.UserName,
		Avatar:        user.Avatar,
		Password:      user.Password,
		Phone:         user.Phone,
		Email:         user.Email,
		Status:        user.Status,
		PhoneVerified: user.PhoneVerified,
		EmailVerified: user.EmailVerified,
		LastLoginAt:   user.LastLoginAt,
	}
}

// CastUserPO2DO 存储转实体
func CastUserPO2DO(userPO *po.UserPO) *entity.User {
	if userPO == nil {
		return nil
	}
	user := &entity.User{
		UserID:        userPO.UserID,
		UserName:      userPO.UserName,
		Avatar:        userPO.Avatar,
		Password:      userPO.Password,
		Phone:         userPO.Phone,
		Email:         userPO.Email,
		Status:        userPO.Status,
		PhoneVerified: userPO.PhoneVerified,
		EmailVerified: userPO.EmailVerified,
		LastLoginAt:   userPO.LastLoginAt,
	}

	// 处理时间字段：如果 PO 中为 nil，Entity 中保持零值；否则解引用
	if userPO.CreatedAt != nil {
		user.CreatedAt = *userPO.CreatedAt
	}
	if userPO.UpdatedAt != nil {
		user.UpdatedAt = *userPO.UpdatedAt
	}

	return user
}
