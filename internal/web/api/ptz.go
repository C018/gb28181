package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gowvp/gb28181/pkg/gbs"
	"github.com/ixugo/goddd/pkg/reason"
	"github.com/ixugo/goddd/pkg/web"
)

// registerPTZAPI 注册云台控制 API
func registerPTZAPI(g gin.IRouter, api IPCAPI, handler ...gin.HandlerFunc) {
	group := g.Group("/ptz", handler...)
	group.POST("/control", web.WrapH(api.ptzControl))
	group.POST("/preset", web.WrapH(api.ptzPreset))
	group.GET("/presets", web.WrapH(api.queryPresets))
}

// ptzControlInput 云台控制输入
type ptzControlInput struct {
	DeviceID  string `json:"device_id" binding:"required"`  // 设备 ID
	ChannelID string `json:"channel_id" binding:"required"` // 通道 ID
	Command   string `json:"command" binding:"required"`    // 控制命令: stop, left, right, up, down, zoom_in, zoom_out, left_up, left_down, right_up, right_down, iris_in, iris_out, focus_in, focus_out
	Speed     int    `json:"speed"`                         // 速度 (0-255), 默认 50
}

// ptzControl 云台方向控制
func (a IPCAPI) ptzControl(c *gin.Context, in *ptzControlInput) (gin.H, error) {
	if in.Speed <= 0 {
		in.Speed = 50
	}
	if in.Speed > 255 {
		in.Speed = 255
	}

	var cmd byte
	switch in.Command {
	case "stop":
		cmd = gbs.PTZCmdStop
	case "left":
		cmd = gbs.PTZCmdLeft
	case "right":
		cmd = gbs.PTZCmdRight
	case "up":
		cmd = gbs.PTZCmdUp
	case "down":
		cmd = gbs.PTZCmdDown
	case "zoom_in":
		cmd = gbs.PTZCmdZoomIn
	case "zoom_out":
		cmd = gbs.PTZCmdZoomOut
	case "left_up":
		cmd = gbs.PTZCmdLeftUp
	case "left_down":
		cmd = gbs.PTZCmdLeftDown
	case "right_up":
		cmd = gbs.PTZCmdRightUp
	case "right_down":
		cmd = gbs.PTZCmdRightDown
	case "iris_in":
		cmd = gbs.PTZCmdIrisIn
	case "iris_out":
		cmd = gbs.PTZCmdIrisOut
	case "focus_in":
		cmd = gbs.PTZCmdFocusIn
	case "focus_out":
		cmd = gbs.PTZCmdFocusOut
	default:
		return nil, reason.ErrBadRequest.SetMsg("不支持的控制命令: " + in.Command)
	}

	speed := byte(in.Speed)
	ptzCmd := gbs.BuildPTZCmd(cmd, speed, speed, 0)

	if err := a.uc.SipServer.PTZControl(in.DeviceID, in.ChannelID, ptzCmd); err != nil {
		return nil, reason.ErrServer.SetMsg(err.Error())
	}

	return gin.H{"msg": "ok"}, nil
}

// ptzPresetInput 预置位控制输入
type ptzPresetInput struct {
	DeviceID    string `json:"device_id" binding:"required"`    // 设备 ID
	ChannelID   string `json:"channel_id" binding:"required"`   // 通道 ID
	Action      string `json:"action" binding:"required"`       // 动作: set, call, delete
	PresetIndex int    `json:"preset_index" binding:"required"` // 预置位编号 (1-255)
}

// ptzPreset 预置位控制
func (a IPCAPI) ptzPreset(c *gin.Context, in *ptzPresetInput) (gin.H, error) {
	if in.PresetIndex < 1 || in.PresetIndex > 255 {
		return nil, reason.ErrBadRequest.SetMsg("预置位编号必须在 1-255 之间")
	}

	var cmd byte
	switch in.Action {
	case "set":
		cmd = gbs.PTZCmdPresetSet
	case "call":
		cmd = gbs.PTZCmdPresetCall
	case "delete":
		cmd = gbs.PTZCmdPresetDelete
	default:
		return nil, reason.ErrBadRequest.SetMsg("不支持的预置位操作: " + in.Action)
	}

	ptzCmd := gbs.BuildPresetCmd(cmd, byte(in.PresetIndex))

	if err := a.uc.SipServer.PTZControl(in.DeviceID, in.ChannelID, ptzCmd); err != nil {
		return nil, reason.ErrServer.SetMsg(err.Error())
	}

	return gin.H{"msg": "ok"}, nil
}

// queryPresetsInput 查询预置位输入
type queryPresetsInput struct {
	DeviceID  string `form:"device_id" binding:"required"`  // 设备 ID
	ChannelID string `form:"channel_id" binding:"required"` // 通道 ID
}

// queryPresets 查询预置位列表
func (a IPCAPI) queryPresets(c *gin.Context, in *queryPresetsInput) (gin.H, error) {
	// TODO: 实现预置位查询功能
	// 需要在 gbs.Server 中添加 QueryPresets 方法
	// 目前返回空列表表示功能待实现
	return gin.H{
		"msg":     "预置位查询功能待实现",
		"presets": []any{},
	}, nil
}
