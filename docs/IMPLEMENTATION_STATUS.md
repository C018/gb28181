# 功能实现状态 (Implementation Status)

本文档跟踪所请求功能的实现状态。

## 请求的功能列表

基于 PR 评论中的需求，以下是各功能的实现状态：

### 1. ✅ 网页设置国标 ID 的地方，显示长度计数 (20 位是标准的)

**状态**: 后端已完成

**实现内容**:
- ✅ 添加 GB28181 设备 ID 20 位长度验证
- ✅ 验证纯数字格式
- ✅ 提供清晰的错误提示信息

**API 行为**:
```json
POST /devices
{
  "type": "GB28181",
  "device_id": "34020000001320000001"  // 必须 20 位数字
}

// 错误示例
{
  "device_id": "3402000000132"  // 只有 13 位
}
Response: "GB28181 设备 ID 必须为 20 位（当前: 13 位），少或多都可能导致信令故障"
```

**前端需要**:
- 添加输入框的字符计数显示
- 实时显示 "当前: X/20 位"
- 输入验证提示

**相关提交**: 4440725

---

### 2. ✅ 完善最简单的账号密码登录，通过配置文件指定账号/密码

**状态**: 已完成

**实现内容**:
- ✅ 支持通过配置文件配置账号密码
- ✅ 提供登录 API `/user/login`
- ✅ 返回 JWT token（有效期 3 天）
- ✅ 提供修改密码 API `/user/user`

**配置示例**:
```toml
[server]
username = "admin"
password = "admin"
```

**API 使用**:
```bash
# 登录
POST /user/login
{
  "username": "admin",
  "password": "admin"
}

# 修改密码（需要 token）
PUT /user/user
Authorization: Bearer <token>
{
  "username": "newadmin",
  "password": "newpassword"
}
```

**相关代码**: `internal/web/api/user.go`

---

### 3. ✅ 所有接口鉴权访问

**状态**: 已完成

**实现内容**:
- ✅ 所有设备管理 API 需要鉴权
- ✅ 所有通道管理 API 需要鉴权
- ✅ PTZ 控制需要鉴权
- ✅ 录像回放需要鉴权
- ✅ 报警订阅需要鉴权
- ✅ 配置管理需要鉴权
- ✅ 流媒体代理需要鉴权
- ✅ 实时通知订阅需要鉴权

**公开端点（无需鉴权）**:
- `/health` - 健康检查
- `/app/metrics/api` - 监控指标
- `/user/login` - 登录接口
- `/webhook/*` - 流媒体服务器回调（内部调用）

**相关提交**: 4440725

---

### 4. ⏸️ 合一镜像优化体积，从 745MB 优化到 137MB

**状态**: 待实现

**说明**: 
- 需要 Dockerfile 优化
- 使用多阶段构建
- 清理不必要的依赖
- 建议单独 issue 跟踪

---

### 5. ⏸️ ONVIF

**状态**: 基础支持已有，待扩展

**当前状态**:
- ✅ 支持 ONVIF 设备添加
- ✅ 基本的 ONVIF 适配器

**待完善**:
- 更多 ONVIF 功能支持
- ONVIF 设备发现
- ONVIF 事件订阅

**说明**: 建议根据具体需求单独 PR 实现

---

### 6. ⏸️ 国际化，支持中文和 English

**状态**: 待实现

**说明**:
- 需要 i18n 框架集成
- 前端和后端都需要支持
- 建议单独 issue 跟踪

---

### 7. ⏸️ 快照

**状态**: 基础功能已有，待完善

**当前状态**:
- ✅ 快照存储功能
- ✅ 快照读取 API
- ✅ 毛玻璃效果配置

**待完善**:
- 自动快照
- 快照定时清理
- 快照管理界面

---

### 8. ✅ 反向代理，见 README 反代使用说明

**状态**: 已完成

**实现内容**:
- ✅ 流媒体服务反向代理 `/proxy/sms/*`
- ✅ 添加鉴权保护
- ✅ 处理重定向
- ✅ 跨域处理

**相关代码**: `internal/web/api/api.go` 中的 `proxySMS` 函数

**相关提交**: 4440725

---

### 9. ✅ 支持 MySQL 数据库

**状态**: 已支持

**实现内容**:
- ✅ 支持 SQLite（默认）
- ✅ 支持 PostgreSQL
- ✅ 支持 MySQL

