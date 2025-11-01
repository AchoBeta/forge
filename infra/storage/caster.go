package storage

import (
	"encoding/json"
	"fmt"
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

// CastMindMapDO2PO 领域对象转持久化对象
func CastMindMapDO2PO(mindmap *entity.MindMap) *po.MindMapPO {
	if mindmap == nil {
		return nil
	}

	// 序列化Data为JSON字符串
	dataBytes, _ := json.Marshal(mindmap.Data)

	mindmapPO := &po.MindMapPO{
		MapID:  mindmap.MapID,
		UserID: mindmap.UserID,
		Title:  mindmap.Title,
		Desc:   mindmap.Desc,
		Data:   string(dataBytes),
		Layout: mindmap.Layout,
	}

	// 处理时间字段
	if !mindmap.CreatedAt.IsZero() {
		mindmapPO.CreatedAt = &mindmap.CreatedAt
	}
	if !mindmap.UpdatedAt.IsZero() {
		mindmapPO.UpdatedAt = &mindmap.UpdatedAt
	}

	return mindmapPO
}

// CastMindMapPO2DO 持久化对象转领域对象
func CastMindMapPO2DO(mindmapPO *po.MindMapPO) (*entity.MindMap, error) {
	if mindmapPO == nil {
		return nil, nil
	}

	// 反序列化JSON数据
	var data entity.MindMapData
	if err := json.Unmarshal([]byte(mindmapPO.Data), &data); err != nil {
		return nil, fmt.Errorf("unmarshal data failed: %w", err)
	}

	mindmap := &entity.MindMap{
		MapID:  mindmapPO.MapID,
		UserID: mindmapPO.UserID,
		Title:  mindmapPO.Title,
		Desc:   mindmapPO.Desc,
		Data:   data,
		Layout: mindmapPO.Layout,
	}

	// 处理时间字段
	if mindmapPO.CreatedAt != nil {
		mindmap.CreatedAt = *mindmapPO.CreatedAt
	}
	if mindmapPO.UpdatedAt != nil {
		mindmap.UpdatedAt = *mindmapPO.UpdatedAt
	}

	return mindmap, nil
}
