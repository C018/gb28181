# Frontend Integration Guide (前端集成指南)

本文档说明后端 API 已准备就绪，等待前端集成的功能。

## 概述

后端已完成以下功能的 API 实现，现需要前端界面支持：

## 1. GB28181 设备 ID 输入验证

### 后端状态
✅ 已完成

### 后端功能
- 自动验证 GB28181 设备 ID 必须为 20 位
- 验证必须为纯数字格式
- 提供清晰的错误提示

### 前端需求
**输入框增强**
```javascript
// 示例实现
const [deviceId, setDeviceId] = useState('');
const maxLength = 20;

const handleInput = (value) => {
  // 只允许数字
  const numericValue = value.replace(/[^0-9]/g, '');
  setDeviceId(numericValue);
};

// UI 显示
<input 
  value={deviceId}
  onChange={(e) => handleInput(e.target.value)}
  maxLength={maxLength}
  placeholder="请输入20位设备ID"
/>
<span className="counter">{deviceId.length}/{maxLength} 位</span>
{deviceId.length !== 0 && deviceId.length !== maxLength && (
  <span className="error">设备ID必须为20位</span>
)}
```

### API 端点
```
POST /devices
{
  "type": "GB28181",
  "device_id": "34020000001320000001"  // 必须 20 位数字
}

// 错误响应示例
{
  "code": 400,
  "msg": "GB28181 设备 ID 必须为 20 位（当前: 13 位），少或多都可能导致信令故障"
}
```

---

## 2. 登录功能

### 后端状态
✅ 已完成

### 后端功能
- 基于配置文件的账号密码认证
- JWT token 认证（有效期 3 天）
- 所有敏感 API 需要 token 认证

### 前端需求
**登录页面**
- 用户名输入框
- 密码输入框（密码类型）
- 登录按钮
- 记住登录状态（token 存储在 localStorage）
- 显示登录错误信息

**示例实现**
```javascript
const login = async (username, password) => {
  try {
    const response = await fetch('/user/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password })
    });
    
    if (!response.ok) {
      throw new Error('登录失败');
    }
    
    const data = await response.json();
    // 存储 token
    localStorage.setItem('token', data.token);
    localStorage.setItem('username', data.user);
    
    // 跳转到主页
    navigate('/dashboard');
  } catch (error) {
    showError(error.message);
  }
};
```

### API 端点
```
POST /user/login
Content-Type: application/json

请求体:
{
  "username": "admin",
  "password": "admin"
}

成功响应:
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": "admin"
}

错误响应:
{
  "code": 401,
  "msg": "用户名或密码错误"
}
```

---

## 3. 修改密码功能

### 后端状态
✅ 已完成

### 后端功能
- 支持修改用户名和密码
- 需要 JWT token 认证
- 更改后自动保存到配置文件
- 下次启动时生效

### 前端需求
**密码修改界面**
- 新用户名输入框
- 新密码输入框
- 确认密码输入框
- 提交按钮
- 成功/错误提示

**示例实现**
```javascript
const changePassword = async (newUsername, newPassword, confirmPassword) => {
  // 验证
  if (newPassword !== confirmPassword) {
    showError('两次密码输入不一致');
    return;
  }
  
  if (newPassword.length < 6) {
    showError('密码长度至少6位');
    return;
  }
  
  try {
    const token = localStorage.getItem('token');
    const response = await fetch('/user/user', {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify({
        username: newUsername,
        password: newPassword
      })
    });
    
    if (!response.ok) {
      throw new Error('修改失败');
    }
    
    const data = await response.json();
    showSuccess(data.msg);
    
    // 更新本地存储的用户名
    localStorage.setItem('username', newUsername);
  } catch (error) {
    showError(error.message);
  }
};
```

### API 端点
```
PUT /user/user
Authorization: Bearer <token>
Content-Type: application/json

请求体:
{
  "username": "newadmin",
  "password": "newpassword"
}

成功响应:
{
  "msg": "凭据更新成功"
}

错误响应:
{
  "code": 401,
  "msg": "token 无效或已过期"
}
```

---

## 4. API 请求通用处理

### Token 认证
所有需要认证的 API 都需要在请求头中携带 token：

