package admin

import (
	"fusionsms/common"
	"fusionsms/model"
	"github.com/gin-gonic/gin"
)

func GetSent(c *gin.Context) {
	db := common.GetDbInfo(c)
	var system model.System
	err := db.QueryRow("select smtp_server,smtp_username,smtp_password,smtp_nickname,smtp_ssl from `system`").
		Scan(&system.SmtpServer, &system.SmtpUsername, &system.SmtpPassword, &system.SmtpNickname, &system.SmtpSSL)
	if err != nil {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "获取系统发信配置失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 1,
		"msg":    "获取系统发信配置成功",
		"data":   system,
	})
}

func SetSent(c *gin.Context) {
	server := c.PostForm("server")
	username := c.PostForm("username")
	password := c.PostForm("password")
	nickname := c.PostForm("nickname")
	ssl := c.PostForm("ssl")
	if server == "" || username == "" || password == "" || nickname == "" || ssl == "" {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "空传参",
		})
		return
	}
	db := common.GetDbInfo(c)
	_, err := db.Exec("update `system` set smtp_server=?,smtp_username=?,smtp_password=?,smtp_ssl=?,smtp_nickname=?", server, username, password, ssl, nickname)
	if err != nil {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 1,
		"msg":    "修改发信配置成功",
	})
}
