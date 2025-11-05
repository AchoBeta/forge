package eino

import (
	"context"
	"fmt"
	"forge/biz/entity"
	"forge/biz/repo"
	"forge/pkg/log/zlog"
	"github.com/cloudwego/eino-ext/components/model/ark"
)

type AiChatClient struct {
	ApiKey    string
	ModelName string
	Client    *ark.ChatModel
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

	return &AiChatClient{ApiKey: apiKey, ModelName: modelName}
}

func (c *AiChatClient) SendMessage(ctx context.Context, messages []*entity.Message) (string, error) {

	input := messagesDo2Input(messages)

	resp, err := c.Client.Generate(ctx, input)
	if err != nil {
		zlog.Errorf("模型调用失败%v", err)
		return "", err
	}
	return resp.Content, nil
}
