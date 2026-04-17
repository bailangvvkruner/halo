# Halo Go

这是一个基于 Go 和现代前端重构中的 Halo 实验版本。

## 当前状态

- 旧项目已整体保存在 [old](file:///d:/Linux/docker/bailangvvkruner/halo/old)
- Go 重构建议书位于 [go-refactor-proposal.md](file:///d:/Linux/docker/bailangvvkruner/halo/old/docs/go-refactor-proposal.md)
- 当前仓库根目录开始承载新的 Go 后端与新前端实现
- 当前已具备文章、页面、分类、标签、菜单、主题、插件、用户、设置、备份等基础数据链路
- 当前已具备 Go 后端编译、前端生产构建与 Dockerfile 多阶段构建定义

## 当前功能

- Go 后端：Gin + GORM + SQLite
- 新前端：React + Vite
- 工作目录：自动创建 `db`、`logs`、`themes`、`plugins`、`attachments`、`backups`
- 基础接口：健康检查、登录、文章、页面、分类、标签、菜单、评论、主题、插件、附件、备份、设置、用户
- 初始化数据：默认主题、默认插件、默认设置、示例文章、示例页面、示例分类、示例标签、示例菜单、示例备份

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
