# LinkChecker

LinkChecker 是一个 Go 项目，用于检测网盘分享链接的有效性与状态，并可选择将检测结果保存/上报。

它提供：

- CLI：支持命令行参数检测、CSV 批量检测
- Server：Gin API 服务 + 内嵌前端静态资源（`web/dist`）
- CLI -> Server 上报：用于落库、列表查询

当前内置支持的网盘：

- 百度网盘（`pan.baidu.com`）
- 夸克网盘（`pan.quark.cn`）

## 环境要求

- Go（见 `go.mod`，当前为 `go 1.24.6`）

可选：

- Task（用于运行 `Taskfile.yml`）
- UPX（用于 `task upx` 压缩二进制）

## 快速开始：CLI

直接检测（命令行模式）：

```bash
go run . check "https://pan.quark.cn/s/xxxx" "https://pan.baidu.com/s/xxxx?pwd=abcd"
```

CSV 批处理（文件模式）：

```bash
go run . check -f ./input.csv -o ./result.csv
```

CSV 格式：

- 每行：`链接,密码`
- 密码可选；没有就留空

### CLI 参数说明

`check` 命令参数（见 `cmd/check.go`）：

- `-f, --file` 输入 CSV 文件路径（格式：`链接,密码`）
- `-o, --output` 输出 CSV 文件路径（仅文件模式生效，默认 `./result.csv`）
- `-s, --server` 服务端上报地址（例如 `http://127.0.0.1:8080`）
- `-t, --token` 服务端鉴权 Token（作为请求头 `token` 发送）
- `-x, --proxy` HTTP 代理（例如 `http://127.0.0.1:7890`）
- `-p, --parallel` 并发检测数量（自动限制为 `CPU核数 * 2` 上限）

补充：

- 检测结果在进程内使用内存缓存 10 分钟（减少重复请求）。
- 当 `--server` 非空且检测结果为“有效”时，会向服务端 `POST /api/v1/report` 上报。

## 快速开始：Server

启动 API 服务：

```bash
go run . server
```

### 配置文件

服务端通过 Viper 加载配置（见 `configs/func.go`）：

- 优先读取 `./config.yaml`
- 其次读取 `./data/config.yaml`

如果配置不存在，会自动生成默认配置到 `data/config.yaml`，并确保以下目录存在：

- `data/`（默认 SQLite 数据库文件所在目录）
- `logs/`（默认日志目录）

默认配置要点（见 `configs/config.go`）：

- Server 端口：`8080`
- Server Token：随机 16 字符（用于上报接口鉴权）
- 数据库：`sqlite3`，路径 `data/data.db`
- 日志：`logs/app.log`

注意：目前实现中（见 `cmd/server.go`），配置文件中的 `server.port` 会覆盖 CLI 的 `--port` 参数。

### API 列表

健康检查：

- `GET /ping` -> `pong`

链接接口（按 IP 限流）：

- `POST /api/v1/link/` 检测单个链接并保存
- `POST /api/v1/link/list` 分页查询保存的链接（支持 keyword）

上报接口（需要 token）：

- `POST /api/v1/report` 保存上报结果（CLI 在 `--server` 模式会调用）

鉴权方式：

- 请求头：`token: <server.token>`
- 缺失或不匹配会返回 400。

### Web 页面

说明：本仓库默认不会公开前端构建产物，`web/dist` 已加入 `.gitignore`。

如果你在私有环境自行构建了前端并放置到 `web/dist`，服务端会将其作为静态资源内嵌（见 `web/web.go`），并对外提供：

- `/`（SPA 入口）
- `/assets/*`、`/js/*`、`/favicon.ico`

除 `/api` 以外的未知路由会回落到 `index.html`，用于 SPA 前端路由。

## 构建与运行

使用 Taskfile（推荐）：

```bash
task build
```

运行测试：

```bash
task test
```

本地运行（会将输出追加到临时日志文件）：

```bash
task run
```

不使用 Taskfile 的普通构建：

```bash
go build -o bin/linkchecker .
```

## 版本信息

查看版本/构建信息：

```bash
go run . version
```

如果使用 `Taskfile.yml` 构建，会通过 `-ldflags` 注入（见 `internal/conf/var.go`）：

- `internal/conf.Version`
- `internal/conf.BuildAt`
- `internal/conf.GitAuthor`
- `internal/conf.GitEmail`
- `internal/conf.GitCommit`

## 目录结构

- `cmd/` Cobra CLI 命令（`check`/`server`/`version`）
- `internal/netdisk/` 各网盘 Provider + Manager
- `internal/router/` Gin 路由 + 前端静态资源挂载
- `internal/service/`、`internal/repo/`、`internal/model/` 业务与持久化（GORM）
- `configs/` 配置结构体 + 加载/默认生成
- `pkg/` 公共包（db/log/request/cache/response 等）
- `web/` 内嵌前端构建产物（`web/dist`）

## License

MIT（见 `LICENSE`）。

## GitHub 自动发布

仓库已支持基于 GitHub Actions 的自动打标签发布：

使用前请先在仓库 `Settings -> Actions -> General -> Workflow permissions` 中启用 `Read and write permissions`，否则 workflow 无法推送 tag 或创建 release。

- 在 GitHub Actions 中手动运行 `Create Release Tag`
- 输入版本号，例如 `v1.2.3`
- Workflow 会创建并推送对应 tag
- tag 推送后会自动触发 `Release`
- `Release` 会执行测试、构建多平台二进制，并创建 GitHub Release

当前发布产物包含：

- Linux：`amd64`、`arm64`、`386`、`armv7`
- Windows：`amd64`、`386`、`arm64`
- macOS：`amd64`、`arm64`
