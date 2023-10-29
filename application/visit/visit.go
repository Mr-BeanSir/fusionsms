package visit

import (
	"fusionsms/common"
	"fusionsms/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"reflect"
	"regexp"
	"time"
)

func Reg(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	captcha := c.PostForm("code")
	if username == "" || password == "" || captcha == "" {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "用户名或密码或验证码为空",
		})
		return
	}
	if matchString, _ := regexp.MatchString("^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(.[a-zA-Z0-9_-]+)+$", username); !matchString {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "用户名请使用邮箱注册",
		})
		return
	}
	db := common.GetDbInfo(c)
	defer db.Close()
	code := ""
	err := db.QueryRow("select code from reg_temp where ip = ?", c.ClientIP()).Scan(&code)
	if err != nil {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "获取邮箱验证码失败",
		})
		return
	}
	if captcha != code {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "与邮箱验证码不同",
		})
		return
	}
	stmt, _ := db.Prepare("insert into user(`username`,`password`,`key`) values(?,?,?)")
	exec, err := stmt.Exec(username, common.Md5V(password), common.RandomString(32))
	if err != nil {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "注册失败，请尝试更换用户名",
		})
		return
	}
	if rows, _ := exec.RowsAffected(); rows == 1 {
		c.JSON(200, gin.H{
			"status": 1,
			"msg":    "注册成功",
		})
		return
	}
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "" || password == "" {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "用户名或密码为空",
		})
		return
	}
	db := common.GetDbInfo(c)
	defer db.Close()
	stmt, _ := db.Prepare("select uid,username,password,`group`,quota,jwt,status from user where username = ? and password = ?")
	var user model.User
	err := stmt.QueryRow(username, common.Md5V(password)).Scan(&user.Uid, &user.Username, &user.Password, &user.Group,
		&user.Quota, &user.Jwt, &user.Status)
	if reflect.DeepEqual(user, model.User{}) {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "错误的用户名或密码",
		})
		return
	}
	if user.Status != 0 {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "用户状态不正常，请联系站长",
		})
		return
	}
	jwts := common.CreateJwt(&model.UserClaims{
		UserName: user.Username,
		PassWord: user.Password,
		Uid:      user.Uid,
		Status:   user.Status,
		Group:    user.Group,
		Quota:    user.Quota,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})
	stmt, _ = db.Prepare("update user set jwt = ? where uid = ?")
	exec, err := stmt.Exec(jwts, user.Uid)
	if err != nil {
		c.JSON(200, gin.H{
			"status": -102,
			"msg":    "登录失败，authorization更新失败",
		})
		return
	}
	upNum, err := exec.RowsAffected()
	if err != nil || upNum != 1 {
		c.JSON(200, gin.H{
			"status": -102,
			"msg":    "登录失败，获取author更新数失败",
		})
	}
	c.JSON(200, gin.H{
		"status": 1,
		"msg":    "登录成功",
		"token":  jwts,
	})
	return
}
func Forgot(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	code := c.PostForm("code")
	if username == "" || password == "" || code == "" {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "用户名或密码活或验证码为空",
		})
		return
	}
	if matchString, _ := regexp.MatchString("^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(.[a-zA-Z0-9_-]+)+$", username); !matchString {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "用户名请使用邮箱注册",
		})
		return
	}
	db := common.GetDbInfo(c)
	stmt, _ := db.Prepare("select count(*) from user where username = ? and code = ?")
	var num int
	stmt.QueryRow(username, code).Scan(&num)
	if num != 1 {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "邮箱或验证码错误",
		})
		return
	}
	stmt, _ = db.Prepare("update user set password = ?,code = '' where username = ? and code = ?")
	exec, err := stmt.Exec(common.Md5V(password), username, code)
	if err != nil {
		c.JSON(200, gin.H{
			"status": -102,
			"msg":    "修改密码失败：" + err.Error(),
		})
		return
	}
	upNum, err := exec.RowsAffected()
	if err != nil || upNum != 1 {
		c.JSON(200, gin.H{
			"status": -102,
			"msg":    "修改密码失败，更新密码时失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 1,
		"msg":    "修改密码成功",
	})
	return
}

func EmailCaptcha(c *gin.Context) {
	email := c.PostForm("email")
	if email == "" {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "错误传参",
		})
		return
	}
	if matchString, _ := regexp.MatchString("^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(.[a-zA-Z0-9_-]+)+$", email); !matchString {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "用户名请使用邮箱注册",
		})
		return
	}
	db := common.GetDbInfo(c)
	code := common.RandomString(4)
	err := common.SendWithContext(c, "验证码邮件", "您的验证码为："+code+"，5分钟有效", email)
	if err != nil {
		//common.Log(c, err.Error())
		log.Println(err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "发送验证码失败，已记录错误",
		})
		return
	}
	_, err = db.Exec("replace into reg_temp values (?,?)", c.ClientIP(), code)
	if err != nil {
		log.Println(err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "设置验证码失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 1,
		"msg":    "验证码发送成功",
	})
}

func ForgotCaptcha(c *gin.Context) {
	email := c.PostForm("email")
	if email == "" {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "错误传参",
		})
		return
	}
	if matchString, _ := regexp.MatchString("^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(.[a-zA-Z0-9_-]+)+$", email); !matchString {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "用户名请使用邮箱注册",
		})
		return
	}
	db := common.GetDbInfo(c)
	code := common.RandomString(4)
	err := common.SendWithContext(c, "验证码邮件", "您的验证码为："+code+"，5分钟有效", email)
	if err != nil {
		//common.Log(c, err.Error())
		log.Println(err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "发送验证码失败，已记录错误",
		})
		return
	}
	_, err = db.Exec("update user set code = ? where username = ?", code, email)
	if err != nil {
		log.Println(err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "设置验证码失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 1,
		"msg":    "验证码发送成功",
	})
}
