# feedsystem_video_go

基于 Go 的视频 Feed 系统后端，提供账号、视频、点赞、评论、Feed 等接口。默认技术栈：`Gin + GORM + MySQL + JWT`。

仓库内附 Postman Collection（`test/*.json`）用于手动/批量调试接口，推荐使用一体化集合 `test/postman.json`。

## 目录结构
- `cmd/`：程序入口（`cmd/main.go`）
- `configs/`：YAML 配置（`configs/config.yaml`）
- `internal/account/`：账号模块
- `internal/video/`：视频/点赞/评论模块
- `internal/feed/`：Feed 模块
- `internal/http/`：Gin 路由注册（`internal/http/router.go`）
- `internal/middleware/`：JWT 中间件
- `test/`：Postman Collection（推荐 `test/postman.json`）

## 快速开始
1. 准备 MySQL 并创建数据库（库名/账号在 `configs/config.yaml` 配置）。
2. 安装依赖：`go mod tidy`
3. 启动服务：`go run ./cmd`
4. Postman 导入：`test/postman.json`，默认 `host=http://localhost:8080`。

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

可选环境变量：
- `JWT_SECRET`：JWT 签名密钥；不设置则使用默认值（仅建议本地调试）。

## 认证说明（与代码一致）
- 认证 Header：`Authorization: Bearer <jwt>`
- JWT 通过后会额外校验：请求 token 必须等于数据库里该账号当前保存的 token（`internal/middleware/jwt.go`），因此：
  - 同一账号再次登录会覆盖 token（旧 token 失效）
  - `/account/logout` 会清空 token（立即失效）
  - `/account/changePassword` 成功后会清空 token（需要重新登录）
  - `/account/rename` 成功后会返回**新 token**并写回数据库（旧 token 立即失效；客户端必须替换保存）
- Feed 接口使用“软鉴权”：可以不带 token；但如果带了 `Authorization`，必须是合法且未撤销的 token，否则会返回 `401`。

## Postman 完整测试流程（建议）
使用一体化集合：`test/postman.json`（含预置变量与自动保存脚本）。

建议运行顺序：
1. Account → Register Account
2. Account → Login (save jwt_token)（自动保存 `jwt_token`）
3. Account → Find By Username (save accountId)（自动保存 `accountId`）
4. Video → Publish Video (save videoId)（自动保存 `videoId`/`authorId`）
5. Like（可选）
6. Comment（可选）
7. Feed（可选；翻页变量会自动更新）

注意：执行 `Account/Rename` 后，集合会把响应里的 `token` 覆盖到 `jwt_token`，否则后续鉴权接口会因为旧 token 失效而 `401`。

## API（路由与鉴权）
路由注册见 `internal/http/router.go`。以下均为 `POST` + JSON body。

### 账号（/account）
| 路径 | 是否需要 JWT | 说明 |
|------|-------------|------|
| `/account/register` | 否 | 注册：`{"username":"alice","password":"pass123"}` |
| `/account/login` | 否 | 登录：`{"username":"alice","password":"pass123"}` → `{"token":"..."}` |
| `/account/changePassword` | 否 | 修改密码：`{"username":"alice","old_password":"pass123","new_password":"newpass456"}`（成功会登出） |
| `/account/findByID` | 否 | `{"id":1}` |
| `/account/findByUsername` | 否 | `{"username":"alice"}` |
| `/account/rename` | 是 | `{"new_username":"alice_new"}` → `{"token":"..."}`（返回新 token） |
| `/account/logout` | 是 | `{}` |

### 视频（/video）
| 路径 | 是否需要 JWT | 说明 |
|------|-------------|------|
| `/video/listByAuthorID` | 否 | `{"author_id":1}` |
| `/video/getDetail` | 否 | `{"id":1}` |
| `/video/publish` | 是 | `{"title":"demo","description":"...","play_url":"http://...","cover_url":"http://..."}`（必填：`title/play_url/cover_url`） |

### 点赞（/like）
| 路径 | 是否需要 JWT | 说明 |
|------|-------------|------|
| `/like/getLikesCount` | 否 | `{"video_id":1}` |
| `/like/isLiked` | 是 | `{"video_id":1}` |
| `/like/like` | 是 | `{"video_id":1}` |
| `/like/unlike` | 是 | `{"video_id":1}` |

### 评论（/comment）
| 路径 | 是否需要 JWT | 说明 |
|------|-------------|------|
| `/comment/listAll` | 否 | `{"video_id":1}` |
| `/comment/publish` | 是 | `{"video_id":1,"content":"hello"}` |
| `/comment/delete` | 是 | `{"comment_id":1}`（仅作者可删） |

### Feed（/feed）
| 路径 | 是否需要 JWT | 说明 |
|------|-------------|------|
| `/feed/listLatest` | 否（可选 JWT） | `{"limit":10,"latest_time":0}` |
| `/feed/listLikesCount` | 否（可选 JWT） | `{"limit":10,"likes_count":0}` |

翻页说明：
- `/feed/listLatest`：`latest_time` 为 Unix 秒时间戳；响应 `next_time` 也是 Unix 秒（`0` 表示没有下一页/无数据）。
- `/feed/listLikesCount`：`likes_count` 表示“上一页最后一条的点赞数”；响应 `next_likes_count_before` 用于下一页请求。

## 数据表（自动迁移）
启动时会执行 `AutoMigrate` 创建/更新表结构（`internal/db/db.go`）：`Account`、`Video`、`Like`、`Comment`。

