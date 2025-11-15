package generationservice

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"forge/biz/entity"
	"forge/biz/repo"
	"forge/biz/types"
	"forge/pkg/log/zlog"
	"forge/util"
)

type GenerationService struct {
	generationRepo repo.IGenerationRepo
	aiChatRepo     repo.AiChatRepo
	mindMapRepo    repo.IMindMapRepo
}

func NewGenerationService(generationRepo repo.IGenerationRepo, aiChatRepo repo.AiChatRepo, mindMapRepo repo.IMindMapRepo) types.IGenerationService {
	return &GenerationService{
		generationRepo: generationRepo,
		aiChatRepo:     aiChatRepo,
		mindMapRepo:    mindMapRepo,
	}
}

// GetBatchWithResults 获取批次及其结果
func (g *GenerationService) GetBatchWithResults(ctx context.Context, batchID string) (*entity.GenerationBatch, []*entity.GenerationResult, error) {
	// 获取用户信息
	user, ok := entity.GetUser(ctx)
	if !ok {
		return nil, nil, fmt.Errorf("无法获取用户信息")
	}

	// 获取批次信息
	batch, err := g.generationRepo.GetGenerationBatch(ctx, batchID, user.UserID)
	if err != nil {
		return nil, nil, err
	}

	// 获取结果列表
	results, err := g.generationRepo.GetGenerationResultsByBatchID(ctx, batchID)
	if err != nil {
		return nil, nil, err
	}

	return batch, results, nil
}

// ListUserBatches 获取用户批次列表
func (g *GenerationService) ListUserBatches(ctx context.Context, userID string, page, pageSize int) ([]*entity.GenerationBatch, int64, error) {
	return g.generationRepo.ListUserGenerationBatches(ctx, userID, page, pageSize)
}

// LabelResultWithSave 标记结果并可能保存导图
func (g *GenerationService) LabelResultWithSave(ctx context.Context, resultID string, label int) (*entity.MindMap, error) {
	// 验证标签值
	if label != -1 && label != 0 && label != 1 {
		return nil, fmt.Errorf("标签值必须是-1、0或1")
	}

	// 更新标签
	if err := g.generationRepo.UpdateGenerationResultLabel(ctx, resultID, label); err != nil {
		return nil, err
	}

	// 如果标记为正值(1)，保存导图到正式系统
	if label == 1 {
		mindMap, err := g.SaveSelectedMindMap(ctx, resultID)
		if err != nil {
			zlog.CtxWarnf(ctx, "保存选中导图失败: %v", err)
			// 返回错误，因为用户期望保存成功
			return nil, fmt.Errorf("保存选中导图失败: %w", err)
		}
		return mindMap, nil
	}

	return nil, nil
}

// LabelResult 标记结果（兼容旧接口）
func (g *GenerationService) LabelResult(ctx context.Context, resultID string, label int) error {
	_, err := g.LabelResultWithSave(ctx, resultID, label)
	return err
}

// SFTRecord SFT训练记录结构
type SFTRecord struct {
	Messages []SFTMessage `json:"messages"`
	Thinking string       `json:"thinking"`
}

type SFTMessage struct {
	Role             string   `json:"role"`
	Content          string   `json:"content"`
	LossWeight       *float64 `json:"loss_weight,omitempty"`
	ReasoningContent *string  `json:"reasoning_content,omitempty"`
}

// ExportSFTData 导出SFT数据
func (g *GenerationService) ExportSFTData(ctx context.Context, startDate, endDate, userID string) (string, error) {
	// 获取已标记的结果
	results, err := g.generationRepo.GetLabeledResults(ctx, userID, startDate, endDate)
	if err != nil {
		return "", err
	}

	if len(results) == 0 {
		return "", nil
	}

	var jsonlLines []string

	// 选择SFT正样本
	selectedResults := g.selectSFTSamples(ctx, results)

	// 为每个选中的结果生成JSONL记录
	for _, result := range selectedResults {

		// 获取对话记录
		conversation, err := g.aiChatRepo.GetConversation(ctx, result.ConversationID, userID)
		if err != nil {
			zlog.CtxWarnf(ctx, "获取对话记录失败 conversationID:%s, err:%v", result.ConversationID, err)
			continue
		}

		// 构建SFT记录
		record, err := g.buildSFTRecord(conversation, result.Label)
		if err != nil {
			zlog.CtxWarnf(ctx, "构建SFT记录失败 conversationID:%s, err:%v", result.ConversationID, err)
			continue
		}

		// 转换为JSON字符串
		jsonBytes, err := json.Marshal(record)
		if err != nil {
			zlog.CtxWarnf(ctx, "序列化SFT记录失败 conversationID:%s, err:%v", result.ConversationID, err)
			continue
		}

		jsonlLines = append(jsonlLines, string(jsonBytes))
		zlog.CtxInfof(ctx, "添加SFT样本 resultID:%s", result.ResultID)
	}

	return strings.Join(jsonlLines, "\n"), nil
}

