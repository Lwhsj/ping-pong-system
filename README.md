# 🏓 乒乓球比赛管理系统

一套完整的乒乓球比赛管理解决方案，包含后端 API 服务、观众/教练 Web 端和裁判桌面客户端，并集成 AI 智能分析能力。

## 系统架构

```
┌─────────────────────────────────────────────────────────────────┐
│                        系统整体架构                               │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌──────────────┐     HTTP/REST     ┌──────────────────────┐   │
│  │  Tennis Web  │ ◄───────────────► │                      │   │
│  │  (Vue 3)    │   轮询比分更新      │                      │   │
│  │  观众/教练端  │                    │   Ping Pong Backend  │   │
│  └──────────────┘                    │       (Go/Gin)       │   │
│                                      │                      │   │
│  ┌──────────────┐     HTTP/REST     │  ┌────────────────┐  │   │
│  │Referee Client│ ◄───────────────► │  │  MySQL (GORM)  │  │   │
│  │ (Electron)  │   记分/视频上传     │  └────────────────┘  │   │
│  │   裁判端     │                    │                      │   │
│  └──────────────┘                    │  ┌────────────────┐  │   │
│                                      │  │ OpenAI LLM API │  │   │
│                                      │  │  (AI 分析)      │  │   │
│                                      │  └────────────────┘  │   │
│                                      └──────────────────────┘   │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

## 技术栈

| 组件 | 技术 | 说明 |
|------|------|------|
| 后端 | Go 1.24 + Gin + GORM | REST API 服务 |
| 数据库 | MySQL | 数据持久化 |
| Web 前端 | Vue 3 + Vite + Element Plus + ECharts | 观众/教练端 |
| 裁判客户端 | Electron + Vue 3 + Vite | 桌面应用 |
| AI 分析 | OpenAI API (GPT-4o-mini) | 比赛智能分析 |
| 视频存储 | 本地文件系统 | WebM 格式回放片段 |

## 功能介绍

### 🎯 裁判客户端 (Referee Client)

比赛控制的核心入口，提供桌面级操作体验：

- **比赛创建** — 选择对阵选手、指定首发球员、自动记录日期
- **实时记分** — 一键记录得分，自动计算比分和发球轮换
- **视频录制** — 使用 MediaRecorder API 逐回合录制比赛视频
- **AI 助手** — 比赛中随时询问 AI 进行战术分析，赛后生成完整总结
- **数据导出** — 一键导出比赛数据为 Excel 文件

### 📺 观众/教练 Web 端 (Tennis Web)

面向观众和教练的实时数据展示：

- **实时比分** — 每秒轮询更新，大屏显示当前比分
- **历史比赛** — 按日期、选手、状态筛选历史比赛记录
- **比赛详情** — 回合历史表格 + ECharts 统计图表可视化
- **数据导出** — 支持导出比赛数据进行离线分析

### 🧠 AI 智能分析

基于 LLM 的比赛分析能力：

- **实时咨询** — 比赛进行中向 AI 提问（战术建议、趋势分析）
- **赛后总结** — 自动生成比赛表现报告
- **上下文感知** — AI 了解完整比赛数据（比分、发球、连续得分等）

## 项目结构

```
ping-pong-system/
├── pingpong-backend/          # Go 后端服务
│   ├── cmd/server/            # 应用入口
│   ├── internal/
│   │   ├── config/            # 环境配置
│   │   ├── database/          # 数据库连接
│   │   ├── handler/           # HTTP 请求处理
│   │   ├── model/             # 数据模型
│   │   ├── service/           # 业务逻辑 + LLM 集成
│   │   ├── router/            # 路由定义
│   │   ├── middleware/        # CORS 中间件
│   │   └── dto/               # 数据传输对象
│   ├── migrations/            # SQL 建表 + 种子数据
│   └── uploads/               # 视频文件存储
├── tennis-web/                # Vue 3 Web 前端
│   └── src/
│       ├── views/             # 页面组件
│       ├── api/               # API 客户端
│       ├── router/            # 路由配置
│       └── layout/            # 布局组件
├── referee-client/            # Electron 裁判客户端
│   ├── electron/              # 主进程
│   └── src/
│       ├── views/             # 页面组件
│       ├── services/          # API 服务
│       └── composables/       # 组合式函数
└── doc/                       # 设计文档
```

## 数据模型

```
┌──────────┐       ┌──────────────┐       ┌───────────┐
│  Player  │       │    Match     │       │   Rally   │
├──────────┤       ├──────────────┤       ├───────────┤
│ id       │◄──┐   │ id           │◄──┐   │ id        │
│ name     │   ├──│ player1_id   │   │   │ match_id ─┤──►Match
│ sex      │   ├──│ player2_id   │   │   │ rally_num │
│ age      │   │   │ date         │   └──│ scorer    │
└──────────┘   │   │ status       │       │ server    │
               │   │ started_at   │       │ timestamp │
               └──│ first_server │       │ video_file│
                   └──────────────┘       │ score_p1  │
                                          │ score_p2  │
                                          └───────────┘
