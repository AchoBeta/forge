# Backend启动指南

## 快速启动

1. 安装依赖:
```bash
pip install -r requirements.txt
```

2. 启动服务:
```bash
python app.py
```

服务将在 http://localhost:8000 启动

## API端点

- `GET /` - 查看API信息
- `GET /health` - 健康检查
- `POST /generate` - 生成思维导图

## API文档

启动服务后访问:
- Swagger UI: http://localhost:8000/docs
- ReDoc: http://localhost:8000/redoc

## 环境变量

- `PORT`: 服务端口 (默认: 8000)

## 开发说明

本系统使用FastAPI框架开发，提供RESTful API接口。

核心功能:
- 思维导图生成算法
- 数据验证和错误处理
- CORS支持
- 结构化响应
