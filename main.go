package main

import (
	"flag"
	"fmt"
	"fusionsms/config"
	"fusionsms/middleware"
	"fusionsms/route"
	"github.com/gin-gonic/gin"
)

func main() {

	flag.StringVar(&config.DatabaseIP, "ip", "localhost", "数据库IP地址（默认本地）")
	flag.StringVar(&config.DatabasePort, "port", "3306", "数据库端口（默认3306）")
	flag.StringVar(&config.DatabaseName, "name", "fusionsms", "数据库用户名（必填）")
	flag.StringVar(&config.DatabaseDbName, "dbname", "fusionsms", "数据库名（必填）")
	flag.StringVar(&config.DatabasePassword, "password", "fusionsms", "数据库密码（必填）")
	flag.StringVar(&config.Key, "key", "", "短信豆账户Key（必填）")
	var h = false
	flag.BoolVar(&h, "h", false, "帮助")
	flag.Parse()
	if h {
		flag.Usage()
		return
	}
	//controller.SetBaseValues(config.Key)
	//gin.SetMode(gin.ReleaseMode)
	service := gin.Default() // 创建engine
	middleware.Core(service) // 设置中间件
	route.Core(service)      // 设置路由
	if err := service.Run(":8888"); err != nil {
		fmt.Println("启动失败，原因：", err)
	}
}
