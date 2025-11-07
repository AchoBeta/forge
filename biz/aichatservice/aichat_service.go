package aichatservice

import (
	"context"
	"errors"
	"forge/biz/entity"
	"forge/biz/repo"
	"forge/biz/types"
	"forge/pkg/log/zlog"
)

var (
	CONVERSATION_ID_NOT_NULL    = errors.New("会话ID不能为空")
	USER_ID_NOT_NULL            = errors.New("用户ID不能为空")
	MAP_ID_NOT_NULL             = errors.New("导图ID不能为空")
	CONVERSATION_TITLE_NOT_NULL = errors.New("会话标题不能为空")
	CONVERSATION_NOT_EXIST      = errors.New("该会话不存在")
	AI_CHAT_PERMISSION_DENIED   = errors.New("会话权限不足")
	MIND_MAP_NOT_EXIST          = errors.New("该导图不存在")
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
		return "", AI_CHAT_PERMISSION_DENIED
	}

	conversation, err := a.aiChatRepo.GetConversation(ctx, req.ConversationID, user.UserID)
	if err != nil {
		return "", err
	}

	//更新导图提示词
	conversation.ProcessSystemPrompt(req.MapData)

	//添加用户聊天记录
	conversation.AddMessage(req.Message, entity.USER)

	//调用ai 返回ai消息
	aiMsg, err := a.einoServer.SendMessage(ctx, conversation.Messages, req.MapData)
	if err != nil {
		return "", err
	}

	//添加ai消息
	conversation.AddMessage(aiMsg, entity.ASSISTANT)

	//更新会话聊天记录
	err = a.aiChatRepo.UpdateConversationMessage(ctx, conversation)
	if err != nil {
		return "", err
	}

	return aiMsg, nil
}

func (a *AiChatService) SaveNewConversation(ctx context.Context, req *types.SaveNewConversationParams) error {
	user, ok := entity.GetUser(ctx)
	if !ok {
		zlog.CtxErrorf(ctx, "未能从上下文中获取用户信息")
		return AI_CHAT_PERMISSION_DENIED
	}

	conversation, err := entity.NewConversation(user.UserID, req.MapID, req.Title)
	if err != nil {
		return err
	}
	//初始化系统提示词
	conversation.ProcessSystemPrompt(req.MapData)

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
		return nil, AI_CHAT_PERMISSION_DENIED
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
		return AI_CHAT_PERMISSION_DENIED
	}

	err := a.aiChatRepo.DeleteConversation(ctx, req.ConversationID, user.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (a *AiChatService) GetConversation(ctx context.Context, req *types.GetConversationParams) (*entity.Conversation, error) {
	user, ok := entity.GetUser(ctx)
	if !ok {
		zlog.CtxErrorf(ctx, "未能从上下文中获取用户信息")
		return nil, AI_CHAT_PERMISSION_DENIED
	}

	conversation, err := a.aiChatRepo.GetConversation(ctx, req.ConversationID, user.UserID)
	if err != nil {
		return nil, err
	}

	return conversation, nil
}

func (a *AiChatService) UpdateConversationTitle(ctx context.Context, req *types.UpdateConversationTitleParams) error {
	user, ok := entity.GetUser(ctx)
	if !ok {
		zlog.CtxErrorf(ctx, "未能从上下文中获取用户信息")
		return AI_CHAT_PERMISSION_DENIED
	}

	conversation, err := a.aiChatRepo.GetConversation(ctx, req.ConversationID, user.UserID)
	if err != nil {
		return err
	}

	conversation.UpdateTitle(req.Title)

	err = a.aiChatRepo.UpdateConversationTitle(ctx, conversation)
	if err != nil {
		return err
	}
	return nil
}
