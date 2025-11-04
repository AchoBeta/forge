package aichatservice

import (
	"context"
	"errors"
	"forge/biz/entity"
	"forge/biz/repo"
	"forge/biz/types"
	"forge/pkg/log/zlog"
)

type AiChatService struct {
	aiChatRepo repo.AiChatRepo
	einoServer repo.EinoServer
}

func NewAiChatService(aiChatRepo repo.AiChatRepo, einoServer repo.EinoServer) *AiChatService {
	return &AiChatService{aiChatRepo: aiChatRepo, einoServer: einoServer}
}

func (a *AiChatService) ProcessUserMessage(ctx context.Context, req *types.ProcessUserMessageParams) (string, error) {
	user, ok := entity.GetUser(ctx)
	if !ok {
		zlog.CtxErrorf(ctx, "未能从上下文中获取用户信息")
		return "", errors.New("会话权限不足")
	}

	conversation, err := a.aiChatRepo.GetConversation(ctx, req.ConversationID, user.UserID)
	if err != nil {
		return "", err
	}

	//添加用户聊天记录
	conversation.AddMessage(req.Message, "user")

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

func (a *AiChatService) SaveNewConversation(ctx context.Context, req *types.SaveNewConversationParams) error {
	user, ok := entity.GetUser(ctx)
	if !ok {
		zlog.CtxErrorf(ctx, "未能从上下文中获取用户信息")
		return errors.New("会话权限不足")
	}

	conversation, err := entity.NewConversation(user.UserID, req.MapID, req.Title)
	if err != nil {
		return err
	}

	err = a.aiChatRepo.SaveConversation(ctx, conversation)
	if err != nil {
		return err
	}
	return nil
}

func (a *AiChatService) GetConversationList(ctx context.Context, req *types.GetConversationListParams) ([]*entity.Conversation, error) {
	user, ok := entity.GetUser(ctx)
	if !ok {
		zlog.CtxErrorf(ctx, "未能从上下文中获取用户信息")
		return nil, errors.New("会话权限不足")
	}

	conversationList, err := a.aiChatRepo.GetMapAllConversation(ctx, req.MapID, user.UserID)
	if err != nil {
		return nil, err
	}

	return conversationList, nil
}

func (a *AiChatService) DelConversation(ctx context.Context, req *types.DelConversationParams) error {
	user, ok := entity.GetUser(ctx)
	if !ok {
		zlog.CtxErrorf(ctx, "未能从上下文中获取用户信息")
		return errors.New("会话权限不足")
	}

	err := a.aiChatRepo.DeleteConversation(ctx, req.ConversationID, user.UserID)
	if err != nil {
		return err
	}

	return nil
}