// buildSFTRecord 构建SFT记录
func (g *GenerationService) buildSFTRecord(conversation *entity.Conversation, label int) (*SFTRecord, error) {
	if len(conversation.Messages) < 2 {
		return nil, fmt.Errorf("对话消息不足")
	}

	var messages []SFTMessage

	// 处理每条消息
	for i, message := range conversation.Messages {
		sftMessage := SFTMessage{
			Role:    strings.ToLower(message.Role),
			Content: message.Content,
		}

		// 设置loss_weight
		if sftMessage.Role == "assistant" && i == len(conversation.Messages)-1 {
			// 只有最后一条assistant消息才设置loss_weight
			if label == 1 {
				lossWeight := 1.0
				sftMessage.LossWeight = &lossWeight
			} else if label == -1 {
				lossWeight := 0.0
				sftMessage.LossWeight = &lossWeight
			}
		}

		messages = append(messages, sftMessage)
	}

	// 构建完整记录
	record := &SFTRecord{
		Messages: messages,
		Thinking: "disabled", // 使用当前模型配置
	}

	return record, nil
}

// ExportDPOData 导出DPO数据
func (g *GenerationService) ExportDPOData(ctx context.Context, startDate, endDate, userID string) (string, error) {
	// 获取已标记的结果（正负样本）
	labeledResults, err := g.generationRepo.GetLabeledResults(ctx, userID, startDate, endDate)
	if err != nil {
		return "", fmt.Errorf("获取已标记结果失败: %w", err)
	}

	if len(labeledResults) == 0 {
		return "", nil
	}

	// 按批次ID分组
	batchGroups := make(map[string][]*entity.GenerationResult)
	for _, result := range labeledResults {
		batchGroups[result.BatchID] = append(batchGroups[result.BatchID], result)
	}

	var dpoRecords []string

	// 为每个批次生成DPO对比对
	for batchID, results := range batchGroups {
		// 分离正负样本
		positiveResults := g.selectPositiveSamples(ctx, results, batchID)
		var negativeResults []*entity.GenerationResult

		for _, result := range results {
			if result.Label == -1 {
				negativeResults = append(negativeResults, result)
			}
		}

		// 生成正负样本对（限制配对数量，避免数据爆炸）
		maxPairsPerPositive := 3 // 每个正样本最多配对3个负样本
		for _, positive := range positiveResults {
			pairCount := 0
			for _, negative := range negativeResults {
				if pairCount >= maxPairsPerPositive {
					break // 限制配对数量
				}
				dpoRecord, err := g.buildDPORecord(ctx, positive, negative, userID)
				if err != nil {
					zlog.CtxWarnf(ctx, "构建DPO记录失败 batchID:%s, err:%v", batchID, err)
					continue
				}
				dpoRecords = append(dpoRecords, dpoRecord)
				pairCount++
			}
			zlog.CtxInfof(ctx, "批次 %s：正样本 %s 生成了 %d 个DPO配对", batchID, positive.ResultID, pairCount)
		}
	}

	return strings.Join(dpoRecords, "\n"), nil
}

// selectPositiveSamples 选择正样本（简化版）
func (g *GenerationService) selectPositiveSamples(ctx context.Context, results []*entity.GenerationResult, batchID string) []*entity.GenerationResult {
	var positiveResults []*entity.GenerationResult

	// 收集所有正样本 (label = 1)
	for _, result := range results {
		if result.Label == 1 {
			positiveResults = append(positiveResults, result)
		}
	}

	zlog.CtxInfof(ctx, "批次 %s：选择正样本 %d 个", batchID, len(positiveResults))
	return positiveResults
}

// selectSFTSamples 选择SFT正样本（简化版）
func (g *GenerationService) selectSFTSamples(ctx context.Context, results []*entity.GenerationResult) []*entity.GenerationResult {
	var positiveResults []*entity.GenerationResult

	// SFT只使用正样本 (label = 1)
	for _, result := range results {
		if result.Label == 1 {
			positiveResults = append(positiveResults, result)
		}
	}

	zlog.CtxInfof(ctx, "SFT导出：选择正样本 %d 个", len(positiveResults))
	return positiveResults
}

// buildDPORecord 构建DPO记录
func (g *GenerationService) buildDPORecord(ctx context.Context, positive, negative *entity.GenerationResult, userID string) (string, error) {
	// 获取正样本对话
	positiveConversation, err := g.aiChatRepo.GetConversation(ctx, positive.ConversationID, userID)
	if err != nil {
		return "", fmt.Errorf("获取正样本对话失败: %w", err)
	}

	// 构建DPO格式记录
	dpoRecord := &DPORecord{
		Messages: make([]DPOMessage, 0, len(positiveConversation.Messages)-1), // 不包含最后一条assistant消息
	}

	// 添加system和user消息
	for _, message := range positiveConversation.Messages {
		if message.Role == entity.ASSISTANT {
			// 跳过assistant消息，最后会特殊处理
			continue
		}
		dpoMessage := DPOMessage{
			Role:    strings.ToLower(message.Role),
			Content: message.Content,
		}
		dpoRecord.Messages = append(dpoRecord.Messages, dpoMessage)
	}

	// 添加最后的assistant消息，包含chosen和rejected
	finalMessage := DPOMessage{
		Role:     "assistant",
		Chosen:   positive.MapJSON,
		Rejected: negative.MapJSON,
	}
	dpoRecord.Messages = append(dpoRecord.Messages, finalMessage)

	// 转换为JSON字符串
	jsonBytes, err := json.Marshal(dpoRecord)
	if err != nil {
		return "", fmt.Errorf("序列化DPO记录失败: %w", err)
	}

	return string(jsonBytes), nil
}

