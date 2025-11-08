package def

import (
	"forge/biz/entity"
	"time"
)

// 请求体
type ProcessUserMessageRequest struct {
	ConversationID string `json:"conversation_id" binding:"required"`
	Content        string `json:"content" binding:"required"`
	MapData        string `json:"map_data"`
}

type ProcessUserMessageResponse struct {
	NewMapJson string `json:"new_map_json"`
	Content    string `json:"content"`
	Success    bool   `json:"success"`
}

type SaveNewConversationRequest struct {
	Title   string `json:"title" binding:"required"`
	MapID   string `json:"map_id" binding:"required"`
	MapData string `json:"map_data"`
}

type SaveNewConversationResponse struct {
	Success bool `json:"success"`
}

type GetConversationListRequest struct {
	MapID string `json:"map_id" binding:"required"`
}

type ConversationData struct {
	ConversationID string    `json:"conversation_id"`
	Title          string    `json:"title"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type GetConversationListResponse struct {
	List    []ConversationData `json:"list"`
	Success bool               `json:"success"`
}

type DelConversationRequest struct {
	ConversationID string `json:"conversation_id" binding:"required"`
}

type DelConversationResponse struct {
	Success bool `json:"success"`
}

type GetConversationRequest struct {
	ConversationID string `json:"conversation_id" binding:"required"`
}

type GetConversationResponse struct {
	Title    string            `json:"title"`
	Messages []*entity.Message `json:"messages"`
	Success  bool              `json:"success"`
}

type UpdateConversationTitleRequest struct {
	ConversationID string `json:"conversation_id" binding:"required"`
	Title          string `json:"title" binding:"required"`
}

type UpdateConversationTitleResponse struct {
	Success bool `json:"success"`
}

type GenerateMindMapRequest struct {
	Text string `json:"text"`
}

type GenerateMindMapResponse struct {
	Success bool   `json:"success"`
	MapJson string `json:"map_json"`
}
