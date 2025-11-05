package storage

import (
	"context"
	"errors"
	"fmt"
	"forge/biz/aichatservice"
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
		return nil, aichatservice.CONVERSATION_ID_NOT_NULL
	} else if userID == "" {
		return nil, aichatservice.USER_ID_NOT_NULL
	}

	var conversationPO po.ConversationPO
	if err := a.db.WithContext(ctx).Model(&po.ConversationPO{}).Where("conversation_id = ? AND user_id = ?", conversationID, userID).First(&conversationPO).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, aichatservice.CONVERSATION_NOT_EXIST
		}
		return nil, fmt.Errorf("数据库出错 :%v", err)

	}

	return CastConversationPO2DO(&conversationPO)
}

func (a *aiChatPersistence) GetMapAllConversation(ctx context.Context, mapID, userID string) ([]*entity.Conversation, error) {
	if mapID == "" {
		return nil, aichatservice.MAP_ID_NOT_NULL
	} else if userID == "" {
		return nil, aichatservice.USER_ID_NOT_NULL
	} else if !checkMapIsExist(ctx, a, mapID) {
		return nil, aichatservice.MIND_MAP_NOT_EXIST
	}

	var conversationPOs []po.ConversationPO
	if err := a.db.WithContext(ctx).Model(&po.ConversationPO{}).Where("map_id = ? AND user_id = ?", mapID, userID).Find(&conversationPOs).Error; err != nil {
		return nil, fmt.Errorf("获取导图会话时 数据库出错 %w", err)
	}

	return CastConversationPOs2DOs(conversationPOs)
}

func (a *aiChatPersistence) SaveConversation(ctx context.Context, conversation *entity.Conversation) error {
	if conversation.ConversationID == "" {
		return aichatservice.CONVERSATION_ID_NOT_NULL
	} else if conversation.UserID == "" {
		return aichatservice.USER_ID_NOT_NULL
	} else if conversation.MapID == "" {
		return aichatservice.MAP_ID_NOT_NULL
	} else if conversation.Title == "" {
		return aichatservice.CONVERSATION_TITLE_NOT_NULL
	} else if !checkMapIsExist(ctx, a, conversation.MapID) {
		return aichatservice.MIND_MAP_NOT_EXIST
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

func (a *aiChatPersistence) UpdateConversationMessage(ctx context.Context, conversation *entity.Conversation) error {
	if conversation.UserID == "" {
		return aichatservice.USER_ID_NOT_NULL
	} else if conversation.MapID == "" {
		return aichatservice.MAP_ID_NOT_NULL
	} else if !checkConversationIsExist(ctx, a, conversation.ConversationID) {
		return aichatservice.CONVERSATION_NOT_EXIST
	}

	conversationPO, err := CastConversationDO2PO(conversation)
	if err != nil {
		return err
	}

	Updates := make(map[string]interface{})
	if conversationPO.Messages != nil {
		Updates["messages"] = conversationPO.Messages
	}

	err = a.db.WithContext(ctx).Model(&po.ConversationPO{}).Where("conversation_id = ? AND user_id = ?", conversationPO.ConversationID, conversationPO.UserID).Updates(Updates).Error
	if err != nil {
		return fmt.Errorf("更新会话时 数据库出错 %w", err)
	}
	return nil
}

func (a *aiChatPersistence) UpdateConversationTitle(ctx context.Context, conversation *entity.Conversation) error {
	if conversation.UserID == "" {
		return aichatservice.USER_ID_NOT_NULL
	} else if conversation.MapID == "" {
		return aichatservice.MAP_ID_NOT_NULL
	} else if !checkConversationIsExist(ctx, a, conversation.ConversationID) {
		return aichatservice.CONVERSATION_NOT_EXIST
	}

	conversationPO, err := CastConversationDO2PO(conversation)
	if err != nil {
		return err
	}
	Updates := make(map[string]interface{})
	if conversationPO.Messages != nil {
		Updates["title"] = conversationPO.Title
	}

	err = a.db.WithContext(ctx).Model(&po.ConversationPO{}).Where("conversation_id = ? AND user_id = ?", conversationPO.ConversationID, conversationPO.UserID).Updates(Updates).Error
	if err != nil {
		return fmt.Errorf("更新会话时 数据库出错 %w", err)
	}
	return nil
}

func (a *aiChatPersistence) DeleteConversation(ctx context.Context, conversationID, userID string) error {
	if conversationID == "" {
		return aichatservice.CONVERSATION_ID_NOT_NULL
	} else if userID == "" {
		return aichatservice.USER_ID_NOT_NULL
	}

	result := a.db.WithContext(ctx).Model(&po.ConversationPO{}).Where("conversation_id = ? AND user_id = ?", conversationID, userID).Delete(&po.ConversationPO{})
	if result.RowsAffected == 0 {
		return aichatservice.CONVERSATION_NOT_EXIST
	}
	if result.Error != nil {
		return fmt.Errorf("删除会话时出错 %w", result.Error)
	}
	return nil
}

func checkMapIsExist(ctx context.Context, a *aiChatPersistence, checkMapID string) bool {
	var check bool
	a.db.WithContext(ctx).Model(&po.MindMapPO{}).Raw("SELECT EXISTS(SELECT 1 FROM achobeta_forge_mindmap WHERE map_id = ?)", checkMapID).Scan(&check)
	return check
}

func checkConversationIsExist(ctx context.Context, a *aiChatPersistence, checkConversationID string) bool {
	var check bool
	a.db.WithContext(ctx).Model(&po.ConversationPO{}).Raw("SELECT EXISTS(SELECT 1 FROM achobeta_forge_conversation WHERE conversation_id = ?)", checkConversationID).Scan(&check)
	return check
}
