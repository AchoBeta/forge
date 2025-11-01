package mindmapservice

import (
	"context"
	"errors"
	"forge/biz/entity"
	"forge/biz/repo"
	"forge/biz/types"
	"forge/pkg/log/zlog"
	"forge/util"
)

// 错误定义
var (
	ErrMindMapNotFound      = errors.New("思维导图不存在")
	ErrMindMapAlreadyExists = errors.New("思维导图已存在")
	ErrInvalidParams        = errors.New("参数无效")
	ErrPermissionDenied     = errors.New("权限不足")
	ErrInternalError        = errors.New("内部错误")
)

// MindMapServiceImpl 思维导图服务实现
type MindMapServiceImpl struct {
	mindMapRepo repo.IMindMapRepo
}

func NewMindMapServiceImpl(mindMapRepo repo.IMindMapRepo) *MindMapServiceImpl {
	return &MindMapServiceImpl{
		mindMapRepo: mindMapRepo,
	}
}

// CreateMindMap 创建思维导图（用户只能创建自己的思维导图）
func (s *MindMapServiceImpl) CreateMindMap(ctx context.Context, req *types.CreateMindMapParams) (*entity.MindMap, error) {
	// 从JWT token上下文中获取用户信息
	user, ok := entity.GetUser(ctx)
	if !ok {
		zlog.CtxErrorf(ctx, "failed to get user from context")
		return nil, ErrPermissionDenied
	}

	// 参数校验
	if req.Title == "" {
		zlog.CtxErrorf(ctx, "title is required")
		return nil, ErrInvalidParams
	}
	if req.Layout == "" {
		zlog.CtxErrorf(ctx, "layout is required")
		return nil, ErrInvalidParams
	}

	// 生成思维导图ID
	mapID, err := util.GenerateStringID()
	if err != nil {
		zlog.CtxErrorf(ctx, "failed to generate map id: %v", err)
		return nil, ErrInternalError
	}

	// 构建实体
	mindMap := &entity.MindMap{
		MapID:  mapID,
		UserID: user.UserID, // 从JWT token中获取的用户ID
		Title:  req.Title,
		Desc:   req.Desc,
		Layout: req.Layout,
		Data:   req.Data,
	}

	// 实体校验
	if err := mindMap.Validate(); err != nil {
		zlog.CtxErrorf(ctx, "mindmap validation failed: %v", err)
		return nil, err
	}

	// 持久化
	if err := s.mindMapRepo.CreateMindMap(ctx, mindMap); err != nil {
		zlog.CtxErrorf(ctx, "failed to create mindmap: %v", err)
		return nil, ErrInternalError
	}

	zlog.CtxInfof(ctx, "mindmap created successfully, mapID: %s, userID: %s", mapID, user.UserID)
	return mindMap, nil
}

// GetMindMap 获取思维导图（用户只能获取自己的思维导图）
func (s *MindMapServiceImpl) GetMindMap(ctx context.Context, mapID string) (*entity.MindMap, error) {
	// 从JWT token上下文中获取用户信息
	user, ok := entity.GetUser(ctx)
	if !ok {
		zlog.CtxErrorf(ctx, "failed to get user from context")
		return nil, ErrPermissionDenied
	}

	// 参数校验
	if mapID == "" {
		zlog.CtxErrorf(ctx, "mapID is required")
		return nil, ErrInvalidParams
	}

	// 构建查询条件（包含用户ID验证）
	query := repo.NewMindMapQueryByID(user.UserID, mapID)

	// 查询思维导图
	mindMap, err := s.mindMapRepo.GetMindMap(ctx, query)
	if err != nil {
		zlog.CtxErrorf(ctx, "failed to get mindmap: %v", err)
		return nil, ErrInternalError
	}
	if mindMap == nil {
		zlog.CtxWarnf(ctx, "mindmap not found or permission denied, mapID: %s, userID: %s", mapID, user.UserID)
		return nil, ErrMindMapNotFound
	}

	zlog.CtxInfof(ctx, "mindmap retrieved successfully, mapID: %s, userID: %s", mapID, user.UserID)
	return mindMap, nil
}

