package controller

import (
	"fusionsms/common"
	"fusionsms/model"
	"github.com/gin-gonic/gin"
	"regexp"
	"strconv"
	"strings"
)

func RemainingLimitPrompt(c *gin.Context) {
	db := common.GetDbInfo(c)
	jwt := common.GetJwtInfo(c)
	stmt, _ := db.Prepare("select phone,limit_quota_phone from user where uid = ?")
	var user model.User
	err := stmt.QueryRow(jwt.Uid).Scan(&user.Phone, &user.LimitQuotaPhone)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "获取余额提示失败，已记录错误",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 1,
		"phone":  user.Phone,
		"num":    user.LimitQuotaPhone,
		"msg":    "获取成功",
	})
}

func RemainingLimitPromptSave(c *gin.Context) {
	num := c.PostForm("num")
	phone := c.PostForm("phone")
	if num == "" || phone == "" {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "空传参",
		})
		return
	}
	match, _ := regexp.Match("^[0-9]{11}$", []byte(phone))
	match1, _ := regexp.Match("^[0-9]+$", []byte(num))
	if !match || !match1 {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "错误传参",
		})
		return
	}
	db := common.GetDbInfo(c)
	jwt := common.GetJwtInfo(c)
	stmt, _ := db.Prepare("update user set phone = ?,limit_quota_phone = ? where uid = ?")
	result, err := stmt.Exec(phone, num, jwt.Uid)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "获取余额提示失败，已记录错误",
		})
		return
	}
	if i, _ := result.RowsAffected(); i == 1 {
		c.JSON(200, gin.H{
			"status": 1,
			"msg":    "修改成功",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": -1,
		"msg":    "修改失败",
	})
}

func WhiteList(c *gin.Context) {
	db := common.GetDbInfo(c)
	jwt := common.GetJwtInfo(c)
	stmt, _ := db.Prepare("select white_ip from user where uid = ?")
	var user model.User
	err := stmt.QueryRow(jwt.Uid).Scan(&user.WhiteIP)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "获取白名单IP失败，已记录错误",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 1,
		"data":   user.WhiteIP,
		"msg":    "获取成功",
	})
}

func WhiteListSave(c *gin.Context) {
	ip := c.PostForm("list")
	match, _ := regexp.Match("^([0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3},?)+$", []byte(ip))
	if ip == "" || !match {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "错误传参",
		})
		return
	}
	db := common.GetDbInfo(c)
	jwt := common.GetJwtInfo(c)
	stmt, _ := db.Prepare("update user set white_ip = ? where uid = ?")
	result, err := stmt.Exec(ip, jwt.Uid)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "获取白名单IP失败，已记录错误",
		})
		return
	}
	if num, _ := result.RowsAffected(); num == 1 {
		c.JSON(200, gin.H{
			"status": 1,
			"msg":    "修改成功",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": -1,
		"msg":    "修改失败",
	})
}

func BlackPhone(c *gin.Context) {
	db := common.GetDbInfo(c)
	jwt := common.GetJwtInfo(c)
	stmt, _ := db.Prepare("select black_phone from user where uid = ?")
	var user model.User
	err := stmt.QueryRow(jwt.Uid).Scan(&user.BlackPhone)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "获取黑名单手机号失败，已记录错误",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 1,
		"data":   user.BlackPhone,
		"msg":    "获取成功",
	})
}

func BlackPhoneSave(c *gin.Context) {
	phone := c.PostForm("phone")
	match, _ := regexp.Match("^([0-9]{11},?)+$", []byte(phone))
	if phone == "" || !match {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "错误传参",
		})
		return
	}
	db := common.GetDbInfo(c)
	jwt := common.GetJwtInfo(c)
	stmt, _ := db.Prepare("update user set black_phone = ? where uid = ?")
	result, err := stmt.Exec(phone, jwt.Uid)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "保存黑名单手机号失败，已记录错误",
		})
		return
	}
	if num, _ := result.RowsAffected(); num == 1 {
		c.JSON(200, gin.H{
			"status": 1,
			"msg":    "修改成功",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": -1,
		"msg":    "修改失败",
	})
}

