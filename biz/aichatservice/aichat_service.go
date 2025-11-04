package aichatservice

import (
	"context"
	"forge/biz/repo"
)

type AiChatService struct {
	aiChatRepo repo.AiChatRepo
	einoServer repo.EinoServer
}

func NewAiChatService(aiChatRepo repo.AiChatRepo, einoServer repo.EinoServer) *AiChatService {
	return &AiChatService{aiChatRepo: aiChatRepo}
}

func (a *AiChatService) ProcessUserMessage(ctx context.Context, conversationID, userID, message string) (string, error) {
	conversation, err := a.aiChatRepo.GetConversation(ctx, conversationID, userID)
	if err != nil {
		return "", err
	}

	//添加用户聊天记录
	conversation.AddMessage(message, "user")

	//调用ai 返回ai消息
	aiMsg, err := a.einoServer.SendMessage(ctx, conversation.Messages)
	if err != nil {
		return "", err
	}

	//添加ai消息
	conversation.AddMessage(aiMsg, "assistant")

	//更新会话
	err = a.aiChatRepo.UpdateConversation(ctx, conversation)
	if err != nil {
		return "", err
	}

	return aiMsg, nil
}
