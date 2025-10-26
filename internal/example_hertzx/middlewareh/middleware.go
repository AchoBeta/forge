package middleware

import (
	"context"
	"forge/internal/example_hertzx/managerh"
	"forge/pkg/log/zlog"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func init() {
	manager.RouteHandler.RegisterMiddleware(manager.LEVEL_GLOBAL, AddTracer, false)
}
