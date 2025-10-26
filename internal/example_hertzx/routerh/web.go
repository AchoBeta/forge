package router

import (
	"fmt"
	"forge/global"
	_ "forge/internal/example_hertzx/apih"
	manager "forge/internal/example_hertzx/managerh"
	_ "forge/internal/example_hertzx/middlewareh"
	"forge/pkg/log/zlog"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func RunServer() {
	h, err := listen()
	if err != nil {
		zlog.Errorf("Listen error: %v", err)
		panic(err.Error())
	}
	h.Spin()
}

func listen() (*server.Hertz, error) {

	h := server.Default(server.WithHostPorts(fmt.Sprintf("%s:%d", global.Config.App.Host, global.Config.App.Port)))
	manager.RouteHandler.Register(h)
	return h, nil
}
