package eino

import (
	"context"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

func (a *AiChatClient) UpdateMindMap(ctx context.Context, params *UpdateMindMapParams) (string, error) {
	message := initToolUpdateMindMap(params.MapJson, params.Requirement)

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
					"map_json": {
						Type:     schema.String,
						Desc:     "传入最开始的系统消息的导图json数据,不要注释、不要 Markdown 包裹",
						Required: true,
					},
				},
			),
		}, a.UpdateMindMap)
	return updateMindMapTool
}
