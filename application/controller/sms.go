package controller

import (
	"database/sql"
	"fmt"
	"fusionsms/common"
	"fusionsms/model"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"io"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

func LocalTest(c *gin.Context) {
	sid := c.PostForm("sid")
	to := c.PostForm("to")
	content := c.PostForm("content")
	if sid == "" || to == "" || content == "" {
		c.JSON(200, gin.H{
			"status": 301,
			"msg":    "错误传参",
			"bool":   false,
		})
		return
	}
	db := common.GetDbInfo(c)
	jwt := common.GetJwtInfo(c)
	if jwt.Status != 0 {
		c.JSON(200, gin.H{
			"status": 123,
			"msg":    "账号状态不正常，请联系站长",
			"bool":   false,
		})
		return
	}
	var (
		userinfo model.User
		signinfo model.Sign
		limit    model.Limit
	)

	stmt, _ := db.Prepare("select `group`,quota,status,phone,limit_quota_phone,black_phone from user where uid = ?")
	stmt1, _ := db.Prepare("select content,status,`key` from sign where sid = ?")
	stmt2, _ := db.Prepare("select num,time from `limit` where sid = ?")
	err := stmt.QueryRow(jwt.Uid).Scan(&userinfo.Group, &userinfo.Quota, &userinfo.Status, &userinfo.Phone, &userinfo.LimitQuotaPhone, &userinfo.BlackPhone)
	err1 := stmt1.QueryRow(sid).Scan(&signinfo.Content, &signinfo.Status, &signinfo.Key)
	rows, err2 := stmt2.Query(sid)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -500,
			"msg":    "获取用户信息错误，错误已记录",
			"bool":   false,
		})
		return
	}
	if err1 != nil {
		common.Log(c, err1.Error())
		c.JSON(200, gin.H{
			"status": -500,
			"msg":    "获取用户信息错误1，错误已记录",
			"bool":   false,
		})
		return
	}
	if err1 != nil {
		common.Log(c, err2.Error())
		c.JSON(200, gin.H{
			"status": -500,
			"msg":    "获取用户信息错误2，错误已记录",
			"bool":   false,
		})
		return
	}
	if userinfo.Status != 0 || signinfo.Status != 1 {
		c.JSON(200, gin.H{
			"status": -401,
			"msg":    "用户被封禁或签名头未审核",
			"bool":   false,
		})
		return
	}
	if userinfo.Quota < 1 {
		c.JSON(200, gin.H{
			"status": 305,
			"msg":    "用户账户额度不足",
			"bool":   false,
		})
		return
	}
	if strings.Contains(userinfo.BlackPhone, to) {
		c.JSON(200, gin.H{
			"status": 307,
			"msg":    "黑名单手机号，不允许发送",
			"bool":   false,
		})
		return
	}
	for rows.Next() {
		err := rows.Scan(&limit.Num, &limit.Time)
		if err != nil {
			c.JSON(200, gin.H{
				"status": 508,
				"msg":    "处理限流操作时错误，请联系站长",
				"bool":   false,
			})
			return
		}
		var (
			formerlyStmt *sql.Stmt
			num          int
		)

		switch limit.Time {
		case "s":
			formerlyStmt, _ = db.Prepare("select count(*) from sent_log where sid = ? and time > date_sub(now(),interval 1 second )")
			break
		case "m":
			formerlyStmt, _ = db.Prepare("select count(*) from sent_log where sid = ? and time > date_sub(now(),interval 1 minute )")
			break
		case "h":
			formerlyStmt, _ = db.Prepare("select count(*) from sent_log where sid = ? and time > date_sub(now(),interval 1 hour )")
			break
		case "d":
			formerlyStmt, _ = db.Prepare("select count(*) from sent_log where sid = ? and time > date_sub(now(),interval 1 day )")
			break
		default:
			c.JSON(200, gin.H{
				"status": 320,
				"msg":    "错误的类型，请检查操作",
				"bool":   false,
			})
			return
		}
		err = formerlyStmt.QueryRow(sid).Scan(&num)
		if err != nil {
			c.JSON(200, gin.H{
				"status": 321,
				"msg":    "错误查询，数据库表不存在？",
				"bool":   false,
			})
			return
		}
		if num > limit.Num {
			c.JSON(200, gin.H{
				"status": 322,
				"msg":    "超过最大发信限制，请等等再发",
				"bool":   false,
			})
			return
		}
	}
	sign := content[strings.Index(content, "【") : strings.Index(content, "】")+3]
	if sign != signinfo.Content {
		c.JSON(200, gin.H{
			"status": 308,
			"msg":    "用户错误，传递签名错误",
			"bool":   false,
		})
		return
	}
	TemplatesStmt, _ := db.Prepare("select content,status from template where sid = ?")
	TemplatesRows, err := TemplatesStmt.Query(sid)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "获取报备模板列表失败，错误已记录",
			"bool":   false,
		})
		return
	}
	var template model.Template
	for TemplatesRows.Next() {
		err := TemplatesRows.Scan(&template.Content, &template.Status)
		if err != nil {
			common.Log(c, err.Error())
			c.JSON(200, gin.H{
				"status": -1,
				"msg":    "扫描字段出错，错误已记录",
				"bool":   false,
			})
			return
		}
		template.Content = strings.ReplaceAll(template.Content, "@", "(.+?)")
		match, err := regexp.Match(template.Content, []byte(content))
		if err != nil {
			c.JSON(200, gin.H{
				"status": -1,
				"msg":    "验证模板错误，请检查传递内容",
				"bool":   false,
			})
			return
		}
		if !match || template.Status != 1 {
			c.JSON(200, gin.H{
				"status": -1,
				"msg":    "模板不存在或模板未审核",
				"bool":   false,
			})
			return
		}
	}
	DescStmt, err := db.Prepare("insert into sent_log(uid, sid, content, phone, time, decrease_num) VALUE (?,?,?,?,now(),?)")
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "添加日志失败",
			"bool":   false,
		})
		return
	}
	DescNum := math.Ceil(float64(utf8.RuneCountInString(content)) / 64)
	result, err := DescStmt.Exec(jwt.Uid, sid, content, to, DescNum)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "添加日志失败2",
			"bool":   false,
		})
		return
	}
	index, err := result.LastInsertId()
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "添加日志失败，获取下标错误",
			"bool":   false,
		})
		return
	}
	BaseValues.Add("sign", signinfo.Key)
	BaseValues.Add("to", to)
	BaseValues.Add("content", content)
	form, err := http.PostForm(BaseUrl+"/Api/Sent", BaseValues)
	if err != nil {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "构造socket错误",
			"bool":   false,
		})
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(form.Body)
	all, err := io.ReadAll(form.Body)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "读取返回错误，错误已记录",
			"bool":   false,
		})
		return
	}
	parse := gjson.Parse(string(all))
	if !parse.Get("bool").Bool() {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    parse.Get("msg").String(),
			"bool":   false,
		})
		return
	}
	stmt, _ = db.Prepare("update sent_log set task_id = ?,status = 3 where id = ?")
	_, err = stmt.Exec(parse.Get("task_id").Int(), index)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "更新发信日志失败，已记录错误",
			"bool":   false,
			"data":   string(all),
		})
		return
	}
	stmt, _ = db.Prepare("update user set quota = quota - ? where uid = ?")
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    fmt.Sprintf("更新用户余额失败，%d减少额度%f", userinfo.Uid, DescNum),
			"bool":   false,
		})
		return
	}
	_, err = stmt.Exec(DescNum, jwt.Uid)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    fmt.Sprintf("更新用户余额失败2，%d减少额度%f", userinfo.Uid, DescNum),
			"bool":   false,
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 1,
		"msg":    parse.Get("msg").String(),
		"bool":   true,
	})
}

