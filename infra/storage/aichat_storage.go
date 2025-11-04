package storage

import (
	"context"
	"errors"
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

func (a *aiChatPersistence) GetConversation(ctx context.Context, conversationID, userID string) (*entity.Conversation, error) {
	if conversationID == "" {
		return nil, fmt.Errorf("会话ID不能为空")
	} else if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	var conversationPO po.ConversationPO
	if err := a.db.WithContext(ctx).Model(&po.ConversationPO{}).Where("conversation_id = ? AND user_id = ?", conversationID, userID).First(&conversationPO).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("该会话不存在")
		}
		return nil, fmt.Errorf("数据库出错 :%v", err)

	}

	return CastConversationPO2DO(&conversationPO)
}

func (a *aiChatPersistence) GetMapAllConversation(ctx context.Context, mapID, userID string) ([]*entity.Conversation, error) {
	if mapID == "" {
		return nil, fmt.Errorf("导图ID不能为空")
	} else if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}
	var conversationPOs []po.ConversationPO
	if err := a.db.WithContext(ctx).Model(&po.ConversationPO{}).Where("map_id = ? AND user_id = ?", mapID, userID).Find(&conversationPOs).Error; err != nil {
		return nil, fmt.Errorf("获取导图会话时 数据库出错 %w", err)
	}

	return CastConversationPOs2DOs(conversationPOs)
}

func (a *aiChatPersistence) SaveConversation(ctx context.Context, conversation *entity.Conversation) error {
	if conversation.ConversationID == "" || conversation.Title == "" || conversation.MapID == "" {
		return fmt.Errorf("会话ID不能为空")
	} else if conversation.UserID == "" {
		return fmt.Errorf("用户ID不能为空")
	} else if conversation.MapID == "" {
		return fmt.Errorf("导图ID不能为空")
	} else if conversation.Title == "" {
		return fmt.Errorf("会话标题不能为空")
	}

	conversationPO, err := CastConversationDO2PO(conversation)
	if err != nil {
		return err
	}
	err = a.db.WithContext(ctx).Model(&po.ConversationPO{}).Create(&conversationPO).Error
	if err != nil {
		return fmt.Errorf("保存会话时，数据库出错 %w", err)
	}
	return nil
}

func (a *aiChatPersistence) UpdateConversation(ctx context.Context, conversation *entity.Conversation) error {
	if conversation.UserID == "" {
		return fmt.Errorf("用户ID不能为空")
	} else if conversation.MapID == "" {
		return fmt.Errorf("导图ID不能为空")
	}

	conversationPO, err := CastConversationDO2PO(conversation)
	if err != nil {
		return err
	}

	Updates := make(map[string]interface{})
	if conversationPO.Title != "" {
		Updates["title"] = conversation.Title
	}
	if conversationPO.Messages != nil {
		Updates["messages"] = conversationPO.Messages
	}

	err = a.db.WithContext(ctx).Model(&po.ConversationPO{}).Where("conversation_id = ? AND user_id = ?", conversationPO.ConversationID, conversationPO.UserID).Updates(Updates).Error
	if err != nil {
		return fmt.Errorf("更新会话时 数据库出错 %w", err)
	}
	return nil
}

func (a *aiChatPersistence) DeleteConversation(ctx context.Context, conversationID, userID string) error {
	if conversationID == "" {
		return fmt.Errorf("会话ID不能为空")
	} else if userID == "" {
		return fmt.Errorf("用户ID不能为空")
	}

	result := a.db.WithContext(ctx).Model(&po.ConversationPO{}).Where("conversation_id = ? AND user_id = ?", conversationID, userID).Delete(&po.ConversationPO{})
	if result.RowsAffected == 0 {
		return fmt.Errorf("该会话不存在")
	}
	if result.Error != nil {
		return fmt.Errorf("删除会话时出错 %w", result.Error)
	}
	return nil
}