```javascript
const apiRequest = async (url, options = {}) => {
  const token = localStorage.getItem('token');
  
  const headers = {
    'Content-Type': 'application/json',
    ...options.headers,
  };
  
  // 添加 token（除了登录接口）
  if (token && url !== '/user/login') {
    headers['Authorization'] = `Bearer ${token}`;
  }
  
  const response = await fetch(url, {
    ...options,
    headers
  });
  
  // 处理 401 未授权
  if (response.status === 401) {
    localStorage.removeItem('token');
    localStorage.removeItem('username');
    // 跳转到登录页
    window.location.href = '/login';
    throw new Error('登录已过期，请重新登录');
  }
  
  return response;
};
```

### 无需认证的公开端点
以下端点不需要 token：
- `/health` - 健康检查
- `/app/metrics/api` - 监控指标
- `/user/login` - 登录接口
- `/webhook/*` - 流媒体服务器回调（内部调用）

### 需要认证的端点
所有其他端点都需要 token 认证，包括：
- 设备管理 `/devices/*`
- 通道管理 `/channels/*`
- PTZ 控制 `/ptz/*`
- 录像回放 `/playback/*`
- 报警订阅 `/alarm/*`
- 配置管理 `/config/*`
- 流媒体代理 `/proxy/sms/*`
- 实时通知订阅 `/notifications/subscribe`

---

## 5. 实时消息通知 (SSE)

### 后端状态
✅ 已完成

### 前端需求
使用 Server-Sent Events (SSE) 接收实时通知

```javascript
const subscribeNotifications = () => {
  const token = localStorage.getItem('token');
  const eventSource = new EventSource(
    `/notifications/subscribe?token=${token}`
  );
  
  eventSource.onmessage = (event) => {
    const notification = JSON.parse(event.data);
    handleNotification(notification);
  };
  
  eventSource.onerror = (error) => {
    console.error('SSE 连接错误:', error);
    eventSource.close();
    // 可以实现重连逻辑
  };
  
  return eventSource;
};

const handleNotification = (notification) => {
  switch (notification.type) {
    case 'device':
      // 设备上线/离线
      updateDeviceStatus(notification.data);
      break;
    case 'stream':
      // 流开始/停止
      updateStreamStatus(notification.data);
      break;
    case 'alarm':
      // 报警事件
      showAlarmNotification(notification.data);
      break;
    case 'ai':
      // AI 检测告警
      showAIAlert(notification.data);
      break;
    case 'recording':
      // 录像开始/停止
      updateRecordingStatus(notification.data);
      break;
  }
};
```

---

## 6. 快照功能

### 后端状态
✅ 已完成（包括毛玻璃效果配置）

### 前端需求
显示通道快照

```javascript
// 获取快照链接
const getSnapshotUrl = (channelId) => {
  const token = localStorage.getItem('token');
  return `/channels/${channelId}/snapshot?token=${token}`;
};

// 刷新快照
const refreshSnapshot = async (channelId) => {
  const token = localStorage.getItem('token');
  const response = await fetch(`/channels/${channelId}/refresh-snapshot`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  const data = await response.json();
  return data.link; // 新的快照链接
};
```

**毛玻璃效果**
- 后端支持配置 `enable_snapshot_blur` 开关
- 前端只需正常显示快照图片
- 模糊效果由后端处理

---

## 7. 播放鉴权

### 后端状态
✅ 已完成

### 说明
- 播放链接包含时效性 token
- 有效期可配置（默认 10 分钟）
- token 格式：`timestamp.MD5(timestamp + secret)`

### 前端处理
```javascript
// 获取播放链接（后端返回的链接已包含 token）
const getPlayUrl = async (deviceId, channelId) => {
  const token = localStorage.getItem('token');
  const response = await fetch(`/channels/${channelId}/play`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({
      device_id: deviceId,
      channel_id: channelId
    })
  });
  const data = await response.json();
  return data.url; // 包含时效 token 的播放 URL
};
```

---

## 总结

### 立即可实现的前端功能

1. **登录页面** - 所有 API 已就绪
2. **密码修改页面** - 所有 API 已就绪  
3. **GB28181 设备添加** - 增加输入验证和字符计数
4. **实时通知** - SSE 订阅和消息处理
5. **快照显示** - 显示和刷新

### 配置说明

前端可以在用户设置中添加以下配置项的界面：
- 毛玻璃效果开关
- 播放链接有效期
- 其他系统配置

### API 文档

完整的 API 文档请参考：
- [API接口文档.md](./API接口文档.md)
- [在线接口文档](https://apifox.com/apidoc/shared-7b67c918-5f72-4f64-b71d-0593d7427b93)

### 联系方式

如有疑问或需要后端配合，请提交 Issue。

---

**最后更新**: 2025-12-11
**PR**: #6 - copilot/finish-incomplete-tasks
