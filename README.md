# Halo Go

这是一个基于 Go 和现代前端重构中的 Halo 实验版本。

## 进展报告

### 已完成功能（P0-P1）

#### P0 进行中

- 评论嵌套回复基础模型：`Reply`
- 用户注册与验证基础模型：`UserRegistration`、`PasswordReset`
- 角色与权限基础模型：`PolicyRule`、`Role`、`RoleRule`、`UserRole`
- 权限中间件基础：按资源与方法进行权限判断
- Reply API：按评论获取回复、创建回复、审核通过、审核拒绝
- 注册 API：`POST /api/register`、`GET /api/register/verify`
- 写操作路由已开始接入权限中间件

#### 初始化与认证

- `/system/setup` 初始化路由，未初始化时自动跳转
- 4 语言初始化页：简体中文、繁体中文、英语、西班牙语
- 语言切换写入 URL 参数并刷新页面
- 确认密码校验
- JWT 登录
- 用户登录态存储（localStorage + Bearer Token）

#### 后台管理

- 侧边导航控制台壳
- 概览、文章、页面、分类标签、评论、用户、插件主题、设置模块切换
- 后台壳深色主题、active 状态、退出登录

#### 后端 API（50+ 接口）

| 模块 | 接口 |
|------|------|
| 健康 | `GET /api/health` |
| 初始化 | `GET /api/setup/status` `POST /api/setup` |
| 登录 | `POST /api/login` |
| 文章 | `GET/POST /api/posts` `GET/PUT/DELETE /api/posts/:id` |
| 页面 | `GET/POST /api/pages` `GET/PUT/DELETE /api/pages/:id` |
| 分类 | `GET/POST /api/categories` `GET/PUT/DELETE /api/categories/:id` |
| 标签 | `GET/POST /api/tags` `GET/PUT/DELETE /api/tags/:id` |
| 菜单 | `GET/POST /api/menus` `GET/PUT/DELETE /api/menus/:id` |
| 评论 | `GET/POST /api/comments` `GET/DELETE /api/comments/:id` |
| 用户 | `GET/POST /api/users` `GET/PUT/DELETE /api/users/:id` |
| 主题 | `GET /api/themes` |
| 插件 | `GET /api/plugins` `POST /api/plugins/:id/enable` `POST /api/plugins/:id/disable` |
| 附件 | `GET /api/attachments` `POST /api/attachments/upload` |
| 备份 | `GET/POST /api/backups` `GET /api/backups/:id/download` |
| 设置 | `GET/PUT /api/settings/:id` |
| 仪表盘 | `GET /api/dashboard/stats` |
| 搜索 | `GET /api/search` |

#### 数据库与基础设施

- 纯 Go SQLite 驱动（`github.com/glebarez/sqlite`，无 CGO 依赖）
- GORM 模型：User、Post、Page、Category、Tag、Menu、Comment、Theme、Plugin、Attachment、Backup、Setting
- 工作目录自动创建：`db/` `logs/` `themes/` `plugins/` `attachments/` `backups/`
- 环境变量配置：`HALO_WORK_DIR` `HALO_DB_PATH` `HALO_ADDR` `HALO_JWT_SECRET`

#### 构建与部署

- Go 交叉编译（`CGO_ENABLED=1 GOOS=linux GOARCH=amd64`）
- 前端多阶段 Docker 构建（Node → Go → Alpine）
- Vite 开发代理（`/api` → 后端地址）
- `.gitignore` 隔离编译产物与本地运行目录

### 待完成功能（按优先级）

#### P0 — 核心缺失

| 功能 | 说明 |
|------|------|
| 评论嵌套回复 | 原版 Comment + Reply 分离，需实现二级评论与通知触发 |
| 用户注册 | 邮箱验证 Token 机制 |
| 用户权限体系 | 角色（super-admin/operator/guest）+ Policy Rule 基础模型 |
| API 端点权限 | 中间件检查用户角色与端点 Policy |

#### P1 — 平台核心能力

| 功能 | 说明 |
|------|------|
| 插件生命周期 | 安装/卸载/启停基础逻辑 |
| 主题渲染 | 模板引擎（推荐 Pongo2）+ 模板上下文 + 主题端点 |
| 全文搜索 | SQLite FTS5 或 Bleve 索引 |
| 邮箱服务 | 发送验证/密码重置邮件 |

#### P2 — 完善功能

| 功能 | 说明 |
|------|------|
| 双因素认证 | TOTP + 恢复码 |
| 设备管理 | 新设备检测 + 登录通知 |
| 附件存储策略 | OSS/S3 支持 + 策略配置 |
| 通知系统 | 通知中心 + 用户偏好 |
| 缩略图 | 图片处理 |

#### P3 — 生态兼容

| 功能 | 说明 |
|------|------|
| OAuth2 社交登录 | 第三方账号绑定 |
| 插件反向代理路由 | 插件注册自定义 API 端点 |
| 插件共享事件 | 插件间事件派发监听 |
| 备份自动调度 | 定时备份 + 增量备份 |
## 当前状态

- 旧项目已整体保存在 [old](file:///d:/Linux/docker/bailangvvkruner/halo/old)
- Go 重构建议书位于 [go-refactor-proposal.md](file:///d:/Linux/docker/bailangvvkruner/halo/old/docs/go-refactor-proposal.md)
- 当前仓库根目录开始承载新的 Go 后端与新前端实现
- 当前已具备文章、页面、分类、标签、菜单、主题、插件、用户、设置、备份等基础数据链路
- 当前已具备完整的初始化认证链路、后台管理壳、50+ 后端 API、纯 Go SQLite 数据库、多语言初始化页


## 当前功能

- Go 后端：Gin + GORM + SQLite
- 新前端：React + Vite
- 工作目录：自动创建 `db`、`logs`、`themes`、`plugins`、`attachments`、`backups`
- 基础接口：健康检查、登录、文章、页面、分类、标签、菜单、评论、主题、插件、附件、备份、设置、用户
- 初始化数据：默认主题、默认插件、默认设置、示例文章、示例页面、示例分类、示例标签、示例菜单、示例备份
### 技术栈

- 后端：Go 1.25 + Gin + GORM + 纯 Go SQLite
- 前端：React 19 + Vite 7 + TypeScript + Axios
- 构建：Go 交叉编译 + Node 多阶段 Docker
- 数据库：SQLite（开发/测试），可切换 PostgreSQL/MySQL

## 本地运行

### 后端

```bash
go run ./cmd/halo
```

默认监听地址：`http://localhost:8090`

### 前端

```bash
cd web
npm install
npm run dev
```

开发前端地址：`http://localhost:3000`

在 Windows PowerShell 下如果遇到脚本执行限制，可使用：

```powershell
& npm.cmd install
& .\node_modules\.bin\vite.cmd
```

## 已验证

- `go build ./cmd/halo`
- `web` 前端生产构建

当前环境未安装 Docker，因此尚未完成镜像实机构建验证，但 [Dockerfile](file:///d:/Linux/docker/bailangvvkruner/halo/Dockerfile) 已按后端和前端多阶段构建方式配置。

## Docker 运行

```bash
docker build -t halo-go .
docker run -d --name halo-go -p 8090:8090 -v ~/.halo2:/root/.halo2 halo-go
```