// DPORecord DPO训练记录结构
type DPORecord struct {
	Messages []DPOMessage `json:"messages"`
}

type DPOMessage struct {
	Role     string `json:"role"`
	Content  string `json:"content,omitempty"`
	Chosen   string `json:"chosen,omitempty"`
	Rejected string `json:"rejected,omitempty"`
}

// ExportSFTDataToFile 导出SFT数据到文件
func (g *GenerationService) ExportSFTDataToFile(ctx context.Context, startDate, endDate, userID string) (string, error) {
	// 获取JSONL数据
	jsonlData, err := g.ExportSFTData(ctx, startDate, endDate, userID)
	if err != nil {
		return "", err
	}

	if jsonlData == "" {
		return "", fmt.Errorf("没有可导出的数据")
	}

	// 生成文件名
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("SFT_Text_Sample_%s_%s.jsonl", userID, timestamp)

	// 这里可以根据需要实现文件写入逻辑
	// 比如写入到临时文件或上传到云存储
	// 当前只返回文件名，实际文件操作可在handler层处理

	return filename, nil
}

// SaveSelectedMindMap 保存选中的导图到正式系统
func (g *GenerationService) SaveSelectedMindMap(ctx context.Context, resultID string) (*entity.MindMap, error) {
	// 获取用户信息
	user, ok := entity.GetUser(ctx)
	if !ok {
		return nil, fmt.Errorf("无法获取用户信息")
	}

	// 获取生成结果
	result, err := g.generationRepo.GetGenerationResult(ctx, resultID)
	if err != nil {
		return nil, fmt.Errorf("获取生成结果失败: %w", err)
	}

	// 检查结果是否属于当前用户（通过批次验证）
	batch, err := g.generationRepo.GetGenerationBatch(ctx, result.BatchID, user.UserID)
	if err != nil {
		return nil, fmt.Errorf("权限验证失败: %w", err)
	}

	// 生成MapID
	mapID, err := util.GenerateStringID()
	if err != nil {
		return nil, fmt.Errorf("生成导图ID失败: %w", err)
	}

	// 解析AI返回的JSON到临时结构
	var aiResponse map[string]interface{}
	if err := json.Unmarshal([]byte(result.MapJSON), &aiResponse); err != nil {
		zlog.CtxErrorf(ctx, "JSON解析失败，完整JSON内容: %s", result.MapJSON)
		return nil, fmt.Errorf("解析导图JSON失败: %w", err)
	}

	// 提取title和layout
	title := fmt.Sprintf("批量生成-%s", batch.BatchID[:8])
	layout := "mindMap"
	desc := "From batch generation"

	if aiTitle, ok := aiResponse["title"].(string); ok && aiTitle != "" {
		title = aiTitle
	}
	if aiLayout, ok := aiResponse["layout"].(string); ok && aiLayout != "" {
		layout = aiLayout
	}

	// 提取root数据并转换为MindMapData
	rootData, exists := aiResponse["root"]
	if !exists {
		return nil, fmt.Errorf("JSON中缺少root字段")
	}

	rootBytes, err := json.Marshal(rootData)
	if err != nil {
		return nil, fmt.Errorf("序列化root数据失败: %w", err)
	}
	var mindMapData entity.MindMapData
	if err := json.Unmarshal(rootBytes, &mindMapData); err != nil {
		return nil, fmt.Errorf("解析root数据失败: %w", err)
	}

	// 创建MindMap实体（模仿原有CreateMindMap实现）
	mindMap := &entity.MindMap{
		MapID:  mapID,
		UserID: user.UserID,
		Title:  title,
		Desc:   desc,
		Layout: layout,
		Data:   mindMapData, // 直接使用解析好的结构化数据
	}

	// 保存到正式导图系统
	if err := g.mindMapRepo.CreateMindMap(ctx, mindMap); err != nil {
		return nil, fmt.Errorf("保存导图失败: %w", err)
	}

	zlog.CtxInfof(ctx, "成功保存选中导图: mapID=%s, resultID=%s", mindMap.MapID, resultID)

	return mindMap, nil
}

// SaveGenerationBatch 保存批次和结果（事务操作）
func (g *GenerationService) SaveGenerationBatch(ctx context.Context, batch *entity.GenerationBatch, results []*entity.GenerationResult, conversations []*entity.Conversation) error {
	return g.generationRepo.SaveGenerationBatch(ctx, batch, results, conversations)
}
