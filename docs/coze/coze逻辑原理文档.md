# Coze 心理助手集成 - 逻辑原理文档

## 一、整体架构

```
┌─────────────────────────────────────────────────────────────────┐
│                        前端 Web 应用                             │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                      后端 API 服务                               │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │ Controller  │→ │  Service    │→ │   Model     │              │
│  │ (路由处理)   │  │ (业务逻辑)   │  │ (数据存储)   │              │
│  └─────────────┘  └─────────────┘  └─────────────┘              │
│         │                │                                       │
│         ▼                ▼                                       │
│  ┌─────────────────────────────┐                                 │
│  │      CozeClient (utils)     │                                 │
│  │      Coze API 封装          │                                 │
│  └─────────────────────────────┘                                 │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Coze API 服务                               │
│              https://api.coze.cn/v3/chat                         │
└─────────────────────────────────────────────────────────────────┘
```

## 二、核心模块说明

### 2.1 数据模型 (models/AIchat.go)

```go
type AIConversation struct {
    Id             int    // 主键
    UserId         int    // 关联用户ID
    ConversationId string // Coze 返回的会话ID
    Title          string // 会话标题（取首条消息前30字）
    CreatedAt      int    // 创建时间
    UpdatedAt      int    // 更新时间
}
```

**设计思路**：
- 本地只存储会话元数据（ID、标题、时间）
- 实际对话内容由 Coze 平台托管
- 通过 `ConversationId` 关联 Coze 会话

### 2.2 Coze 客户端 (utils/coze.go)

#### 核心结构

```go
type CozeClient struct {
    AccessToken string      // SAT 令牌
    BaseURL     string      // API 地址
    BotID       string      // 机器人ID
    HTTPClient  *http.Client
}
```

#### 流式对话原理

```
客户端请求 ──→ POST /v3/chat?conversation_id=xxx
                    │
                    ▼
            SSE 流式响应
                    │
    ┌───────────────┼───────────────┐
    │               │               │
    ▼               ▼               ▼
event:conversation  event:message   event:done
.chat.created       .delta
    │               │
    │               ▼
    │         逐字返回内容
    ▼
返回会话ID
```

**关键点**：
1. `conversation_id` 必须作为 **URL 查询参数** 传递
2. 响应格式为 SSE (Server-Sent Events)
3. 每行格式：`event:事件名` + `data:JSON数据`

#### 事件类型处理

| 事件 | 说明 | 处理方式 |
|------|------|----------|
| `conversation.chat.created` | 会话创建 | 提取 conversation_id |
| `conversation.message.delta` | 消息增量 | 提取 content 逐字显示 |
| `conversation.message.completed` | 消息完成 | 标记对话结束 |
| `conversation.chat.completed` | 对话完成 | 无需处理 |
| `error` | 错误 | 返回错误信息 |

### 2.3 业务服务 (service/aichatService.go)

#### 会话管理流程

```
新对话流程：
用户发送消息 → Coze 创建会话 → 返回 conversation_id → 保存到本地数据库

继续对话流程：
用户发送消息 → 携带 conversation_id → Coze 继续该会话 → 更新本地时间戳
```

#### 核心函数

| 函数 | 功能 |
|------|------|
| `GetCozeClient()` | 获取 Coze 客户端（单例） |
| `GetUserConversations()` | 获取用户会话列表 |
| `CreateConversation()` | 创建新会话记录 |
| `UpdateConversationTime()` | 更新会话时间 |
| `DeleteConversation()` | 删除会话 |
| `GenerateTitle()` | 生成会话标题 |

### 2.4 控制器 (controller/aichatController.go)

#### 流式响应处理

```go
func ChatStream(c *gin.Context) {
    // 1. 设置 SSE 响应头
    c.Header("Content-Type", "text/event-stream")
    c.Header("Cache-Control", "no-cache")
    
    // 2. 调用 Coze 流式接口
    cozeClient.ChatStream(req, func(event, data string) error {
        // 3. 根据事件类型推送 SSE
        c.SSEvent("message", gin.H{...})
        c.Writer.Flush()  // 立即发送
    })
    
    // 4. 保存/更新会话记录
    if isNewConversation {
        service.CreateConversation(...)
    } else {
        service.UpdateConversationTime(...)
    }
}
```

## 三、数据流向

### 3.1 发送消息流程

```
1. 前端 POST /api/v1/ai/chat
   {
     "conversationId": "xxx",  // 可选，新对话为空
     "message": "你好"
   }
   
2. 后端验证用户 → 调用 Coze API
   
3. Coze 返回 SSE 流：
   event: conversation.chat.created
   data: {"conversation_id":"7607xxx",...}
   
   event: conversation.message.delta
   data: {"content":"你",...}
   
   event: conversation.message.delta
   data: {"content":"好",...}
   
   event: conversation.message.completed
   data: {...}
   
4. 后端转发 SSE 给前端：
   data: {"type":"conversation","conversationId":"7607xxx"}
   data: {"type":"content","content":"你"}
   data: {"type":"content","content":"好"}
   data: {"type":"done"}
   
5. 后端保存会话记录到数据库
```

### 3.2 获取历史消息流程

```
1. 前端 GET /api/v1/ai/messages?conversationId=xxx

2. 后端验证会话归属 → 调用 Coze 消息列表 API

3. Coze 返回消息列表：
   {
     "code": 0,
     "data": [
       {"role":"assistant","type":"answer","content":"你好呀",...},
       {"role":"user","type":"question","content":"你好",...}
     ]
   }

4. 后端过滤并反转顺序 → 返回给前端
```

## 四、配置说明

```toml
[coze]
access_token = 'sat_xxx'  # 服务访问令牌（永久有效）
bot_id = '7548xxx'        # 机器人ID
```

**获取方式**：
1. 登录 https://coze.cn
2. 进入「扣子 API」→「授权」
3. 创建服务访问令牌（SAT）
4. Bot ID 在智能体开发页面 URL 中获取

## 五、安全考虑

1. **用户隔离**：每个用户只能访问自己的会话
2. **会话验证**：操作会话前验证归属权
3. **令牌安全**：SAT 存储在配置文件，不暴露给前端
4. **权限控制**：所有 API 需要登录认证
