package loop

import (
	"context"
	"forge/biz/entity"
	"forge/constant"
	"forge/pkg/log/zlog"
	cozeloop "github.com/coze-dev/cozeloop-go"
)

var client cozeloop.Client

func MustInitLoop() {
	// todo需要初始化
	// ref: https://github.com/coze-dev/cozeloop-go
	_client, err := cozeloop.NewClient()
	if err != nil {
		panic(err)
	}
	client = _client

}

func GetNewSpan(ctx context.Context, spanName string, spanType constant.LoopSpanType, opts ...cozeloop.StartSpanOption) (sCtx context.Context, sp cozeloop.Span) {
	// 集成业务信息，可扩展
	defer func() {
		user, ok := entity.GetUser(ctx)
		if ok {
			sp.SetUserID(ctx, user.UserID)
		}
		logid, ok := zlog.GetLogId(ctx)
		if ok {
			sp.SetLogID(ctx, logid)
		}
	}()
	sCtx, sp = client.StartSpan(ctx, spanName, spanType.String(), opts...)
	return sCtx, sp
}

func SetSpanAllInOne(ctx context.Context, sp cozeloop.Span, input, output any, err error) {
	sp.SetInput(ctx, input)
	sp.SetOutput(ctx, output)
	sp.SetError(ctx, err)
	sp.Finish(ctx)
	client.Flush(ctx) // todo这样会有性能问题，但是我们量级太小了无所谓
}
func SetSpanInput(ctx context.Context, sp cozeloop.Span, input any) {
	sp.SetInput(ctx, input)
}
func SetSpanOutput(ctx context.Context, sp cozeloop.Span, output any) {
	sp.SetOutput(ctx, output)
}
func SetSpanFinish(ctx context.Context, sp cozeloop.Span) {
	sp.Finish(ctx)
	client.Flush(ctx) // todo这样会有性能问题，但是我们量级太小了无所谓
}
func SetSpanInputWithTags(ctx context.Context, sp cozeloop.Span, input any, tags map[string]interface{}) {
	sp.SetInput(ctx, input)
	sp.SetTags(ctx, tags)
}
