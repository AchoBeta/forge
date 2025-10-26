# Frontend说明

## 启动方式

### 方式1: 直接打开文件
在浏览器中直接打开 `index.html`

### 方式2: 使用本地服务器
```bash
python -m http.server 3000
```
然后访问 http://localhost:3000

## 文件说明

- `index.html` - 主页面结构
- `style.css` - 样式和动画
- `app.js` - 前端业务逻辑

## 配置

在 `app.js` 中可以修改API地址:
```javascript
const API_BASE_URL = 'http://localhost:8000';
```

## 浏览器要求

支持现代浏览器:
- Chrome/Edge 90+
- Firefox 88+
- Safari 14+

## 功能特性

- 响应式设计
- 流畅的动画效果
- 交互式思维导图展示
- 实时参数调整
