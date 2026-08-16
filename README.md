# 论坛系统（Forum）

一个纯 Go 标准库实现的论坛 REST API 服务，采用标准 Go 工程目录结构，内存存储，零第三方依赖。

## 目录结构

```
origin/
├── cmd/server/          # 程序入口
├── internal/
│   ├── app/             # 依赖装配
│   ├── config/          # 配置加载
│   ├── model/           # 领域模型与校验
│   ├── store/           # 数据访问接口 + 内存实现
│   ├── service/         # 业务逻辑层
│   └── handler/         # HTTP 处理器层
└── pkg/
    ├── httpx/           # HTTP 响应工具
    ├── idgen/           # ID / 短码生成
    └── logger/          # 分级日志
```

## 运行

```bash
go run ./cmd/server
PORT=8081 go run ./cmd/server
```

默认监听 `:8080`。

## 测试

```bash
go test ./...
```

## API 接口

### 用户

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/users` | 注册用户 |
| GET | `/api/users` | 用户列表 |
| GET | `/api/users/{id}` | 用户详情 |
| PUT | `/api/users/{id}` | 更新资料 |

### 版块

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/boards` | 创建版块 |
| GET | `/api/boards` | 版块列表 |
| GET | `/api/boards/{id}` | 版块详情 |
| PUT | `/api/boards/{id}` | 更新版块 |
| DELETE | `/api/boards/{id}` | 删除版块 |

### 主题

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/topics` | 发帖 |
| GET | `/api/topics?board_id=&keyword=&featured=&page=&size=` | 主题列表（置顶优先） |
| GET | `/api/topics/{id}` | 主题详情（累加浏览数） |
| PUT | `/api/topics/{id}` | 编辑主题 |
| DELETE | `/api/topics/{id}` | 删除主题 |
| PATCH | `/api/topics/{id}/flag` | 更新标记 `{"pinned":true,"featured":false,"closed":false}` |
| POST | `/api/topics/{id}/like` | 点赞 `{"user_id":"u1"}` |
| DELETE | `/api/topics/{id}/like?user_id=u1` | 取消点赞 |

### 回帖

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/topics/{id}/replies` | 回帖 |
| GET | `/api/topics/{id}/replies` | 回帖列表 |
| DELETE | `/api/replies/{id}` | 删除回帖 |

### 统计

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/stats` | 全局统计 |
| GET | `/api/stats/boards` | 按版块统计 |
| GET | `/api/stats/hot?limit=10` | 热门主题 |
