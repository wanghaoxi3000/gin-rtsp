package main

import (
	"ginrtsp/conf"
	"ginrtsp/server"
	"ginrtsp/service"
)

func main() {
	// 从配置文件读取配置, 初始化各个模块
	conf.Init()

	// 装载路由
	r := server.NewRouter()

	go service.WsManager.Start()
	r.Run(":3000")
}
