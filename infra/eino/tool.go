package eino

import (
	"context"
	"fmt"
	"forge/biz/entity"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

func (a *AiChatClient) UpdateMindMap(ctx context.Context, params *UpdateMindMapParams) (string, error) {
	conversation, ok := entity.GetConversation(ctx)
	if !ok {
		return "", fmt.Errorf("未能从上下文中获取到导图数据")
	}
	//fmt.Println(conversation.MapData)
	message := initToolUpdateMindMap(conversation.MapData, params.Requirement)

	resp, err := a.ToolAiClient.Generate(ctx, message)
	if err != nil {
		return "", err
	}
	return resp.Content, nil
}

func (a *AiChatClient) CreateUpdateMindMapTool() tool.InvokableTool {
	updateMindMapTool := utils.NewTool(
		&schema.ToolInfo{
			Name: "update_mind_map",
			Desc: "用于修改导图,需要修改导图时调用该工具,返回完整新导图JSON",
			ParamsOneOf: schema.NewParamsOneOfByParams(
				map[string]*schema.ParameterInfo{
					"requirement": {
						Type:     schema.String,
						Desc:     "需要工具修改导图的需求，例如「把 root.children[0].data.text 改成『新产品』」",
						Required: true,
					},
				},
			),
		}, a.UpdateMindMap)
	return updateMindMapTool
}
