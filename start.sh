#!/bin/bash
# 一键启动脚本

echo "🔥 启动 Forge 羲图系统..."

# 检查Python
if ! command -v python3 &> /dev/null; then
    echo "❌ 未找到Python3，请先安装Python"
    exit 1
fi

# 安装依赖
echo "📦 安装后端依赖..."
cd backend
pip install -r requirements.txt

# 启动后端
echo "🚀 启动后端服务..."
python3 app.py &
BACKEND_PID=$!

# 等待后端启动
sleep 3

# 启动前端
echo "🌐 启动前端服务..."
cd ../frontend
python3 -m http.server 3000 &
FRONTEND_PID=$!

echo ""
echo "✅ 系统启动成功！"
echo ""
echo "后端服务: http://localhost:8000"
echo "前端界面: http://localhost:3000"
echo "API文档: http://localhost:8000/docs"
echo ""
echo "按 Ctrl+C 停止服务"

# 等待用户中断
trap "kill $BACKEND_PID $FRONTEND_PID; exit" INT
wait
