# Coze 心理助手 API 文档 - 前端使用指南

## 一、API 概览

所有 API 都需要用户登录认证，请在请求中携带 Cookie 或 Token。

**Base URL**: `/api/v1/ai`

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/conversations` | 获取会话列表 |
| GET | `/messages` | 获取会话消息 |
| POST | `/chat` | 发送消息（流式响应） |
| POST | `/conversation/del` | 删除会话 |

---

## 二、API 详细说明

### 2.1 获取会话列表

**请求**
```
GET /api/v1/ai/conversations
```

**响应**
```json
{
  "code": 200,
  "msg": "ok",
  "data": [
    {
      "id": 1,
      "userId": 123,
      "conversationId": "7607315027620479011",
      "title": "你好，我想聊聊...",
      "createdAt": 1771216055,
      "updatedAt": 1771216055
    }
  ]
}
```

**前端使用示例**
```javascript
async function loadConversations() {
  const res = await fetch('/api/v1/ai/conversations', {
    credentials: 'include'  // 携带 Cookie
  });
  const data = await res.json();
  if (data.code === 200) {
    // 渲染会话列表
    renderConversationList(data.data);
  }
}
```

---

### 2.2 获取会话消息

**请求**
```
GET /api/v1/ai/messages?conversationId=7607315027620479011
```

**参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| conversationId | string | 是 | 会话ID |

**响应**
```json
{
  "code": 200,
  "msg": "ok",
  "data": [
    {
      "role": "user",
      "content": "你好",
      "createdAt": 1771216055
    },
    {
      "role": "assistant",
      "content": "你好呀，今天过得怎么样？",
      "createdAt": 1771216056
    }
  ]
}
```

**前端使用示例**
```javascript
async function loadMessages(conversationId) {
  const res = await fetch(`/api/v1/ai/messages?conversationId=${conversationId}`, {
    credentials: 'include'
  });
  const data = await res.json();
  if (data.code === 200) {
    // 渲染消息列表
    data.data.forEach(msg => {
      addMessageToUI(msg.role, msg.content);
    });
  }
}
```

---

### 2.3 发送消息（流式响应）

**请求**
```
POST /api/v1/ai/chat
Content-Type: application/json

{
  "conversationId": "7607315027620479011",  // 可选，新对话不传
  "message": "你好"
}
```

**参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| conversationId | string | 否 | 会话ID，不传则创建新对话 |
| message | string | 是 | 用户消息 |

**响应格式**：SSE (Server-Sent Events)

**事件类型**

| type | 说明 | 数据 |
|------|------|------|
| `conversation` | 新会话创建 | `{ conversationId: "xxx" }` |
| `content` | AI 回复内容（增量） | `{ content: "你" }` |
| `done` | 回复完成 | 无 |
| `error` | 发生错误 | `{ error: "错误信息" }` |

**前端使用示例**
```javascript
async function sendMessage(conversationId, message) {
  const res = await fetch('/api/v1/ai/chat', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify({
      conversationId: conversationId || '',  // 空字符串表示新对话
      message: message
    })
  });

  // 读取 SSE 流
  const reader = res.body.getReader();
  const decoder = new TextDecoder();
  let fullContent = '';

  while (true) {
    const { done, value } = await reader.read();
    if (done) break;

    const text = decoder.decode(value);
    const lines = text.split('\n');

    for (const line of lines) {
      if (line.startsWith('data:')) {
        try {
          const json = JSON.parse(line.substring(5).trim());
          
          switch (json.type) {
            case 'content':
              // 逐字追加内容
              fullContent += json.content;
              updateAssistantMessage(fullContent);
              break;
              
            case 'conversation':
              // 新会话创建，保存 conversationId
              currentConversationId = json.conversationId;
              break;
              
            case 'done':
              // 回复完成
              console.log('回复完成');
              break;
              
            case 'error':
              // 错误处理
              console.error('错误:', json.error);
              break;
          }
        } catch (e) {
          // 解析错误，忽略
        }
      }
    }
  }
}
```

---

### 2.4 删除会话

**请求**
```
POST /api/v1/ai/conversation/del
Content-Type: application/json

