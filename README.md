# feedsystem_video_go

## 项目简介
基于 Go 的视频 Feed 系统后端，提供账号注册/登录、视频发布与查询、点赞/取消点赞、Feed 拉取等接口。默认使用 `Gin + GORM + MySQL + JWT`。

仓库内附 Postman Collection（见 `test/*.json`）方便调试 API。

## 目录结构
- `cmd/`：程序入口（`cmd/main.go`）
- `configs/`：YAML 配置（`configs/config.yaml`）
- `internal/account/`：账号模块（实体/仓储/服务/HTTP handler）
- `internal/video/`：视频与点赞模块
- `internal/feed/`：Feed 拉取模块
- `internal/http/`：Gin 路由注册
- `internal/middleware/`：JWT 中间件
- `test/`：Postman Collection

## 快速开始
1. 准备 MySQL，并创建数据库（默认名 `feedsystem`，可在 `configs/config.yaml` 修改）。
2. 安装依赖：`go mod tidy`
3. 启动服务：`go run ./cmd`（端口默认 `8080`）
4. 导入 Postman：`test/*.json`，设置变量 `host=http://localhost:8080`、`jwt_token` 为登录返回的 token。

## 配置
配置文件：`configs/config.yaml`

```yaml
server:
  port: 8080

database:
  host: localhost
  port: 3306
  user: root
  password: 123456
  dbname: feedsystem
```

## 认证说明（与代码一致）
- Header：`Authorization: Bearer <jwt>`
- JWT 解析成功后，还会校验该 token 是否等于数据库里账号当前保存的 token（见 `internal/middleware/jwt.go`），因此：
  - 同一账号再次登录会覆盖 token，旧 token 失效
  - `/account/logout` 会清空 token，token 立即失效
  - `/account/changePassword` 成功后也会清空 token（需要重新登录拿新 token）

## API（路由与鉴权）
路由注册见 `internal/http/router.go`。

### 账号（/account）
| 方法 | 路径 | 是否需要 JWT | 功能 | 请求体示例 |
|------|------|-------------|------|-----------|
| POST | `/account/register` | 否 | 注册账号 | `{"username":"alice","password":"pass123"}` |
| POST | `/account/login` | 否 | 登录，返回 `token` | `{"username":"alice","password":"pass123"}` |
| POST | `/account/changePassword` | 否 | 修改指定用户名密码（会校验旧密码） | `{"username":"alice","old_password":"pass123","new_password":"newpass456"}` |
| POST | `/account/findByID` | 否 | 按 ID 查询账号 | `{"id":1}` |
| POST | `/account/findByUsername` | 否 | 按用户名查询账号 | `{"username":"alice"}` |
| POST | `/account/rename` | 是 | 修改当前登录用户用户名 | `{"new_username":"alice_new"}` |
| POST | `/account/logout` | 是 | 注销（撤销当前 token） | `{}` |

### 视频（/video）
| 方法 | 路径 | 是否需要 JWT | 功能 | 请求体示例 |
|------|------|-------------|------|-----------|
| POST | `/video/listByAuthorID` | 否 | 按作者 ID 列出视频（不分页） | `{"author_id":1}` |
| POST | `/video/getDetail` | 否 | 查询视频详情 | `{"id":1}` |
| POST | `/video/publish` | 是 | 发布视频（作者取自 JWT） | `{"title":"demo","description":"...","play_url":"http://...","cover_url":"http://..."}` |

说明：发布视频时 `title`、`play_url`、`cover_url` 为必填字段（见 `internal/video/video_service.go`）。

### 点赞（/like）
| 方法 | 路径 | 是否需要 JWT | 功能 | 请求体示例 |
|------|------|-------------|------|-----------|
| POST | `/like/getLikesCount` | 否 | 获取视频点赞总数 | `{"video_id":1}` |
| POST | `/like/isLiked` | 是 | 当前用户是否点赞 | `{"video_id":1}` |
| POST | `/like/like` | 是 | 点赞 | `{"video_id":1}` |
| POST | `/like/unlike` | 是 | 取消点赞 | `{"video_id":1}` |

### Feed（/feed）
| 方法 | 路径 | 是否需要 JWT | 功能 | 请求体示例 |
|------|------|-------------|------|-----------|
| POST | `/feed/listLatest` | 否 | 拉取最新 Feed（按 `create_time` 倒序），`limit` 最大 50 | `{"limit":10,"latest_time":0}` |

说明：
- `latest_time` 是 Unix 秒时间戳；当 `latest_time=0` 或不传时，返回最新一页（见 `internal/feed/handler.go`）。
- 响应里的 `next_time` 是 JSON 序列化后的 `time.Time`（RFC3339 字符串，见 `internal/feed/service.go`）。下一页请求时，需要把它转换为 Unix 秒时间戳再传回 `latest_time`。

## 数据表（自动迁移）
启动时会执行 `AutoMigrate` 创建/更新表结构：`Account`、`Video`、`Like`（见 `internal/db/db.go`）。
