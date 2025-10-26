# Forge 羲图系统 - 基于生成式AI的思维导图系统

🔥 **AchoBeta Forge** 是一个基于生成式AI技术的智能思维导图生成系统。

## 功能特性

- 🤖 **AI智能生成**: 基于输入主题自动生成结构化思维导图
- 🎨 **可视化展示**: 美观的思维导图可视化界面
- ⚙️ **灵活配置**: 支持自定义深度和分支数量
- 🚀 **快速响应**: 实时生成，即时反馈
- 🌐 **现代化架构**: FastAPI后端 + 原生JavaScript前端

## 项目结构

```
forge/
├── backend/                # 后端服务
│   ├── app.py             # FastAPI应用主文件
│   └── requirements.txt   # Python依赖
├── frontend/              # 前端界面
│   ├── index.html        # 主页面
│   ├── style.css         # 样式文件
│   └── app.js            # 前端逻辑
└── README.md             # 项目说明
```

## 快速开始

### 环境要求

- Python 3.8+
- pip (Python包管理器)
- 现代浏览器 (Chrome, Firefox, Safari, Edge)

### 安装步骤

#### 1. 安装后端依赖

```bash
cd backend
pip install -r requirements.txt
```

#### 2. 启动后端服务

```bash
cd backend
python app.py
```

后端服务将在 `http://localhost:8000` 启动

#### 3. 访问前端界面

在浏览器中打开 `frontend/index.html` 文件，或使用本地Web服务器：

```bash
# 使用Python内置服务器
cd frontend
python -m http.server 3000
```

然后访问 `http://localhost:3000`

## 使用说明

1. **输入主题**: 在输入框中输入想要展开的主题（例如：人工智能、机器学习等）
2. **调整参数**: 
   - **深度**: 控制思维导图的层级深度（1-5层）
   - **分支数**: 控制每个节点的子节点数量（1-8个）
3. **生成导图**: 点击"生成思维导图"按钮
4. **查看结果**: 系统将自动生成并展示思维导图

## API文档

### 生成思维导图

**接口**: `POST /generate`

**请求体**:
```json
{
  "topic": "人工智能",
  "depth": 3,
  "branches": 3
}
```

**响应**:
```json
{
  "root": {
    "id": "root",
    "text": "人工智能",
    "level": 0,
    "children": [...]
  },
  "metadata": {
    "topic": "人工智能",
    "depth": 3,
    "branches": 3,
    "generated_by": "Forge AI System"
  }
}
```

### 健康检查

**接口**: `GET /health`

**响应**:
```json
{
  "status": "healthy",
  "service": "forge-mindmap"
}
```

## 技术栈

### 后端
- **FastAPI**: 现代化的Python Web框架
- **Pydantic**: 数据验证和序列化
- **Uvicorn**: ASGI服务器

### 前端
- **HTML5**: 页面结构
- **CSS3**: 样式和动画
- **JavaScript (ES6+)**: 交互逻辑和数据处理

## 扩展功能

本系统提供了基础的思维导图生成功能，可以通过以下方式进行扩展：

1. **集成真实AI模型**:
   - OpenAI GPT-3/4 API
   - 本地部署的大语言模型（如LLaMA、ChatGLM等）
   - Azure OpenAI Service

2. **增强功能**:
   - 导出功能（PNG、PDF、JSON等格式）
   - 思维导图编辑功能
   - 协作功能
   - 历史记录保存
   - 自定义主题和样式

3. **性能优化**:
   - 缓存机制
   - 批量生成
   - 异步处理

## 贡献指南

欢迎提交Issue和Pull Request！

## 开源协议

MIT License

## 联系方式

- 项目: AchoBeta Forge 羲图系统
- 年份: 2025
- 小组: G组
- 用途: 复试项目

---

**注意**: 本项目为教育和演示目的，实际生产环境使用时请确保配置适当的AI模型和安全措施。