{
  "conversationId": "7607315027620479011"
}
```

**响应**
```json
{
  "code": 200,
  "msg": "删除成功喵",
  "data": null
}
```

**前端使用示例**
```javascript
async function deleteConversation(conversationId) {
  const res = await fetch('/api/v1/ai/conversation/del', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify({ conversationId })
  });
  const data = await res.json();
  if (data.code === 200) {
    // 从列表中移除
    removeConversationFromUI(conversationId);
  }
}
```

---

## 三、完整前端示例

### 3.1 HTML 结构

```html
<div class="chat-container">
  <!-- 会话列表 -->
  <div class="sidebar">
    <button onclick="newChat()">新对话</button>
    <div id="conversationList"></div>
  </div>
  
  <!-- 聊天区域 -->
  <div class="main-chat">
    <div id="messages"></div>
    <textarea id="input" placeholder="输入消息..."></textarea>
    <button onclick="send()">发送</button>
  </div>
</div>
```

### 3.2 JavaScript 完整代码

```javascript
const API_BASE = '/api/v1/ai';
let currentConversationId = '';
let isStreaming = false;

// 初始化
document.addEventListener('DOMContentLoaded', () => {
  loadConversations();
});

// 加载会话列表
async function loadConversations() {
  const res = await fetch(`${API_BASE}/conversations`, {
    credentials: 'include'
  });
  const data = await res.json();
  if (data.code === 200) {
    renderConversations(data.data);
  }
}

// 渲染会话列表
function renderConversations(conversations) {
  const container = document.getElementById('conversationList');
  container.innerHTML = conversations.map(conv => `
    <div class="conversation-item" 
         data-id="${conv.conversationId}"
         onclick="selectConversation('${conv.conversationId}')">
      ${conv.title || '新对话'}
    </div>
  `).join('');
}

// 选择会话
async function selectConversation(id) {
  currentConversationId = id;
  document.getElementById('messages').innerHTML = '';
  await loadMessages(id);
}

// 加载消息
async function loadMessages(conversationId) {
  const res = await fetch(`${API_BASE}/messages?conversationId=${conversationId}`, {
    credentials: 'include'
  });
  const data = await res.json();
  if (data.code === 200) {
    data.data.forEach(msg => {
      addMessage(msg.role, msg.content);
    });
  }
}

// 新对话
function newChat() {
  currentConversationId = '';
  document.getElementById('messages').innerHTML = '';
}

// 发送消息
async function send() {
  const input = document.getElementById('input');
  const message = input.value.trim();
  if (!message || isStreaming) return;
  
  isStreaming = true;
  input.value = '';
  
  addMessage('user', message);
  const assistantMsg = addMessage('assistant', '');
  
  try {
    const res = await fetch(`${API_BASE}/chat`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({
        conversationId: currentConversationId,
        message: message
      })
    });

    const reader = res.body.getReader();
    const decoder = new TextDecoder();
    let fullContent = '';

    while (true) {
      const { done, value } = await reader.read();
      if (done) break;

      const text = decoder.decode(value);
      const lines = text.split('\n');

      for (const line of lines) {
        if (line.startsWith('data:')) {
          try {
            const json = JSON.parse(line.substring(5).trim());
            
            if (json.type === 'content') {
              fullContent += json.content;
              assistantMsg.textContent = fullContent;
            } else if (json.type === 'conversation') {
              currentConversationId = json.conversationId;
              loadConversations();  // 刷新列表
            } else if (json.type === 'error') {
              assistantMsg.textContent = '错误: ' + json.error;
            }
          } catch (e) {}
        }
      }
    }
  } catch (err) {
    assistantMsg.textContent = '网络错误: ' + err.message;
  } finally {
    isStreaming = false;
  }
}

// 添加消息到界面
function addMessage(role, content) {
  const container = document.getElementById('messages');
  const div = document.createElement('div');
  div.className = `message ${role}`;
  div.textContent = content;
  container.appendChild(div);
  return div;
}

// 回车发送
document.getElementById('input').addEventListener('keydown', (e) => {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault();
    send();
  }
});
```

---

## 四、注意事项

### 4.1 认证要求

所有 API 都需要用户登录，未登录会返回：
```json
{
  "code": 233,
  "msg": "未登录或登录已过期喵"
}
```

### 4.2 流式响应处理

1. 使用 `ReadableStream` 读取响应体
2. 每行以 `data:` 开头，后面是 JSON
3. `content` 事件是增量的，需要累加

### 4.3 会话 ID 管理

- 新对话：不传 `conversationId` 或传空字符串
- 继续对话：传入之前返回的 `conversationId`
- 首次对话会返回 `conversation` 事件，需要保存

### 4.4 错误处理

| code | 说明 |
|------|------|
| 200 | 成功 |
| 233 | 业务已知错误 |
| 500 | 服务器内部错误 |
