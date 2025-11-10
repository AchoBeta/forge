package eino

import (
	"context"
	"errors"
	"fmt"
	"forge/biz/entity"
	"forge/biz/repo"
	"forge/biz/types"
	"forge/pkg/log/zlog"
	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
)

type AiChatClient struct {
	ApiKey       string
	ModelName    string
	Agent        compose.Runnable[[]*schema.Message, types.AgentResponse]
	ToolAiClient *ark.ChatModel
}

type State struct {
	Content   string
	ToolCalls []schema.ToolCall
}

func initState(ctx context.Context) *State {
	return &State{
		Content: "",
	}
}

func NewAiChatClient(apiKey, modelName string) repo.EinoServer {
	ctx := context.Background()

	var aiChatClient AiChatClient

	//初始化工具专用模型
	toolModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey:   apiKey,
		Model:    modelName,
		Thinking: &model.Thinking{Type: model.ThinkingTypeDisabled},
	})
	if toolModel == nil || err != nil {
		zlog.Errorf("ToolAi模型连接失败: %v", err)
		panic(fmt.Errorf("ToolAi模型连接失败: %v", err))
	}

	//toolAiClient = toolModel

	aiChatClient.ApiKey = apiKey
	aiChatClient.ModelName = modelName
	aiChatClient.ToolAiClient = toolModel

	//构建agent
	aiChatModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey:   apiKey,
		Model:    modelName,
		Thinking: &model.Thinking{Type: model.ThinkingTypeDisabled},
	})
	if aiChatModel == nil || err != nil {
		zlog.Errorf("ai模型连接失败: %v", err)
		panic(fmt.Errorf("ai模型连接失败: %v", err))
	}
	updateMindMapTool := aiChatClient.CreateUpdateMindMapTool()
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

	ToolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
		Tools: []tool.BaseTool{
			updateMindMapTool,
		},
	})

	if err != nil {
		zlog.Errorf("创建工具节点失败: %v", err)
		panic("创建工具节点失败," + err.Error())
	}

	//分支中的lambda 用于对其前后输入输出
	lambda1 := compose.InvokableLambda(func(ctx context.Context, input *schema.Message) (output []*schema.Message, err error) {
		output = make([]*schema.Message, 0)
		output = append(output, input)
		return output, nil
	})

	//分支结束统一进入的lambda 用于处理输出的数据
	lambda2 := compose.InvokableLambda(func(ctx context.Context, input []*schema.Message) (output types.AgentResponse, err error) {
		//fmt.Println("lambda测试：", input)

		if len(input) == 0 {
			return types.AgentResponse{}, errors.New("agent出错")
		}

		output = types.AgentResponse{}

		if input[len(input)-1].Role == schema.Tool {
			output.NewMapJson = input[len(input)-1].Content
			output.ToolCallID = input[len(input)-1].ToolCallID

		}
		_ = compose.ProcessState[*State](ctx, func(ctx context.Context, state *State) error {
			output.Content = state.Content
			output.ToolCalls = state.ToolCalls
			return nil
		})
		return output, nil
	})

	//chatModel执行完之后把 输出存一下
	chatModelPostHandler := func(ctx context.Context, input *schema.Message, state *State) (output *schema.Message, err error) {
		//fmt.Println("工具使用测试:", input)
		state.ToolCalls = input.ToolCalls
		state.Content = input.Content
		return input, nil
	}

	g := compose.NewGraph[[]*schema.Message, types.AgentResponse](compose.WithGenLocalState(initState))

	err = g.AddChatModelNode("model", aiChatModel, compose.WithStatePostHandler(chatModelPostHandler))
	if err != nil {
		panic("添加节点失败," + err.Error())
	}

	err = g.AddToolsNode("tools", ToolsNode)
	if err != nil {
		panic("添加节点失败," + err.Error())
	}

	err = g.AddLambdaNode("lambda1", lambda1)
	if err != nil {
		panic("添加节点失败" + err.Error())
	}
	err = g.AddLambdaNode("lambda2", lambda2)
	if err != nil {
		panic("添加节点失败" + err.Error())
	}

	//开始连接这些节点

	err = g.AddEdge(compose.START, "model")
	if err != nil {
		panic("添加边失败" + err.Error())
	}

	//创建边一个分支
	err = g.AddBranch("model", compose.NewGraphBranch(func(ctx context.Context, in *schema.Message) (endNode string, err error) {
		if len(in.ToolCalls) > 0 {
			return "tools", nil
		}
		return "lambda1", nil
	}, map[string]bool{
		"tools":   true,
		"lambda1": true,
	}))
	if err != nil {
		panic("创建分支失败" + err.Error())
	}

	err = g.AddEdge("tools", "lambda2")
	if err != nil {
		panic("创建边失败" + err.Error())
	}

	err = g.AddEdge("lambda1", "lambda2")
	if err != nil {
		panic("创建边失败" + err.Error())
	}

	err = g.AddEdge("lambda2", compose.END)
	if err != nil {
		panic("创建边失败" + err.Error())
	}

	agent, err := g.Compile(ctx)
	if err != nil {
		panic("编译错误" + err.Error())
	}

	aiChatClient.Agent = agent

	return &aiChatClient
}

func (a *AiChatClient) SendMessage(ctx context.Context, messages []*entity.Message) (types.AgentResponse, error) {

	input := messagesDo2Input(messages)

	resp, err := a.Agent.Invoke(ctx, input)

	if err != nil {
		zlog.Errorf("模型调用失败%v", err)
		return types.AgentResponse{}, err
	}
	return resp, nil
}

// 传入文本生成导图
func (a *AiChatClient) GenerateMindMap(ctx context.Context, text, userID string) (string, error) {
	message := initGenerateMindMapMessage(text, userID)

	resp, err := a.ToolAiClient.Generate(ctx, message)
	if err != nil {
		zlog.Errorf("模型调用失败%v", err)
		return "", err
	}
	return resp.Content, nil
}
