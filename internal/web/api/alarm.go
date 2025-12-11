package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ixugo/goddd/pkg/reason"
	"github.com/ixugo/goddd/pkg/web"
)

// registerAlarmAPI 注册报警 API
func registerAlarmAPI(g gin.IRouter, api IPCAPI, handler ...gin.HandlerFunc) {
	group := g.Group("/alarms", handler...)
	group.POST("/subscribe", web.WrapH(api.alarmSubscribe))
	group.POST("/unsubscribe", web.WrapH(api.alarmUnsubscribe))
}

// alarmSubscribeInput 报警订阅输入
type alarmSubscribeInput struct {
	DeviceID      string `json:"device_id" binding:"required"` // 设备 ID
	ExpireSeconds int    `json:"expire_seconds"`               // 订阅有效期(秒), 默认 3600
}

// alarmSubscribe 报警订阅
func (a IPCAPI) alarmSubscribe(c *gin.Context, in *alarmSubscribeInput) (gin.H, error) {
	if in.ExpireSeconds <= 0 {
		in.ExpireSeconds = 3600 // 默认 1 小时
	}

	if err := a.uc.SipServer.AlarmSubscribe(in.DeviceID, in.ExpireSeconds); err != nil {
		return nil, reason.ErrServer.SetMsg(err.Error())
	}

	// 发送订阅成功通知
	a.uc.NotificationAPI.hub.NotifyAlarmSubscriptionChanged(in.DeviceID, true)

	return gin.H{"msg": "ok", "expires": in.ExpireSeconds}, nil
}

// alarmUnsubscribeInput 取消报警订阅输入
type alarmUnsubscribeInput struct {
	DeviceID string `json:"device_id" binding:"required"` // 设备 ID
}

// alarmUnsubscribe 取消报警订阅
func (a IPCAPI) alarmUnsubscribe(c *gin.Context, in *alarmUnsubscribeInput) (gin.H, error) {
	if err := a.uc.SipServer.AlarmUnsubscribe(in.DeviceID); err != nil {
		return nil, reason.ErrServer.SetMsg(err.Error())
	}

	// 发送取消订阅通知
	a.uc.NotificationAPI.hub.NotifyAlarmSubscriptionChanged(in.DeviceID, false)

	return gin.H{"msg": "ok"}, nil
}
