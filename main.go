package main

import (
	"fmt"
	"ginrtsp/conf"
	"ginrtsp/server"
	"ginrtsp/service"
	"os"
)

func main() {
	// 从配置文件读取配置, 初始化各个模块
	conf.Init()
	port := "3000"
	if len(os.Getenv("RTSP_PORT")) > 0 {
		port = os.Getenv("RTSP_PORT")
	}
	port = ":" + port

	fmt.Println("RTSP_PORT", port)

	// 装载路由
	r := server.NewRouter()

	go service.WsManager.Start()
	r.Run(port)
}
