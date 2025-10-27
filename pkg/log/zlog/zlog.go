package zlog

import (
	"context"
	"fmt"
	"forge/constant"
	"github.com/bytedance/gg/gslice"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logKey string

const loggerKey logKey = "logger"
const logDetail logKey = "log_detail"

var logger *zap.Logger

// WithLogKey
//
//	@Description:给指定context添加字段 实现类似traceid作用
//	@param ctx
//	@param fields
//	@return context.Context
func WithLogKey(ctx context.Context, fields ...zapcore.Field) context.Context {
	ctx = context.WithValue(ctx, loggerKey, withContext(ctx).With(fields...))
	// 获取之前的logdetail
	detail := make([]zapcore.Field, 0)
	_detail := ctx.Value(logDetail)
	if _detail != nil {
		detail = _detail.([]zapcore.Field)
	}
	// 深拷贝防止污染
	detail = gslice.Clone(detail)
	detail = append(detail, fields...)
	return context.WithValue(ctx, loggerKey, detail)
}

// 通过ctx获得logid
func GetLogId(ctx context.Context) (string, bool) {
	detail := make([]zapcore.Field, 0)
	_detail := ctx.Value(logDetail)
	if _detail == nil {
		return "", false
	}
	detail = _detail.([]zapcore.Field)
	find := gslice.Find(detail, func(field zapcore.Field) bool {
		return field.String == constant.LOGID
	})
	logid, ok := find.Get()
	if !ok {
		return "", false
	}
	return logid.String, false
}

func InitLogger(zapLogger *zap.Logger) {
	logger = zapLogger
}

// 从指定的context返回一个zap实例
func withContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return logger
	}
	if ctxLogger, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		return ctxLogger
	}
	return logger
}

func Infof(format string, v ...interface{}) {
	logger.Info(fmt.Sprintf(format, v...))
}

func Errorf(format string, v ...interface{}) {
	logger.Error(fmt.Sprintf(format, v...))
}

func Warnf(format string, v ...interface{}) {
	logger.Warn(fmt.Sprintf(format, v...))
}

func Debugf(format string, v ...interface{}) {
	logger.Debug(fmt.Sprintf(format, v...))
}

func Panicf(format string, v ...interface{}) {
	logger.Panic(fmt.Sprintf(format, v...))
}

func Fatalf(format string, v ...interface{}) {
	logger.Fatal(fmt.Sprintf(format, v...))
}

// 下面的logger方法会携带trace id

func CtxInfof(ctx context.Context, format string, v ...interface{}) {
	withContext(ctx).Info(fmt.Sprintf(format, v...))
}

func CtxErrorf(ctx context.Context, format string, v ...interface{}) {
	withContext(ctx).Error(fmt.Sprintf(format, v...))
}

func CtxWarnf(ctx context.Context, format string, v ...interface{}) {
	withContext(ctx).Warn(fmt.Sprintf(format, v...))
}

func CtxDebugf(ctx context.Context, format string, v ...interface{}) {
	withContext(ctx).Debug(fmt.Sprintf(format, v...))
}

func CtxPanicf(ctx context.Context, format string, v ...interface{}) {
	withContext(ctx).Panic(fmt.Sprintf(format, v...))
}

func CtxFatalf(ctx context.Context, format string, v ...interface{}) {
	withContext(ctx).Fatal(fmt.Sprintf(format, v...))
}

func CtxAllInOne(ctx context.Context, action string, input, output any, err error) {
	if err != nil {
		withContext(ctx).Error(action+" failed", zap.Any("input", input), zap.Any("output", output), zap.Error(err))
	} else {
		withContext(ctx).Info(action+" succeed", zap.Any("input", input), zap.Any("output", output))
	}
}