func Send(c *gin.Context) {
	signKey := c.PostForm("sign")
	key := c.PostForm("key")
	to := c.PostForm("to")
	content := c.PostForm("content")
	if key == "" || to == "" || content == "" || signKey == "" {
		c.JSON(200, gin.H{
			"status": 301,
			"msg":    "错误传参",
			"bool":   false,
		})
		return
	}
	db := common.GetDbInfo(c)
	var (
		userinfo model.User
		signinfo model.Sign
		limit    model.Limit
	)

	stmt, _ := db.Prepare("select `group`,quota,status,phone,limit_quota_phone,black_phone,uid from user where `key` = ?")
	stmt1, _ := db.Prepare("select content,status,`key`,sid from sign where `key` = ?")
	stmt2, _ := db.Prepare("select num,time from `limit` where sid = ?")
	err := stmt.QueryRow(key).Scan(&userinfo.Group, &userinfo.Quota, &userinfo.Status, &userinfo.Phone, &userinfo.LimitQuotaPhone, &userinfo.BlackPhone, &userinfo.Uid)
	if err != nil {
		c.JSON(200, gin.H{
			"status": -500,
			"msg":    "获取用户信息错误1",
			"bool":   false,
		})
		common.Log(c, err.Error())
		return
	}
	if !common.SignIsOwnWithSignKey(db, strconv.Itoa(userinfo.Uid), signKey) {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "无法操作不属于自己的签名头",
			"bool":   false,
		})
		return
	}
	err1 := stmt1.QueryRow(signKey).Scan(&signinfo.Content, &signinfo.Status, &signinfo.Key, &signinfo.Sid)
	if err1 != nil {
		c.JSON(200, gin.H{
			"status": -500,
			"msg":    "获取用户信息错误2",
			"bool":   false,
		})
		common.Log(c, err.Error())
		return
	}
	rows, err2 := stmt2.Query(signinfo.Sid)
	if err2 != nil {
		c.JSON(200, gin.H{
			"status": -500,
			"msg":    "获取用户信息错误3",
			"bool":   false,
		})
		common.Log(c, err.Error())
		return
	}
	if userinfo.Status != 0 || signinfo.Status != 1 {
		c.JSON(200, gin.H{
			"status": -401,
			"msg":    "用户被封禁或签名头未审核",
			"bool":   false,
		})
		return
	}
	if userinfo.Quota < 1 {
		c.JSON(200, gin.H{
			"status": 305,
			"msg":    "用户账户额度不足",
			"bool":   false,
		})
		return
	}
	if strings.Contains(userinfo.BlackPhone, to) {
		c.JSON(200, gin.H{
			"status": 307,
			"msg":    "黑名单手机号，不允许发送",
			"bool":   false,
		})
		return
	}
	for rows.Next() {
		err := rows.Scan(&limit.Num, &limit.Time)
		if err != nil {
			c.JSON(200, gin.H{
				"status": 508,
				"msg":    "处理限流操作时错误，请联系站长",
				"bool":   false,
			})
			return
		}
		var (
			formerlyStmt *sql.Stmt
			num          int
		)

		switch limit.Time {
		case "s":
			formerlyStmt, _ = db.Prepare("select count(*) from sent_log where sid = ? and time > date_sub(now(),interval 1 second )")
			break
		case "m":
			formerlyStmt, _ = db.Prepare("select count(*) from sent_log where sid = ? and time > date_sub(now(),interval 1 minute )")
			break
		case "h":
			formerlyStmt, _ = db.Prepare("select count(*) from sent_log where sid = ? and time > date_sub(now(),interval 1 hour )")
			break
		case "d":
			formerlyStmt, _ = db.Prepare("select count(*) from sent_log where sid = ? and time > date_sub(now(),interval 1 day )")
			break
		default:
			c.JSON(200, gin.H{
				"status": 320,
				"msg":    "错误的类型，请检查操作",
				"bool":   false,
			})
			return
		}
		err = formerlyStmt.QueryRow(signinfo.Sid).Scan(&num)
		if err != nil {
			c.JSON(200, gin.H{
				"status": 321,
				"msg":    "错误查询，数据库表不存在？",
				"bool":   false,
			})
			return
		}
		if num > limit.Num {
			c.JSON(200, gin.H{
				"status": 322,
				"msg":    "超过最大发信限制，请等等再发",
				"bool":   false,
			})
			return
		}
	}
	sign := content[strings.Index(content, "【") : strings.Index(content, "】")+3]
	if sign != signinfo.Content {
		c.JSON(200, gin.H{
			"status": 308,
			"msg":    "用户错误，传递签名错误",
			"bool":   false,
		})
		return
	}
	TemplatesStmt, _ := db.Prepare("select content,status from template where sid = ?")
	TemplatesRows, err := TemplatesStmt.Query(signinfo.Sid)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "获取报备模板列表失败，错误已记录",
			"bool":   false,
		})
		return
	}
	var template model.Template
	for TemplatesRows.Next() {
		err := TemplatesRows.Scan(&template.Content, &template.Status)
		if err != nil {
			common.Log(c, err.Error())
			c.JSON(200, gin.H{
				"status": -1,
				"msg":    "扫描字段出错，错误已记录",
				"bool":   false,
			})
			return
		}
		template.Content = strings.ReplaceAll(template.Content, "@", "(.+?)")
		match, err := regexp.Match(template.Content, []byte(content))
		if err != nil {
			c.JSON(200, gin.H{
				"status": -1,
				"msg":    "验证模板错误，请检查传递内容",
				"bool":   false,
			})
			return
		}
		if !match || template.Status != 1 {
			c.JSON(200, gin.H{
				"status": -1,
				"msg":    "模板不存在或模板未审核",
				"bool":   false,
			})
			return
		}
	}
	DescStmt, err := db.Prepare("insert into sent_log(uid, sid, content, phone, time, decrease_num) VALUE (?,?,?,?,now(),?)")
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "添加日志失败",
			"bool":   false,
		})
		return
	}
	DescNum := math.Ceil(float64(utf8.RuneCountInString(content)) / 64)
	result, err := DescStmt.Exec(userinfo.Uid, signinfo.Sid, content, to, DescNum)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "添加日志失败2",
			"bool":   false,
		})
		return
	}
	index, err := result.LastInsertId()
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "添加日志失败，获取下标错误",
			"bool":   false,
		})
		return
	}
	BaseValues.Add("sign", signinfo.Key)
	BaseValues.Add("to", to)
	BaseValues.Add("content", content)
	form, err := http.PostForm(BaseUrl+"/Api/Sent", BaseValues)
	if err != nil {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "构造socket错误",
			"bool":   false,
		})
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(form.Body)
	all, err := io.ReadAll(form.Body)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "读取返回错误，错误已记录",
			"bool":   false,
		})
		return
	}
	parse := gjson.Parse(string(all))
	if !parse.Get("bool").Bool() {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    parse.Get("msg").String(),
			"bool":   false,
		})
		return
	}
	stmt, _ = db.Prepare("update sent_log set task_id = ?,status = 3 where id = ?")
	_, err = stmt.Exec(parse.Get("task_id").Int(), index)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "更新发信日志失败，已记录错误",
			"bool":   false,
			"data":   string(all),
		})
		return
	}
	stmt, _ = db.Prepare("update user set quota = quota - ? where uid = ?")
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    fmt.Sprintf("更新用户余额失败，%d减少额度%f", userinfo.Uid, DescNum),
			"bool":   false,
		})
		return
	}
	_, err = stmt.Exec(DescNum, userinfo.Uid)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    fmt.Sprintf("更新用户余额失败2，%d减少额度%f", userinfo.Uid, DescNum),
			"bool":   false,
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 1,
		"msg":    parse.Get("msg").String(),
		"bool":   true,
	})
}
