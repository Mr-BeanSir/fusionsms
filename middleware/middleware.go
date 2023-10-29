package middleware

import (
	"database/sql"
	"fusionsms/common"
	"fusionsms/config"
	"fusionsms/model"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

func Core(service *gin.Engine) {
	service.Use(SetHeader)
	service.Use(Auth)
}

func SetHeader(c *gin.Context) {
	method := c.Request.Method
	origin := c.Request.Header.Get("Origin")
	if origin != "" {
		c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
	}
	if method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
	}
}

func Auth(c *gin.Context) {
	// 添加屏蔽的uri
	arrUri := [...]string{
		"/user/reg",
		"/user/login",
		"/user/forgot",
		"/user/emailCaptcha",
		"/user/forgotCaptcha",
		"/send",
		"/controller/Api/ReceiveTemplateStatus",
		"/controller/Api/ReceiveSignStatus",
		"/controller/Api/ReceiveSentStatus",
	}
	// 访问路径为屏蔽uri直接通过
	for _, value := range arrUri {
		if c.Request.RequestURI == value {
			db, err := config.GetDb()
			if err != nil {
				c.JSON(http.StatusForbidden, gin.H{
					"status": "-1",
					"msg":    "数据库连接错误",
				})
				c.Abort()
				return
			}
			c.Set("createJwt", common.CreateJwt)
			c.Set("db", db)
			c.Next()
			return
		}
	}
	// 获取jwt
	token := c.Request.Header.Get("Authorization")
	if token == "nil" {
		c.JSON(http.StatusForbidden, gin.H{
			"status": "-101",
			"msg":    "请携带authorization a",
		})
		c.Abort()
		return
	}
	// 解析jwt并验证
	jwt, err := common.ParseJwt(token)
	if err != nil || !jwt.Valid {
		log.Println(jwt.Valid)
		c.JSON(http.StatusForbidden, gin.H{
			"status": "-101",
			"msg":    "错误的authorization b",
		})
		c.Abort()
		return
	}
	// 获取数据库 并设置到context
	db, err := config.GetDb()
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"status": "-1",
			"msg":    "数据库连接错误",
		})
		c.Abort()
		return
	}
	// 查找authorization是否存在于数据库，方便管理
	stmt, _ := db.Prepare("select count(0) from user where jwt = ?")
	query := stmt.QueryRow(token)
	var temp struct {
		num int
	}
	_ = query.Scan(&temp.num)
	if temp.num == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"status": "-101",
			"msg":    "错误的authorization c",
		})
		c.Abort()
		return
	}
	if c.Request.RequestURI[0:6] == "/admin" {
		jwt1, _ := jwt.Claims.(*model.UserClaims)
		if !common.IsAdmin(db, jwt1.Uid) {
			c.JSON(200, gin.H{
				"status": -1,
				"msg":    "您似乎并没有管理员权限，或数据库读取错误",
			})
			c.Abort()
			return
		}
	}
	c.Set("db", db)
	c.Set("claims", jwt)
	c.Next()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
}
