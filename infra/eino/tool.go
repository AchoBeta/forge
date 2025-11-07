package eino

import (
	"context"
	"forge/pkg/log/zlog"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

func UpdateMindMap(ctx context.Context, params *UpdateMindMapParams) (string, error) {
	message := initToolUpdateMindMap(params.MapData, params.Requirement)

	resp, err := params.AiClient.Generate(ctx, message)
	if err != nil {
		zlog.Errorf("模型调用失败%v", err)
		return "", err
	}

	return resp.Content, nil
}

func CreateUpdateMindMapTool() tool.InvokableTool {
	updateMindMapTool := utils.NewTool(
		&schema.ToolInfo{
			Name: "update_mind_map",
			Desc: "用于修改导图,需要修改导图时调用该工具,返回完整新导图JSON",
			ParamsOneOf: schema.NewParamsOneOfByParams(
				map[string]*schema.ParameterInfo{
					"requirement": {
						Type:     schema.String,
						Desc:     "修改导图的需求，例如「把 root.children[0].data.text 改成『新产品』」",
						Required: true,
					},
				},
			),
		}, UpdateMindMap)
	return updateMindMapTool
}