```

## API 概览

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/match/start` | 创建比赛 |
| GET | `/api/match/:id` | 获取比赛信息 |
| GET | `/api/match/:id/current` | 获取实时比分 |
| POST | `/api/match/:id/finish` | 结束比赛 |
| GET | `/api/matches` | 历史比赛列表 |
| GET | `/api/match/:id/detail` | 回合详情 |
| GET | `/api/match/:id/stats` | 比赛统计 |
| GET | `/api/match/:id/export` | 导出 Excel |
| POST | `/api/rally` | 记录得分 |
| POST | `/api/upload/video` | 上传视频 |
| GET | `/api/video/:fileName` | 获取视频流 |
| POST | `/api/agent/match/:id/analyze` | AI 分析比赛 |
| POST | `/api/agent/chat` | AI 对话 |
| GET | `/api/players` | 选手列表 |

## 快速开始

### 环境要求

- Go 1.24+
- Node.js 18+
- MySQL 8.0+

### 1. 数据库初始化

```bash
mysql -u root -p -e "CREATE DATABASE pingpong"
mysql -u root -p pingpong < pingpong-backend/migrations/schema.sql
mysql -u root -p pingpong < pingpong-backend/migrations/seed.sql
```

### 2. 启动后端

```bash
cd pingpong-backend
cp .env.example .env   # 编辑配置
go run ./cmd/server
```

服务默认运行在 `http://localhost:8080`

### 3. 启动 Web 前端

```bash
cd tennis-web
npm install
npm run dev
```

开发服务器运行在 `http://localhost:5173`

### 4. 启动裁判客户端

```bash
cd referee-client
npm install
npm run dev
```

自动打开 Electron 桌面窗口。

### 环境变量配置

在 `pingpong-backend/.env` 中配置：

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `APP_PORT` | 8080 | 服务端口 |
| `DB_HOST` | 127.0.0.1 | 数据库地址 |
| `DB_PORT` | 3306 | 数据库端口 |
| `DB_NAME` | pingpong | 数据库名称 |
| `DB_USER` | root | 数据库用户 |
| `DB_PASSWORD` | — | 数据库密码 |
| `UPLOAD_DIR` | uploads | 视频存储目录 |
| `MAX_UPLOAD_MB` | 50 | 视频最大尺寸 (MB) |
| `AGENT_ENABLED` | true | 启用 AI 分析 |
| `LLM_BASE_URL` | https://api.openai.com/v1 | LLM API 地址 |
| `LLM_API_KEY` | — | LLM API 密钥 |
| `LLM_MODEL` | gpt-4o-mini | LLM 模型 |

## 工作流程

```
裁判创建比赛 → 选手对阵开始 → 裁判逐回合记分（同时录像）
     │                              │
     │                              ▼
     │              Web 端实时展示比分（1秒轮询）
     │                              │
     ▼                              ▼
比赛结束 → AI 生成赛后分析报告 → 数据导出 Excel
```

## 开发说明

- 后端使用标准 Go 项目布局，业务逻辑集中在 `internal/service/`
- 前端统一使用 Vue 3 Composition API + `<script setup>` 语法
- 裁判端通过 Electron 打包为桌面应用，使用 hash 路由模式
- Web 端通过 Vite proxy 或 Nginx 代理后端 API
- 视频以 WebM 格式按回合存储，文件名关联 Rally 记录

