package def

// 创建请求
type CreateMindMapReq struct {
	Title  string      `json:"title" binding:"required,max=100"`
	Desc   string      `json:"desc" binding:"max=500"`
	Layout string      `json:"layout" binding:"required"`
	Data   MindMapData `json:"data" binding:"required"`
}

// 列表查询请求
type ListMindMapsReq struct {
	Title    string `form:"title"`
	Layout   string `form:"layout"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=20"`
}

// 更新请求
type UpdateMindMapReq struct {
	Title  *string      `json:"title,omitempty" binding:"omitempty,max=100"`
	Desc   *string      `json:"desc,omitempty" binding:"omitempty,max=500"`
	Layout *string      `json:"layout,omitempty"`
	Data   *MindMapData `json:"data,omitempty"`
}

// 思维导图DTO
type MindMapDTO struct {
	MapID     string      `json:"map_id"`
	Title     string      `json:"title"`
	Desc      string      `json:"desc"`
	Layout    string      `json:"layout"`
	Data      MindMapData `json:"data"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
}

// 节点数据DTO
type NodeData struct {
	Text string `json:"text"`
	// 可扩展其他节点属性，如颜色、图标等
}

// 思维导图数据DTO - 递归树结构
type MindMapData struct {
	Data     NodeData      `json:"data"`     // 节点数据
	Children []MindMapData `json:"children"` // 子节点（递归结构）
}

// 响应DTO
type CreateMindMapResp struct {
	*MindMapDTO
}

type GetMindMapResp struct {
	*MindMapDTO
}

type ListMindMapsResp struct {
	List     []*MindMapDTO `json:"list"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}

type UpdateMindMapResp struct {
	Success bool `json:"success"`
}

type DeleteMindMapResp struct {
	Success bool `json:"success"`
}
