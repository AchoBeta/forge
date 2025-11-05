package def

import (
	"forge/biz/entity"
	"time"
)

// 请求体
type ProcessUserMessageRequest struct {
	ConversationID string `json:"conversation_id" binding:"required"`
	Content        string `json:"content" binding:"required"`
}

type ProcessUserMessageResponse struct {
	Content string `json:"content"`
	Success bool   `json:"success"`
}

type SaveNewConversationRequest struct {
	Title string `json:"title" binding:"required"`
	MapID string `json:"map_id" binding:"required"`
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
