package repo

import "forge/biz/entity"

type AiChatRepo interface {
	//获取某个会话
	GetConversation(conversationID, userID string) (*entity.Conversation, error)

	//获取某个导图的所有会话
	GetMapAllConversation(mapID, userID string) ([]*entity.Conversation, error)

	//保存某个会话实体
	SaveConversation(conversation *entity.Conversation) error

	//更新某个会话
	UpdateConversation(conversation *entity.Conversation) error

	//删除某个会话
	DeleteConversation(conversationID, userID string) error
}

type EinoServer interface {
	//向ai发送消息
	SendMessage(messages []*entity.Message) (string, error)
}
