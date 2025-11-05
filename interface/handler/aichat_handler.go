package handler

import (
	"context"
	"forge/interface/caster"
	"forge/interface/def"
)

func (h *Handler) SendMessage(ctx context.Context, req *def.ProcessUserMessageRequest) (*def.ProcessUserMessageResponse, error) {

	//转 biz层 参数
	params := caster.CastProcessUserMessageReq2Params(req)

	aiMsg, err := h.AiChatService.ProcessUserMessage(ctx, params)
	if err != nil {
		return nil, err
	}

	resp := &def.ProcessUserMessageResponse{
		Content: aiMsg,
		Success: true,
	}

	return resp, nil
}

func (h *Handler) SaveNewConversation(ctx context.Context, req *def.SaveNewConversationRequest) (*def.SaveNewConversationResponse, error) {
	params := caster.CastSaveNewConversationReq2Params(req)

	err := h.AiChatService.SaveNewConversation(ctx, params)
	if err != nil {
		return nil, err
	}

	resp := &def.SaveNewConversationResponse{
		Success: true,
	}
	return resp, nil
}

func (h *Handler) GetConversationList(ctx context.Context, req *def.GetConversationListRequest) (*def.GetConversationListResponse, error) {
	params := caster.CastGetConversationListReq2Params(req)

	conversations, err := h.AiChatService.GetConversationList(ctx, params)
	if err != nil {
		return nil, err
	}

	resp := &def.GetConversationListResponse{
		Success: true,
		List:    caster.CastConversationsDOs2Resp(conversations),
	}

	return resp, nil

}

func (h *Handler) DelConversation(ctx context.Context, req *def.DelConversationRequest) (*def.DelConversationResponse, error) {
	params := caster.CastDelConversationReq2Params(req)

	err := h.AiChatService.DelConversation(ctx, params)
	if err != nil {
		return nil, err
	}

	resp := &def.DelConversationResponse{
		Success: true,
	}
	return resp, nil
}

func (h *Handler) GetConversation(ctx context.Context, req *def.GetConversationRequest) (*def.GetConversationResponse, error) {
	params := caster.CastGetConversationReq2Params(req)

	conversation, err := h.AiChatService.GetConversation(ctx, params)

	if err != nil {
		return nil, err
	}

	resp := &def.GetConversationResponse{
		Success:  true,
		Title:    conversation.Title,
		Messages: conversation.Messages,
	}

	return resp, nil
}