// ListMindMaps 获取思维导图列表（用户只能获取自己的思维导图列表）
func (s *MindMapServiceImpl) ListMindMaps(ctx context.Context, req *types.ListMindMapsParams) ([]*entity.MindMap, int64, error) {
	// 从JWT token上下文中获取用户信息
	user, ok := entity.GetUser(ctx)
	if !ok {
		zlog.CtxErrorf(ctx, "failed to get user from context")
		return nil, 0, ErrPermissionDenied
	}

	// 构建查询条件（强制包含用户ID）
	query := repo.NewMindMapQueryForList(user.UserID, req.Page, req.PageSize)

	// 添加可选筛选条件
	if req.Title != "" {
		query.Title = req.Title
	}
	if req.Layout != "" {
		query.Layout = req.Layout
	}

	// 查询列表
	mindMaps, total, err := s.mindMapRepo.ListMindMaps(ctx, query)
	if err != nil {
		zlog.CtxErrorf(ctx, "failed to list mindmaps: %v", err)
		return nil, 0, ErrInternalError
	}

	zlog.CtxInfof(ctx, "mindmaps listed successfully, userID: %s, count: %d, total: %d", user.UserID, len(mindMaps), total)
	return mindMaps, total, nil
}

// UpdateMindMap 更新思维导图（用户只能更新自己的思维导图）
func (s *MindMapServiceImpl) UpdateMindMap(ctx context.Context, mapID string, req *types.UpdateMindMapParams) error {
	// 从JWT token上下文中获取用户信息
	user, ok := entity.GetUser(ctx)
	if !ok {
		zlog.CtxErrorf(ctx, "failed to get user from context")
		return ErrPermissionDenied
	}

	// 参数校验
	if mapID == "" {
		zlog.CtxErrorf(ctx, "mapID is required")
		return ErrInvalidParams
	}

	// 首先验证思维导图存在且属于当前用户
	existingMindMap, err := s.GetMindMap(ctx, mapID)
	if err != nil {
		return err // GetMindMap已经包含权限验证
	}
	if existingMindMap == nil {
		return ErrMindMapNotFound
	}

	// 构建更新信息
	updateInfo := &repo.MindMapUpdateInfo{
		MapID:  mapID,
		UserID: user.UserID, // 确保只能更新自己的思维导图
		Title:  req.Title,
		Desc:   req.Desc,
		Layout: req.Layout,
		Data:   req.Data,
	}

	// 如果有标题更新，进行校验
	if req.Title != nil {
		if *req.Title == "" {
			return ErrInvalidParams
		}
		if len(*req.Title) > 100 {
			return entity.ErrTitleTooLong
		}
	}

	// 如果有描述更新，进行校验
	if req.Desc != nil && len(*req.Desc) > 500 {
		return entity.ErrDescTooLong
	}

	// 如果有布局更新，进行校验
	if req.Layout != nil && *req.Layout == "" {
		return entity.ErrInvalidLayout
	}

	// 执行更新
	if err := s.mindMapRepo.UpdateMindMap(ctx, updateInfo); err != nil {
		zlog.CtxErrorf(ctx, "failed to update mindmap: %v", err)
		return ErrInternalError
	}

	zlog.CtxInfof(ctx, "mindmap updated successfully, mapID: %s, userID: %s", mapID, user.UserID)
	return nil
}

// DeleteMindMap 删除思维导图（用户只能删除自己的思维导图）
func (s *MindMapServiceImpl) DeleteMindMap(ctx context.Context, mapID string) error {
	// 从JWT token上下文中获取用户信息
	user, ok := entity.GetUser(ctx)
	if !ok {
		zlog.CtxErrorf(ctx, "failed to get user from context")
		return ErrPermissionDenied
	}

	// 参数校验
	if mapID == "" {
		zlog.CtxErrorf(ctx, "mapID is required")
		return ErrInvalidParams
	}

	// 首先验证思维导图存在且属于当前用户
	existingMindMap, err := s.GetMindMap(ctx, mapID)
	if err != nil {
		return err // GetMindMap已经包含权限验证
	}
	if existingMindMap == nil {
		return ErrMindMapNotFound
	}

	// 执行删除（软删除）
	if err := s.mindMapRepo.DeleteMindMap(ctx, mapID); err != nil {
		zlog.CtxErrorf(ctx, "failed to delete mindmap: %v", err)
		return ErrInternalError
	}

	zlog.CtxInfof(ctx, "mindmap deleted successfully, mapID: %s, userID: %s", mapID, user.UserID)
	return nil
}
