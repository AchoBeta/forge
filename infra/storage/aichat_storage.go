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

type aiChatPersistence struct {
	db *gorm.DB
}

var cp *aiChatPersistence

func InitAiChatStorage() {
	db := database.ForgeDB()

	if err := db.AutoMigrate(&po.ConversationPO{}); err != nil {
		panic(fmt.Sprintf("自动建表失败 :%v", err))
	}

	cp = &aiChatPersistence{db: db}
}

func GetAiChatPersistence() repo.AiChatRepo { return cp }

func (a aiChatPersistence) GetConversation(ctx context.Context, conversationID, userID string) (*entity.Conversation, error) {
	if conversationID == "" || userID == "" {
		return nil, fmt.Errorf("会话ID和用户ID不能为空")
	}

	var conversationPO po.ConversationPO
	if err := a.db.WithContext(ctx).Model(&po.ConversationPO{}).Where("conversation_id = ? AND user_id = ?", conversationID, userID).First(&conversationPO).Error; err != nil {
		return nil, fmt.Errorf("该会话不存在")
	}

	return CastConversationPO2DO(&conversationPO)
}

func (a aiChatPersistence) GetMapAllConversation(ctx context.Context, mapID, userID string) ([]*entity.Conversation, error) {
	//TODO implement me
	panic("implement me")
}

func (a aiChatPersistence) SaveConversation(ctx context.Context, conversation *entity.Conversation) error {
	//TODO implement me
	panic("implement me")
}

func (a aiChatPersistence) UpdateConversation(ctx context.Context, conversation *entity.Conversation) error {
	//TODO implement me
	panic("implement me")
}

func (a aiChatPersistence) DeleteConversation(ctx context.Context, conversationID, userID string) error {
	//TODO implement me
	panic("implement me")
}
