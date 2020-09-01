package api

import (
	"bufio"
	"ginrtsp/service"

	"github.com/gin-gonic/gin"
)

// PlayRTSP 启动 FFMPEG 播放 RTSP 流
func PlayRTSP(c *gin.Context) {
	srv := &service.RTSPTransSrv{}
	if err := c.ShouldBind(srv); err != nil {
		c.JSON(400, errorRequest(err))
		return
	}

	ret := srv.Service()
	c.JSON(ret.Code, ret)
}

// Mpeg1Video 接收 mpeg1vido 数据流
func Mpeg1Video(c *gin.Context) {
	bodyReader := bufio.NewReader(c.Request.Body)

	for {
		data, err := bodyReader.ReadBytes('\n')
		if err != nil {
			break
		}

		service.WsManager.Groupbroadcast(c.Param("channel"), data)
	}
}

// Wsplay 通过 websocket 播放 mpegts 数据
func Wsplay(c *gin.Context) {
	service.WsManager.RegisterClient(c)
}
