package server

import (
	"ginrtsp/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 处理跨域请求,支持options访问
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

// NewRouter Gin 路由配置
func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(Cors())

	// 路由
	r.GET("/ping", api.Ping)
	route := r.Group("/stream")
	{
		route.POST("/play", api.PlayRTSP)
		route.POST("/upload/:channel", api.Mpeg1Video)
		route.GET("/live/:channel", api.Wsplay)
	}

	return r
}
