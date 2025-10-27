package middleware

import (
	"forge/constant"
	"forge/pkg/log/zlog"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// AddTracer
//
//	@Description: add traced in logger
//	@return app.HandlerFunc
func AddTracer() gin.HandlerFunc {
	return func(gCtx *gin.Context) {
		// Trace ID 存在于 HTTP Header "X-Trace-ID" 中 nginx会注入，如果不注入自己注入一个
		logID := gCtx.Request.Header.Get("X-Request-ID")
		if logID == "" {
			logID = uuid.New().String()
			gCtx.Request.Header.Set("X-Request-ID", logID)
		}
		// todo后面集成到coze罗盘平台链路追踪 https://loop.coze.cn/open/docs/cozeloop/sdk
		// 增加Logid
		ctx := gCtx.Request.Context()
		ctx = zlog.WithLogKey(ctx, zap.String(constant.LOGID, logID))
		gCtx.Request.WithContext(ctx)
		gCtx.Next()
		return
	}
}
