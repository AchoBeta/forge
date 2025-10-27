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
	return &po.UserPO{
		UserID:   user.UserID,
		UserName: user.Name,
		Password: user.Password,
	}
}
