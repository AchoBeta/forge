package main

import (
	"forge/initalize"
	router "forge/internal/example_hertzx/routerh"
	"forge/pkg/log/zlog"
)

func main() {

	initalize.Init()
	// 释放资源 todo优雅退出
	defer initalize.Eve()
	router.RunServer()
	zlog.Infof("server is done")

}
