package types

import (
	"context"
	"forge/biz/entity"
	"mime/multipart"

	"github.com/cloudwego/eino/schema"
)

type IAiChatService interface {
	//处理用户消息
	ProcessUserMessage(ctx context.Context, req *ProcessUserMessageParams) (AgentResponse, error)

	//保存新的会话
	SaveNewConversation(ctx context.Context, req *SaveNewConversationParams) (string, error)

	//获取该导图的所有会话
	GetConversationList(ctx context.Context, req *GetConversationListParams) ([]*entity.Conversation, error)

	//删除某会话
	DelConversation(ctx context.Context, req *DelConversationParams) error

	//获取某会话的详细信息
	GetConversation(ctx context.Context, req *GetConversationParams) (*entity.Conversation, error)

	//更新某会话的标题
	UpdateConversationTitle(ctx context.Context, req *UpdateConversationTitleParams) error

	//生成导图
	GenerateMindMap(ctx context.Context, req *GenerateMindMapParams) (string, error)

	//批量生成导图（Pro版本）
	GenerateMindMapPro(ctx context.Context, req *GenerateMindMapProParams) (*entity.GenerationBatch, []*entity.GenerationResult, []*entity.Conversation, error)
}

type ProcessUserMessageParams struct {
	ConversationID string
	Message        string
	MapData        string
}

type SaveNewConversationParams struct {
	Title   string
	MapID   string
	MapData string
}

type GetConversationListParams struct {
	MapID string
}

type DelConversationParams struct {
	ConversationID string
}

type GetConversationParams struct {
	ConversationID string
}

type UpdateConversationTitleParams struct {
	ConversationID string
	Title          string
}

type AgentResponse struct {
	NewMapJson string            `json:"new_map_json"`
	Content    string            `json:"content"`
	ToolCallID string            `json:"tool_call_id"`
	ToolCalls  []schema.ToolCall `json:"tool_calls"`
}

type GenerateMindMapParams struct {
	Text string
	File *multipart.FileHeader
}

// GenerationResultWithParams 带生成参数的结果
type GenerationResultWithParams struct {
	JSON        string   // 生成的JSON
	Temperature *float64 // 温度参数
	TopP        *float64 // Top-P参数
	Strategy    int      // 生成策略
	Error       error    // 生成错误
}