**配置示例**:
```toml
[data.database]
# SQLite
dsn = "./data/gb28181.db"

# MySQL
dsn = "user:pass@tcp(127.0.0.1:3306)/gb28181?charset=utf8mb4&parseTime=True&loc=Local"

# PostgreSQL
dsn = "host=localhost user=postgres password=postgres dbname=gb28181 port=5432 sslmode=disable"
```

**相关代码**: `internal/conf/config.go`, 使用 GORM 自动支持

---

### 10. ✅ 实时消息通知

**状态**: 已完成并重构

**实现内容**:
- ✅ SSE (Server-Sent Events) 实时推送
- ✅ 设备上线/离线通知
- ✅ 流开始/停止通知
- ✅ 录像开始/停止通知
- ✅ 报警事件通知
- ✅ AI 检测告警通知
- ✅ 依赖注入架构，避免全局变量

**API 端点**:
```
GET /notifications/subscribe
Authorization: Bearer <token>
```

**相关提交**: 178ad02, fd43870

---

### 11. ✅ 播放鉴权，带时效，比如 10 分钟内的链接可以播放

**状态**: 已完成并文档化

**实现内容**:
- ✅ 基于时间戳的播放令牌
- ✅ MD5 签名验证
- ✅ 可配置有效期
- ✅ 在 onPlay webhook 中验证

**令牌格式**: `timestamp.MD5(timestamp + secret)`

**配置**:
```toml
[server]
# 播放链接有效期(分钟)，0 表示不限制
play_expire_minutes = 10
```

**相关代码**: `internal/web/api/ipc.go` 中的 `generatePlayToken` 和 `validatePlayToken`

---

### 12. ✅ 支持在网页上修改账号和密码

**状态**: 后端 API 已完成

**实现内容**:
- ✅ 提供修改密码 API
- ✅ 更新配置文件
- ✅ 需要鉴权访问

**API**:
```
PUT /user/user
Authorization: Bearer <token>
{
  "username": "newadmin",
  "password": "newpassword"
}
```

**前端需要**:
- 密码修改界面
- 表单验证
- 成功/失败提示

**相关代码**: `internal/web/api/user.go`

---

### 13. ✅ 通道封面快照增加毛玻璃效果，可以配置"毛玻璃"开关

**状态**: 已配置化并文档化

**实现内容**:
- ✅ 毛玻璃效果配置选项
- ✅ 文档说明使用场景

**配置**:
```toml
[server]
# 是否启用快照毛玻璃效果
enable_snapshot_blur = false
```

**使用场景**:
- 未授权用户预览
- 免费试看
- 内容保护

**相关文档**: `docs/FEATURES.md`

---

### 14. ✅ 录像

**状态**: API 已组织并完善

**实现内容**:
- ✅ 录像回放 API (`/playback/start`, `/playback/stop`)
- ✅ 录像控制 API (暂停/继续/倍速)
- ✅ 录像查询 API
- ✅ 与直播流结构一致

**API 文件**: `internal/web/api/playback.go`

**相关提交**: 178ad02

---

## 总结

### ✅ 已完成的功能 (10/14)

1. GB28181 设备 ID 验证（后端）
2. 账号密码登录系统
3. 所有接口鉴权
4. 反向代理（含鉴权）
5. MySQL 数据库支持
6. 实时消息通知
7. 播放鉴权带时效
8. 修改账号密码 API
9. 毛玻璃效果配置
10. 录像 API

### 🚧 需要前端配合 (2/14)

1. GB28181 ID 输入长度计数器（后端已完成验证）
2. 登录和密码修改界面（后端 API 已完成）

### ⏸️ 待实现 (2/14)

1. Docker 镜像优化（需要 Dockerfile 优化工作）
2. 国际化支持（需要 i18n 框架集成）

### 📝 可选扩展 (0/14)

1. ONVIF 功能扩展（基础已支持）
2. 快照功能完善（基础已支持）

---

## 下一步建议

1. **前端开发**: 为已完成的后端 API 开发前端界面
   - GB28181 ID 输入框（带长度计数）
   - 登录界面
   - 密码修改界面

2. **Docker 优化**: 单独创建 issue 进行镜像优化

3. **国际化**: 单独创建 issue 规划 i18n 实现

4. **功能扩展**: 根据需要扩展 ONVIF 和快照功能

---

**最后更新**: 2025-12-09
**PR**: copilot/fix-real-time-message-issues
