# GB28181 功能说明文档

本文档介绍 GB28181 项目中的各项功能实现和配置说明。

## 实时通知系统 (Real-time Notifications)

### 概述
实时通知系统基于 Server-Sent Events (SSE) 实现，用于向前端客户端推送实时消息，包括设备状态、流状态、录像状态、报警事件和 AI 检测告警等。

### 架构设计
- **NotificationHub**: 通知中心，管理所有 SSE 连接
- **NotificationAPI**: 通知 API，提供订阅接口
- **依赖注入**: 通过 Wire 依赖注入框架创建和管理通知中心实例，避免全局变量污染

### 支持的通知类型
- `device_online/device_offline`: 设备上线/离线
- `stream_start/stream_stop`: 流开始/停止
- `record_start/record_stop`: 录像开始/停止
- `error`: 错误通知
- `alarm`: 报警事件
- `alarm_subscribed`: 报警订阅状态变更
- `ai_alert`: AI 检测告警

### 使用方式
客户端通过 `/notifications/subscribe` 端点订阅实时通知：

```
GET /notifications/subscribe
Authorization: Bearer <token>
```

## 播放令牌系统 (Play Token)

### 概述
播放令牌用于控制播放链接的有效期，防止未授权访问和链接滥用。

### 设计原理
- **令牌格式**: `timestamp.sign`
  - `timestamp`: 过期时间戳（秒）
  - `sign`: `MD5(timestamp + secret)`
- **验证流程**: 
  1. 解析令牌获取时间戳和签名
  2. 验证签名是否正确
  3. 检查时间戳是否过期

### 配置
在 `config.toml` 中配置：

```toml
[server]
# 播放链接有效期(分钟)，0 表示不限制
play_expire_minutes = 60
```

### 鉴权集成
- **onPlay 事件**: 在 ZLMediaKit webhook 的 `on_play` 事件中验证令牌
- **onStreamNotFound 事件**: 流不存在时也会进行令牌验证（如果配置了有效期）
- **RTMP 鉴权**: 与 RTMP 推流鉴权独立，使用不同的密钥和验证机制

### 注意事项
- 令牌验证在 `on_play` webhook 中执行
- RTMP 推流使用独立的 `rtmp_secret` 进行鉴权
- 建议为不同的服务使用不同的密钥

## 快照毛玻璃效果 (Snapshot Blur)

### 概述
快照毛玻璃效果用于在未授权或预览模式下显示模糊的视频快照，保护视频内容隐私。

### 配置
在 `config.toml` 中配置：

```toml
[server]
# 是否启用快照毛玻璃效果
enable_snapshot_blur = false
```

### 实现方式
- 快照图片在返回前端之前进行模糊处理
- 仅在配置启用且满足特定条件时应用
- 不影响正常授权播放的画面质量

### 使用场景
- 未授权用户预览
- 免费试看功能
- 内容保护

## 云台控制 (PTZ Control)

### 概述
云台控制功能支持 GB28181 设备的云台方向控制和预置位管理。

### API 端点
```
POST /ptz/control       # 云台方向控制
POST /ptz/preset        # 预置位设置/调用/删除
GET  /ptz/presets       # 查询预置位列表 (待实现)
```

### 方向控制
支持的控制命令：
- `stop`: 停止
- `left/right/up/down`: 左/右/上/下
- `zoom_in/zoom_out`: 放大/缩小
- `left_up/left_down/right_up/right_down`: 组合方向
- `iris_in/iris_out`: 光圈控制
- `focus_in/focus_out`: 焦距控制

### 预置位管理
- `set`: 设置预置位
- `call`: 调用预置位
- `delete`: 删除预置位
- 预置位编号范围: 1-255

### 待实现功能
- 预置位查询: 查询设备当前配置的预置位列表

## 录像回放 (Playback)

### 概述
录像回放功能支持查询和播放 GB28181 设备端存储的历史录像。

### API 端点
```
POST /playback/start     # 开始回放
POST /playback/stop      # 停止回放
POST /playback/control   # 回放控制（暂停/继续/倍速）
GET  /playback/records   # 查询录像信息
```

### 回放控制
- `play`: 继续播放
- `pause`: 暂停播放
- `scale`: 倍速播放（支持 0.5x, 1x, 2x, 4x 等）

### 架构说明
录像回放的结构设计与直播流结构保持一致，确保：
- 统一的流管理接口
- 一致的启动/停止机制
- 相同的错误处理流程

## 报警订阅 (Alarm Subscription)

### 概述
报警订阅功能允许系统订阅 GB28181 设备的报警事件。

### API 端点
```
POST /alarms/subscribe      # 订阅报警
POST /alarms/unsubscribe    # 取消订阅
```

### 订阅参数
- `device_id`: 设备 ID
- `expire_seconds`: 订阅有效期（秒），默认 3600 秒（1 小时）

### 通知集成
订阅/取消订阅操作会触发实时通知，通过 SSE 推送到前端：
- `alarm_subscribed`: 报警订阅状态变更通知
- `alarm`: 报警事件通知

## AI 检测服务 (AI Detection)

### 概述
AI 检测服务提供视频智能分析能力，支持行人检测、车辆检测、人脸识别等功能。

### 推理模式
支持两种推理模式：
1. **本地推理 (local)**: 在服务器本地运行 AI 模型
2. **远程 API (remote)**: 调用远程 AI 服务接口

### 配置
在 `config.toml` 中配置：

