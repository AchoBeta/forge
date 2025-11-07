package eino

import (
	"context"
	"fmt"
	"forge/biz/entity"
	"forge/biz/repo"
	"forge/pkg/log/zlog"
	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/schema"
)

type AiChatClient struct {
	ApiKey     string
	ModelName  string
	Client     *ark.ChatModel
	ToolClient *ark.ChatModel
}

func NewAiChatClient(apiKey, modelName string) repo.EinoServer {
	ctx := context.Background()
	aiChatModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: apiKey,
		Model:  modelName,
	})
	if aiChatModel == nil || err != nil {
		zlog.Errorf("ai模型连接失败: %v", err)
		panic(fmt.Errorf("ai模型连接失败: %v", err))
	}
	updateMindMapTool := CreateUpdateMindMapTool()
	infoTool, err := updateMindMapTool.Info(ctx)
	if err != nil {
		zlog.Errorf("ai绑定工具失败: %v", err)
		panic(fmt.Errorf("ai绑定工具失败: %v", err))
	}

	infosTool := []*schema.ToolInfo{
		infoTool,
	}
	err = aiChatModel.BindTools(infosTool)
	if err != nil {
		zlog.Errorf("ai绑定工具失败: %v", err)
		panic(fmt.Errorf("ai绑定工具失败: %v", err))
	}

	toolModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: apiKey,
		Model:  modelName,
	})
	if toolModel == nil || err != nil {
		zlog.Errorf("ToolAi模型连接失败: %v", err)
		panic(fmt.Errorf("ToolAi模型连接失败: %v", err))
	}

	return &AiChatClient{ApiKey: apiKey, ModelName: modelName, Client: aiChatModel, ToolClient: toolModel}
}

func (a *AiChatClient) SendMessage(ctx context.Context, messages []*entity.Message, mapData string) (string, error) {

	input := messagesDo2Input(messages)

	resp, err := a.Client.Generate(ctx, input)
	if err != nil {
		zlog.Errorf("模型调用失败%v", err)
		return "", err
	}
	return resp.Content, nil
}
