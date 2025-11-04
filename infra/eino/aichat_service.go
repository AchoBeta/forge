package eino

import (
	"context"
	"forge/biz/entity"
	"forge/biz/repo"
	"forge/pkg/log/zlog"
	"github.com/cloudwego/eino-ext/components/model/ark"
)

type AiChatClient struct {
	ApiKey    string
	ModelName string
}

func NewAiChatClient(apiKey, modelName string) repo.EinoServer {
	return &AiChatClient{ApiKey: apiKey, ModelName: modelName}
}

func (c *AiChatClient) SendMessage(ctx context.Context, messages []*entity.Message) (string, error) {

	aiChatModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: c.ApiKey,
		Model:  c.ModelName,
	})
	if aiChatModel == nil || err != nil {
		zlog.Errorf("ai模型连接失败：%v", err)
		return "", err
	}

	input := messagesDo2Input(messages)

	resp, err := aiChatModel.Generate(ctx, input)
	if err != nil {
		zlog.Errorf("模型调用失败%v", err)
		return "", err
	}
	return resp.Content, nil
}