```toml
[ai]
# 是否启用 AI 检测服务
enabled = false

# 推理模式: local(本地推理) / remote(远程API)
inference_mode = "remote"

# 远程 AI 服务地址
endpoint = "http://localhost:8080"

# API 密钥
api_key = ""

# 请求超时(秒)
timeout = 30

# 模型类型: yolov5/yolov8/custom
model_type = "yolov8"

# 本地模型文件路径
model_path = "./models/yolov8n.onnx"

# 推理设备类型: cpu/cuda/mps
device_type = "cpu"
```

### 检测类型
- `pedestrian`: 行人检测
- `vehicle`: 车辆检测
- `face`: 人脸识别
- `object`: 通用物体检测

### 告警规则
支持配置检测告警规则：
- 检测类型过滤
- 置信度阈值
- 检测区域限制
- 告警冷却时间

### API 端点
```
POST /ai/detect           # 执行检测
GET  /ai/rules            # 列出告警规则
POST /ai/rules            # 创建告警规则
GET  /ai/rules/:id        # 获取告警规则
PUT  /ai/rules/:id        # 更新告警规则
DELETE /ai/rules/:id      # 删除告警规则
GET  /ai/status           # 获取服务状态
```

### 实现说明
- AI 检测功能为可选功能，默认禁用
- 需要根据实际需求选择合适的推理引擎
- 本地推理需要额外集成 ONNX Runtime 或其他推理框架
- 远程 API 模式需要部署独立的 AI 推理服务

## Go 流媒体服务器 (GoLive)

### 概述
GoLive 是一个纯 Go 语言实现的流媒体服务器，作为 ZLMediaKit 的可选替代方案。

### 配置
在 `config.toml` 中配置：

```toml
[golive]
# 是否启用 Go 流媒体服务
enabled = false

# RTMP 端口
rtmp_port = 1936

# RTSP 端口
rtsp_port = 8555

# HTTP-FLV 端口
http_flv_port = 8088

# HLS 端口
hls_port = 8088

# 公网 IP
public_ip = ""

# 启用推流鉴权
enable_auth = false

# 推流鉴权密钥
auth_secret = ""

# HLS 分片时长(秒)
hls_fragment = 2

# HLS 窗口大小
hls_window = 6

# 录像存储路径
record_path = "./records"

# 启用录像
enable_record = false
```

### 支持的协议
- RTMP: 推拉流
- RTSP: 推拉流
- HTTP-FLV: 拉流
- HLS: 拉流

### 实现状态
⚠️ **注意**: GoLive 当前处于早期开发阶段，功能不完整：
- 基础框架已实现
- HTTP 服务器已启动
- **实际的流处理逻辑待实现**
- 需要集成 LAL 或其他 Go 流媒体库

### PS 流支持
- GoLive 尚未实现完整的 PS 流解析和处理
- GB28181 的 PS 流目前仍需使用 ZLMediaKit 处理

### 使用建议
- **生产环境**: 建议继续使用成熟稳定的 ZLMediaKit
- **开发测试**: 可以启用 GoLive 进行功能测试
- **未来计划**: 逐步完善 GoLive 功能，实现纯 Go 技术栈

### 负载均衡集成
GoLive 当前未在 SMS (流媒体服务器管理) 的负载均衡中初始化，需要在以下位置集成：
- 流媒体服务器注册
- 负载均衡策略
- 健康检查机制

## 技术栈说明

### 流媒体服务器
- **ZLMediaKit** (推荐): C++ 实现，功能完整，性能优秀
  - 完整支持 GB28181 PS 流
  - 支持 RTMP、RTSP、HLS、HTTP-FLV、WebRTC 等多种协议
  - 提供丰富的 Webhook 回调接口
  
- **GoLive** (实验性): Go 语言实现，功能待完善
  - 纯 Go 技术栈，便于集成和扩展
  - 当前处于开发阶段，功能不完整
  - 待集成 LAL 等 Go 流媒体库

### AI 推理
- **本地推理**: 需要集成推理引擎
  - ONNX Runtime: 跨平台支持
  - TensorRT: NVIDIA GPU 加速
  - OpenVINO: Intel 硬件优化
  
- **远程 API**: 调用独立的 AI 服务
  - 支持自定义 AI 服务接口
  - 便于水平扩展和资源管理
  - 需要确保网络延迟可接受

## 架构原则

### 分层架构
- **API 层**: 提供 HTTP/WebSocket 接口
- **业务层**: 实现业务逻辑
- **数据层**: 数据持久化和缓存
- **适配器层**: 协议适配和设备通信

### 依赖原则
- **高层不依赖低层**: API 层不依赖业务层实现
- **面向接口编程**: 通过接口定义契约
- **依赖注入**: 使用 Wire 进行依赖管理
- **避免循环依赖**: 保持清晰的模块关系

### 代码规范
- **职责单一**: 每个文件/模块负责一个功能
- **接口清晰**: 明确的输入输出定义
- **错误处理**: 统一的错误处理机制
- **文档完善**: 关键功能提供文档说明

## 未来规划

### 短期目标
1. 完善预置位查询功能
2. 优化录像回放稳定性
3. 完善 AI 检测功能文档
4. 改进实时通知的消息顺序保证

### 中期目标
1. 完善 GoLive 流媒体服务器实现
2. 支持更多 AI 检测模型
3. 优化播放令牌验证机制
4. 实现更灵活的鉴权策略

### 长期目标
1. 实现纯 Go 语言的完整流媒体处理
2. 支持边缘计算和分布式部署
3. 提供更丰富的视频分析能力
4. 完善监控和运维工具

## 相关文档
- [API 接口文档](./API接口文档.md)
- [README](../README.md)
- [配置文件示例](../configs/)
