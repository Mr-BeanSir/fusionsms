package controller

import (
	"fusionsms/common"
	"fusionsms/model"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"strings"
)

func AddSign(c *gin.Context) {
	content := c.PostForm("content")
	if content == "" {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "签名头内容为空",
		})
		return
	}
	jwt := common.GetJwtInfo(c)
	db := common.GetDbInfo(c)
	defer db.Close()
	stmt, _ := db.Prepare("select count(*) from sign where content = ?")
	var num int
	stmt.QueryRow(content).Scan(&num)
	if num != 0 {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "此签名头数据库中已存在",
		})
		return
	}
	signID, signKey, md5, err := ApiAddSign(content)
	if err != nil {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    err.Error(),
		})
		return
	}
	stmt, _ = db.Prepare("insert into sign(uid, content, `key`,`super_id`,`md5`) value(?,?,?,?,?)")
	result, err := stmt.Exec(jwt.Uid, content, signKey, signID, md5)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "添加签名头失败，原因已记录请联系站长",
		})
		return
	}
	if rows, err := result.RowsAffected(); rows != 1 {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "添加签名头失败，原因已记录请联系站长",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 1,
		"msg":    "添加签名头成功！",
	})
}

func GetSignList(c *gin.Context) {
	jwt := common.GetJwtInfo(c)
	db := common.GetDbInfo(c)
	stmt, _ := db.Prepare("select sid,content,`key`,status from sign where uid = ?")
	var (
		sign     model.Sign
		signList []model.Sign
	)
	rows, err := stmt.Query(jwt.Uid)
	if err != nil {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "获取签名列表失败",
		})
		return
	}
	for rows.Next() {
		sign = model.Sign{}
		err := rows.Scan(&sign.Sid, &sign.Content, &sign.Key, &sign.Status)
		if err != nil {
			common.Log(c, err.Error())
			c.JSON(200, gin.H{
				"status": -1,
				"msg":    "获取签名列表失败，原因已记录，请联系站长",
			})
			return
		}
		signList = append(signList, sign)
	}
	c.JSON(200, gin.H{
		"status": 1,
		"data":   signList,
		"msg":    "获取成功",
	})
}

func ResetKey(c *gin.Context) {
	sign := c.PostForm("sign")
	jwt := common.GetJwtInfo(c)
	db := common.GetDbInfo(c)
	if !common.SignIsOwnWithSignKey(db, strconv.Itoa(jwt.Uid), sign) {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "你无法操作不属于你的签名",
		})
		return
	}
	stmt, err := db.Prepare("update sign set `key` = ?, md5 = ? where `key` = ?")
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "重置密钥失败，错误日志已记录，请联系站长",
		})
		return
	}
	superIdStmt, err := db.Prepare("select super_id from sign where `key` = ?")
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "重置密钥失败，错误日志已记录，请联系站长",
		})
		return
	}
	var superId string
	superIdStmt.QueryRow(sign).Scan(&superId)
	resetSign, md5, err := ApiResetSign(superId)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    err.Error(),
		})
		return
	}
	result, err := stmt.Exec(resetSign, md5, sign)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "重置密钥失败，错误日志已记录，请联系站长",
		})
		return
	}
	if num, _ := result.RowsAffected(); num == 1 {
		c.JSON(200, gin.H{
			"status": 1,
			"msg":    "重置密钥成功",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": -1,
		"msg":    "重置密钥失败",
	})
}

func GetSignContent(c *gin.Context) {
	sid := c.Param("sid")
	db := common.GetDbInfo(c)
	jwt := common.GetJwtInfo(c)
	if !common.SignIsOwn(db, strconv.Itoa(jwt.Uid), sid) {
		c.JSON(200, gin.H{
			"status": -2,
			"msg":    "你无法操作不属于你的签名",
		})
		return
	}
	stmt, _ := db.Prepare("select content from sign where sid = ?")
	var content string
	stmt.QueryRow(sid).Scan(&content)
	c.JSON(200, gin.H{
		"status":  1,
		"content": content,
		"msg":     "获取成功",
	})
}

