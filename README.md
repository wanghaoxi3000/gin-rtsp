# gin-rtsp
基于 [**JSMpeg**](https://github.com/phoboslab/jsmpeg/) 的原理，在HTML页面上直接播放RTSP视频流，使用Gin框架开发。


## 主要模块
- API 接口：接收FFMPEG的推流数据和客户端的HTTP请求，将客户端需要播放的RTSP地址转换为一个对应的WebSocket地址，客户端通过这个WebSocket地址便可以直接播放视频，为了及时释放不再观看的视频流，这里设计为客户端播放时需要在每隔60秒的时间里循环请求这个接口，超过指定时间没有收到请求的话后台便会关闭这个视频流。

- FFMPEG 视频转换：收到前端的请求后，启动一个Goroutine调用系统的FFMPEG命令转换指定的RTSP视频流并推送到后台对应的接口，自动结束已超时转换任务。

- WebSocket Manager：管理WebSocket客户端，将请求同一WebSocket地址的客户端添加到一个Group中，向各个Group广播对应的RTSP视频流，删除Group中已断开连接的客户端，释放空闲的Group。


## 注意
- 需要摄像头的码流为H264码流


## 编译
**项目需要运行在安装有FFMPEG程序的环境中。**

### Docker
Dockerfile 中已经封装好了需要的环境，可以使用Docker build后，以Docker的方式运行。
```
$ docker build -t ginrtsp .
$ docker run -td -p 3000:3000 ginrtsp
```

### 本地编译

**linux**

```shell
go build -o ./bin/linux/rtsp-relay
```

**windows**

```shell
go build -o ./bin/windows/rtsp-relay.exe
```

### 环境变量
- RTSP_PORT 默认为3000
- RTSP_CORS 默认为false, 设置为true时跨域，其他不跨域


## 测试
### 使用内置的FFMPEG转换
将需要播放的RTSP流地址提交到 /stream/play 接口，例如：
```
POST /stream/play
{
   "url": "rtsp://admin:password@192.168.3.10:554/cam/realmonitor?channel=1&subtype=0"
}
```

后台可以正常转换此RTSP地址时便会返回一个对应的地址，例如：
```
{
    "code": 0,
    "data": {
        "path": "/stream/live/5b96bff4-bdb2-3edb-9d6e-f96eda03da56"
    },
    "msg": "success"
}
```

编辑`html`文件夹下view-stream.html文件，将script部分的url修改为此地址，在浏览器中打开，便可以看到视频了。

### 手动运行FFMPEG
由于后台转换RTSP的进程在超过60秒没有请求后便会停止，也可以通过手动运行ffmpeg命令，来更方便地在测试状态下查看视频。
```
ffmpeg -rtsp_transport tcp -re -i 'rtsp://admin:password@192.168.3.10:554/cam/realmonitor?channel=1&subtype=0' -q 0 -f mpegts -c:v mpeg1video -an -s 960x540 http://127.0.0.1:3000/stream/upload/test
```

通过如上命令，运行之后在view-stream.html文件的url中填入对应的地址为`/stream/live/test`，在浏览器中打开查看视频。


### 显示效果

![](./video-example.png)


## 参考
[JSMpeg – MPEG1 Video & MP2 Audio Decoder in JavaScript](https://github.com/phoboslab/jsmpeg/)
