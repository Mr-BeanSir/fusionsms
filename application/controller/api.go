package controller

import (
	"fusionsms/common"
	"fusionsms/config"
	"fusionsms/model"
	"github.com/gin-gonic/gin"
)

func ReceiveSignStatus(c *gin.Context) {
	key := c.PostForm("key")
	if key != common.Md5V(config.Key+config.ApiKey) {
		c.String(200, "error")
		return
	}
	sign := c.PostForm("sign")
	status := c.PostForm("status")
	if sign == "" || status == "" {
		c.String(200, "error1")
		return
	}
	db := common.GetDbInfo(c)
	stmt, _ := db.Prepare("update sign set status = ? where `key` = ?")
	result, err := stmt.Exec(status, sign)
	if err != nil {
		c.String(200, "error2")
		return
	}
	if num, _ := result.RowsAffected(); num == 1 {
		c.String(200, "ok")
		return
	}
	c.String(200, "error3")
}

func ReceiveTemplateStatus(c *gin.Context) {
	key := c.PostForm("key")
	if key != common.Md5V(config.Key+config.ApiKey) {
		c.String(200, "error")
		return
	}
	super_tid := c.PostForm("super_tid")
	status := c.PostForm("status")
	reason := c.PostForm("reason")
	if super_tid == "" || status == "" {
		c.String(200, "error1")
		return
	}
	db := common.GetDbInfo(c)
	stmt, _ := db.Prepare("update template set status = ?,reason = ? where super_tid = ?")
	result, err := stmt.Exec(status, reason, super_tid)
	if err != nil {
		c.String(200, "error2")
		return
	}
	if num, _ := result.RowsAffected(); num == 1 {
		c.String(200, "ok")
		return
	}
	c.String(200, "error3")
}

func ReceiveSentStatus(c *gin.Context) {
	key := c.PostForm("key")
	if key != common.Md5V(config.Key+config.ApiKey) {
		c.String(200, "error")
		return
	}
	taskId := c.PostForm("task_id")
	receive := c.PostForm("receive")
	receiveTime := c.PostForm("receive_time")
	boolean := c.PostForm("bool")
	if taskId == "" || boolean == "" {
		c.String(200, "error1")
		return
	}
	db := common.GetDbInfo(c)
	if boolean == "true" {
		stmt, _ := db.Prepare("update sent_log set status = 1,receive = ?,receive_time=? where task_id = ?")
		result, err := stmt.Exec(receive, receiveTime, taskId)
		if err != nil {
			c.String(200, "error2")
			return
		}
		if num, _ := result.RowsAffected(); num == 1 {
			c.String(200, "ok")
			return
		}
	} else {
		var log model.SentLogModel
		stmt, _ := db.Prepare("select decrease_num,uid from sent_log where task_id = ?")
		err := stmt.QueryRow(taskId).Scan(&log.DecreaseNum, &log.Uid)
		if err != nil {
			c.String(200, err.Error())
			return
		}
		stmt, _ = db.Prepare("update sent_log set status = 2,receive = ?,receive_time=?,decrease_num=0 where task_id = ?")
		result, err := stmt.Exec(receive, receiveTime, taskId)
		if err != nil {
			c.String(200, err.Error())
			return
		}
		if num, _ := result.RowsAffected(); num != 1 {
			c.String(200, err.Error())
			return
		}
		stmt, _ = db.Prepare("update user set quota = quota + ? where uid = ?")
		result, err = stmt.Exec(log.DecreaseNum, log.Uid)
		if err != nil {
			c.String(200, err.Error())
			return
		}
		if num, _ := result.RowsAffected(); num == 1 {
			c.String(200, "ok")
			return
		}
	}

	c.String(200, "error3")
}
