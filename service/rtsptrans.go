package service

import (
	"fmt"
	"ginrtsp/serializer"

	"ginrtsp/util"
	"os/exec"
	"strings"
	"sync"

	"time"

	uuid "github.com/satori/go.uuid"
)

// RTSPTransSrv RTSP 转换服务 struct
type RTSPTransSrv struct {
	URL string `form:"url" json:"url" binding:"required,min=1"`
}

// processMap FFMPEG 进程刷新通道，未在指定时间刷新的流将会被关闭
var processMap sync.Map

// Service RTSP 转换服务
func (service *RTSPTransSrv) Service() *serializer.Response {
	simpleString := strings.Replace(service.URL, "//", "/", 1)
	splitList := strings.Split(simpleString, "/")

	if splitList[0] != "rtsp:" && len(splitList) < 2 {
		return &serializer.Response{
			Code: 400,
			Msg:  "不是有效的 RTSP 地址",
		}
	}

	// 多个客户端需要播放相同的RTSP流地址时，保证返回WebSocket地址相同
	processCh := uuid.NewV3(uuid.NamespaceURL, splitList[1]).String()
	if ch, ok := processMap.Load(processCh); ok {
		*ch.(*chan int) <- 1
	} else {
		reflush := make(chan int)
		go runFFMPEG(service.URL, processCh, &reflush)
	}

	playURL := fmt.Sprintf("/stream/live/%s", processCh)
	return serializer.BuildRTSPPlayPathResponse(playURL)
}

func runFFMPEG(rtsp string, playCh string, ch *chan int) {
	processMap.Store(playCh, ch)
	defer func() {
		processMap.Delete(playCh)
		util.Log().Info("Stop translate rtsp %v", rtsp)
	}()

	params := []string{
		"-rtsp_transport",
		"tcp",
		"-re",
		"-i",
		rtsp,
		"-q",
		"5",
		"-f",
		"mpegts",
		"-fflags",
		"nobuffer",
		"-c:v",
		"mpeg1video",
		"-an",
		"-s",
		"960x540",
		fmt.Sprintf("http://127.0.0.1:3000/stream/upload/%s", playCh),
	}

	cmd := exec.Command("ffmpeg", params...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	stdin, err := cmd.StdinPipe()
	if err != nil {
		util.Log().Error("Get ffmpeg stdin err:%v", err)
		return
	}
	defer stdin.Close()

	err = cmd.Start()
	if err != nil {
		util.Log().Error("Start ffmpeg err:%v", err.Error)
		return
	}
	util.Log().Info("Translate rtsp %v to %v", rtsp, playCh)

	for {
		select {
		case <-*ch:
			util.Log().Info("reflush channel %s rtsp %v", playCh, rtsp)

		case <-time.After(60 * time.Second):
			stdin.Write([]byte("q"))
			err = cmd.Wait()
			if err != nil {
				util.Log().Error("Run ffmpeg err:%v", err.Error)
			}
			return
		}
	}
}
