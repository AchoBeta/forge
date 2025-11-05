package caster

import (
	"forge/biz/entity"
	"forge/biz/types"
	"forge/interface/def"
)

func CastProcessUserMessageReq2Params(req *def.ProcessUserMessageRequest) *types.ProcessUserMessageParams {
	if req == nil {
		return nil
	}

	return &types.ProcessUserMessageParams{
		ConversationID: req.ConversationID,
		Message:        req.Content,
	}
}

func CastSaveNewConversationReq2Params(req *def.SaveNewConversationRequest) *types.SaveNewConversationParams {
	if req == nil {
		return nil
	}
	return &types.SaveNewConversationParams{
		Title: req.Title,
		MapID: req.MapID,
	}
}

func CastGetConversationListReq2Params(req *def.GetConversationListRequest) *types.GetConversationListParams {
	if req == nil {
		return nil
	}
	return &types.GetConversationListParams{
		MapID: req.MapID,
	}
}

func CastConversationsDOs2Resp(conversations []*entity.Conversation) []def.ConversationData {
	if conversations == nil {
		return nil
	}

	conversationsData := make([]def.ConversationData, len(conversations))

	for i, conversation := range conversations {
		conversationsData[i] = def.ConversationData{
			ConversationID: conversation.ConversationID,
			Title:          conversation.Title,
			CreatedAt:      conversation.CreatedAt,
			UpdatedAt:      conversation.UpdatedAt,
		}
	}

	return conversationsData
}

func CastDelConversationReq2Params(req *def.DelConversationRequest) *types.DelConversationParams {
	if req == nil {
		return nil
	}

	return &types.DelConversationParams{
		ConversationID: req.ConversationID,
	}
}

func CastGetConversationReq2Params(req *def.GetConversationRequest) *types.GetConversationParams {
	if req == nil {
		return nil
	}
	return &types.GetConversationParams{
		ConversationID: req.ConversationID,
	}
}

func CastUpdateConversationTitleReq2Params(req *def.UpdateConversationTitleRequest) *types.UpdateConversationTitleParams {
	if req == nil {
		return nil
	}
	return &types.UpdateConversationTitleParams{
		Title:          req.Title,
		ConversationID: req.ConversationID,
	}
}