func GetSign(c *gin.Context) {
	sid := c.Param("sid")
	db := common.GetDbInfo(c)
	jwt := common.GetJwtInfo(c)
	if !common.SignIsOwn(db, strconv.Itoa(jwt.Uid), sid) {
		c.JSON(200, gin.H{
			"status": -2,
			"msg":    "你无法操作不属于你的签名",
		})
		return
	}
	stmt, _ := db.Prepare("select tid,content,status,reason,super_tid from template where sid = ?")
	rows, err := stmt.Query(sid)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "获取签名模板失败，错误日志已记录，请联系站长",
		})
		return
	}
	var (
		template  model.Template
		templates []model.Template
	)

	for rows.Next() {
		err := rows.Scan(&template.Tid, &template.Content, &template.Status, &template.Reason)
		if err != nil {
			common.Log(c, err.Error())
			c.JSON(200, gin.H{
				"status": -1,
				"msg":    "获取签名模板失败，错误日志已记录，请联系站长",
			})
			return
		}
		templates = append(templates, template)
	}
	c.JSON(200, gin.H{
		"status": 1,
		"data":   templates,
		"msg":    "获取签名模板成功",
	})
}

func AddSignTemplate(c *gin.Context) {
	sid := c.PostForm("sid")
	template := c.PostForm("content")
	templateSplit := strings.Split(template, "\n")
	if template == "" || len(templateSplit) < 1 {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "报备模板传递为空",
		})
		return
	}
	db := common.GetDbInfo(c)
	jwt := common.GetJwtInfo(c)
	if !common.SignIsOwn(db, strconv.Itoa(jwt.Uid), sid) {
		c.JSON(200, gin.H{
			"status": -2,
			"msg":    "你无法操作不属于你的签名",
		})
		return
	}
	stmt, _ := db.Prepare("select `super_id`,content from sign where sid = ?")
	var sign model.Sign
	stmt.QueryRow(sid).Scan(&sign.SuperId, &sign.Content)
	i := 0
	for _, v := range templateSplit {
		super_tid, err := ApiAddTemplate(strconv.Itoa(sign.SuperId), v)
		if err != nil {
			c.JSON(200, gin.H{
				"status": -1,
				"msg":    err.Error(),
			})
			return
		}
		stmt, _ := db.Prepare("insert into template(sid, uid, content,super_tid) value(?,?,?,?)")
		_, err = stmt.Exec(sid, jwt.Uid, sign.Content+v, super_tid[0].Int())
		if err != nil {
			common.Log(c, err.Error())
			c.JSON(200, gin.H{
				"status": -1,
				"msg":    "添加签名模板失败，错误日志已记录，请联系站长",
			})
			return
		}
		i++
	}
	c.JSON(200, gin.H{
		"status": 1,
		"msg":    "添加成功",
	})
}

func DeleteTemplate(c *gin.Context) {
	tid := c.PostForm("tid")
	if tid == "" {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "传递参数为空",
		})
		return
	}
	db := common.GetDbInfo(c)
	jwt := common.GetJwtInfo(c)
	if !common.TemplateIsOwn(db, strconv.Itoa(jwt.Uid), tid) {
		c.JSON(200, gin.H{
			"status": -2,
			"msg":    "你无法操作不属于你的签名",
		})
		return
	}
	var template model.Template
	var sign model.Sign
	stmt, _ := db.Prepare("select sid,super_tid from template where tid = ?")
	stmt.QueryRow(tid).Scan(&template.Sid, &template.SuperTid)
	stmt, _ = db.Prepare("select super_id from sign where sid = ?")
	stmt.QueryRow(template.Sid).Scan(&sign.SuperId)
	log.Println(sign.SuperId, template.SuperTid)
	err := ApiDeleteTemplate(strconv.Itoa(sign.SuperId), template.SuperTid)
	if err != nil {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    err.Error(),
		})
		return
	}
	stmt, _ = db.Prepare("delete from template where tid = ?")
	result, err := stmt.Exec(tid)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "删除签名模板失败，错误日志已记录，请联系站长",
		})
		return
	}
	if num, _ := result.RowsAffected(); num != 1 {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "删除签名模板失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 1,
		"msg":    "删除成功",
	})
}