func Limit(c *gin.Context) {
	db := common.GetDbInfo(c)
	jwt := common.GetJwtInfo(c)
	stmt, _ := db.Prepare("select id,uid,sid, num, time from `limit` where uid = ?")
	var (
		limit  model.Limit
		limits []model.Limit
	)
	rows, err := stmt.Query(jwt.Uid)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "获取限流列表失败，已记录错误",
		})
		return
	}
	for rows.Next() {
		err := rows.Scan(&limit.ID, &limit.Uid, &limit.Sid, &limit.Num, &limit.Time)
		if err != nil {
			common.Log(c, err.Error())
			c.JSON(200, gin.H{
				"status": -1,
				"msg":    "获取限流列表失败，已记录错误",
			})
			return
		}
		limits = append(limits, limit)
	}
	c.JSON(200, gin.H{
		"status": 1,
		"data":   limits,
		"msg":    "获取成功",
	})
}

func AddLimit(c *gin.Context) {
	sid := c.PostForm("sid")
	num := c.PostForm("num")
	ty_pe := c.PostForm("type")
	if sid == "" || num == "" || !strings.Contains("s,m,h,d", ty_pe) {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "错误传参",
		})
		return
	}
	db := common.GetDbInfo(c)
	jwt := common.GetJwtInfo(c)
	if !common.SignIsOwn(db, strconv.Itoa(jwt.Uid), sid) {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "无法操作不属于自己的签名头",
		})
		return
	}
	stmt, _ := db.Prepare("insert into `limit`(sid, uid, num, time) value (?,?,?,?)")
	result, err := stmt.Exec(sid, jwt.Uid, num, ty_pe)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "添加限流规则失败，已记录错误",
		})
		return
	}
	if num, _ := result.RowsAffected(); num == 1 {
		c.JSON(200, gin.H{
			"status": 1,
			"msg":    "添加成功",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": -1,
		"msg":    "添加失败",
	})
}

func EditLimit(c *gin.Context) {
	id := c.PostForm("id")
	sid := c.PostForm("sid")
	num := c.PostForm("num")
	ty_pe := c.PostForm("type")
	if sid == "" || id == "" || num == "" || !strings.Contains("s,m,h,d", ty_pe) {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "错误传参",
		})
		return
	}
	db := common.GetDbInfo(c)
	jwt := common.GetJwtInfo(c)
	if !common.SignIsOwn(db, strconv.Itoa(jwt.Uid), sid) {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "无法操作不属于自己的签名头",
		})
		return
	}
	stmt, _ := db.Prepare("update `limit` set sid = ?,num = ?,time = ? where id = ?")
	result, err := stmt.Exec(sid, num, ty_pe, id)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "修改限流规则失败，已记录错误",
		})
		return
	}
	if num, _ := result.RowsAffected(); num == 1 {
		c.JSON(200, gin.H{
			"status": 1,
			"msg":    "修改成功",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": -1,
		"msg":    "修改失败",
	})
}

func DeleteLimit(c *gin.Context) {
	id := c.PostForm("id")
	if id == "" {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "错误传参",
		})
		return
	}
	db := common.GetDbInfo(c)
	jwt := common.GetJwtInfo(c)
	if !common.LimitIsOwn(db, strconv.Itoa(jwt.Uid), id) {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "无法操作不属于自己的规则",
		})
		return
	}
	stmt, _ := db.Prepare("delete from `limit` where id = ?")
	result, err := stmt.Exec(id)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "删除限流规则失败，已记录错误",
		})
		return
	}
	if num, _ := result.RowsAffected(); num == 1 {
		c.JSON(200, gin.H{
			"status": 1,
			"msg":    "删除成功",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": -1,
		"msg":    "删除失败",
	})
}
