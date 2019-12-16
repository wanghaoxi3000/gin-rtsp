package conf

import (
	"ginrtsp/util"
	"os"

	"github.com/joho/godotenv"
)

// Init 初始化配置项
func Init() {
	// 从本地读取环境变量
	godotenv.Load()

	if os.Getenv("GIN_MODE") == "release" {
		util.BuildLogger("info")
	} else {
		util.BuildLogger("debug")
	}
}
