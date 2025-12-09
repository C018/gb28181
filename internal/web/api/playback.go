package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gowvp/gb28181/internal/core/ipc"
	"github.com/gowvp/gb28181/internal/core/sms"
	"github.com/gowvp/gb28181/pkg/gbs"
	"github.com/ixugo/goddd/pkg/reason"
	"github.com/ixugo/goddd/pkg/web"
)

// registerPlaybackAPI 注册录像回放 API
func registerPlaybackAPI(g gin.IRouter, api IPCAPI, handler ...gin.HandlerFunc) {
	group := g.Group("/playback", handler...)
	group.POST("/start", web.WrapH(api.startPlayback))
	group.POST("/stop", web.WrapH(api.stopPlayback))
	group.POST("/control", web.WrapH(api.playbackControl))
	group.GET("/records", web.WrapH(api.queryRecordInfo))
}

// startPlaybackInput 开始回放输入
type startPlaybackInput struct {
	ChannelID string `json:"channel_id" binding:"required"` // 通道 ID
	StartTime int64  `json:"start_time" binding:"required"` // 开始时间戳(秒)
	EndTime   int64  `json:"end_time" binding:"required"`   // 结束时间戳(秒)
}

// startPlaybackOutput 开始回放输出
type startPlaybackOutput struct {
	StreamID string `json:"stream_id"`
	App      string `json:"app"`
	Stream   string `json:"stream"`
}

// startPlayback 开始录像回放
func (a IPCAPI) startPlayback(c *gin.Context, in *startPlaybackInput) (*startPlaybackOutput, error) {
	ch, err := a.ipc.GetChannel(c.Request.Context(), in.ChannelID)
	if err != nil {
		return nil, err
	}

	svr, err := a.uc.SMSAPI.smsCore.GetMediaServer(c.Request.Context(), sms.DefaultMediaServerID)
	if err != nil {
		return nil, err
	}

	streamID := "playback_" + ch.ID + "_" + time.Now().Format("20060102150405")

	if err := a.uc.SipServer.Playback(&gbs.PlaybackInput{
		Channel: &ipc.Channel{
			ID:        ch.ID,
			DeviceID:  ch.DeviceID,
			ChannelID: ch.ChannelID,
		},
		SMS:        svr,
		StreamMode: 1, // TCP 被动
		StartTime:  in.StartTime,
		EndTime:    in.EndTime,
	}); err != nil {
		return nil, reason.ErrServer.SetMsg(err.Error())
	}

	return &startPlaybackOutput{
		StreamID: streamID,
		App:      "rtp",
		Stream:   streamID,
	}, nil
}

// stopPlaybackInput 停止回放输入
type stopPlaybackInput struct {
	ChannelID string `json:"channel_id" binding:"required"` // 通道 ID
}

// stopPlayback 停止录像回放
func (a IPCAPI) stopPlayback(c *gin.Context, in *stopPlaybackInput) (gin.H, error) {
	ch, err := a.ipc.GetChannel(c.Request.Context(), in.ChannelID)
	if err != nil {
		return nil, err
	}

	if err := a.uc.SipServer.StopPlayback(c.Request.Context(), &gbs.StopPlaybackInput{
		Channel: &ipc.Channel{
			ID:        ch.ID,
			DeviceID:  ch.DeviceID,
			ChannelID: ch.ChannelID,
		},
	}); err != nil {
		return nil, reason.ErrServer.SetMsg(err.Error())
	}

	return gin.H{"msg": "ok"}, nil
}

// playbackControlInput 回放控制输入
type playbackControlInput struct {
	DeviceID  string  `json:"device_id" binding:"required"`  // 设备 ID
	ChannelID string  `json:"channel_id" binding:"required"` // 通道 ID
	Action    string  `json:"action" binding:"required"`     // 动作: play, pause, scale
	Scale     float64 `json:"scale,omitempty"`               // 倍速 (0.5, 1, 2, 4 等)
}

// playbackControl 回放控制 (暂停/继续/倍速)
func (a IPCAPI) playbackControl(c *gin.Context, in *playbackControlInput) (gin.H, error) {
	var control gbs.PlaybackControl
	switch in.Action {
	case "play":
		control = gbs.PlaybackControlPlay
	case "pause":
		control = gbs.PlaybackControlPause
	case "scale":
		control = gbs.PlaybackControlScale
		if in.Scale <= 0 {
			in.Scale = 1.0
		}
	default:
		return nil, reason.ErrBadRequest.SetMsg("不支持的控制操作: " + in.Action)
	}

	if err := a.uc.SipServer.PlaybackControl(in.DeviceID, in.ChannelID, control, in.Scale); err != nil {
		return nil, reason.ErrServer.SetMsg(err.Error())
	}

	return gin.H{"msg": "ok"}, nil
}

// queryRecordInfoInput 查询录像信息输入
type queryRecordInfoInput struct {
	DeviceID  string `form:"device_id" binding:"required"`  // 设备 ID
	ChannelID string `form:"channel_id" binding:"required"` // 通道 ID
	StartTime int64  `form:"start_time" binding:"required"` // 开始时间戳(秒)
	EndTime   int64  `form:"end_time" binding:"required"`   // 结束时间戳(秒)
}

// queryRecordInfo 查询设备端录像信息
func (a IPCAPI) queryRecordInfo(c *gin.Context, in *queryRecordInfoInput) (*gbs.Records, error) {
	records, err := a.uc.SipServer.QueryRecordInfo(in.DeviceID, in.ChannelID, in.StartTime, in.EndTime)
	if err != nil {
		return nil, reason.ErrServer.SetMsg(err.Error())
	}
	return records, nil
}
