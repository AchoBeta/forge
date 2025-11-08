package repo

import (
	"context"
	"forge/biz/entity"
	"forge/biz/types"
)

type AiChatRepo interface {
	//获取某个会话
	GetConversation(ctx context.Context, conversationID, userID string) (*entity.Conversation, error)

	//获取某个导图的所有会话
	GetMapAllConversation(ctx context.Context, mapID, userID string) ([]*entity.Conversation, error)

	//保存某个会话实体
	SaveConversation(ctx context.Context, conversation *entity.Conversation) error

	//更新某个会话的聊天记录
	UpdateConversationMessage(ctx context.Context, conversation *entity.Conversation) error

	//更新某个会话的标题
	UpdateConversationTitle(ctx context.Context, conversation *entity.Conversation) error

	//删除某个会话
	DeleteConversation(ctx context.Context, conversationID, userID string) error
}

type EinoServer interface {
	//向ai发送消息
	SendMessage(ctx context.Context, messages []*entity.Message) (types.AgentResponse, error)
}
