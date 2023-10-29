package controller

import (
	"fusionsms/common"
	"fusionsms/model"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetLogList(c *gin.Context) {
	start := c.PostForm("start")
	end := c.PostForm("end")
	content := c.PostForm("content")
	taskID := c.PostForm("task_id")
	ty_pe := c.PostForm("type")
	pages := c.PostForm("pages")
	if start == "" || end == "" || ty_pe == "" {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "传递参数有空",
		})
		return
	}
	db := common.GetDbInfo(c)
	jwt := common.GetJwtInfo(c)
	var args []any
	args = append(args, jwt.Uid, start, end)
	sql := "select id, sid, content, phone, time, task_id, receive, receive_time, decrease_num, status from sent_log where uid = ? and time between ? and ? "
	count := "select count(*) from sent_log where uid = ? and time between ? and ? "
	if ty_pe != "all" {
		sql += "and status = ?"
		count += "and status = ?"
		args = append(args, ty_pe)
	}
	if content != "" {
		sql += " and `content` like CONCAT('%',?,'%')"
		count += " and `content` like CONCAT('%',?,'%')"
		args = append(args, content)
	}
	if taskID != "" {
		sql += " and task_id = ?"
		count += " and task_id = ?"
		args = append(args, taskID)
	}
	sql += " order by id desc limit 50"
	count_args := args
	if pages != "" {
		atoi, err := strconv.Atoi(pages)
		if err != nil {
			c.JSON(200, gin.H{
				"status": -1,
				"msg":    "传递参数错误",
			})
			return
		}
		sql += " offset ?"
		args = append(args, (atoi-1)*50)
	}
	stmt, err := db.Prepare(sql)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "获取日志列表失败1，已记录错误",
		})
		return
	}
	rows, err := stmt.Query(args...)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "获取日志列表失败2，已记录错误",
		})
		return
	}
	stmt1, _ := db.Prepare(count)
	var (
		sentLog  model.SentLogModel
		sentLogs []model.SentLogModel
		num      int
	)
	for rows.Next() {
		err := rows.Scan(&sentLog.Id, &sentLog.Sid, &sentLog.Content, &sentLog.Phone, &sentLog.Time, &sentLog.TaskId, &sentLog.Receive, &sentLog.ReceiveTime, &sentLog.DecreaseNum, &sentLog.Status)
		if err != nil {
			common.Log(c, err.Error())
			c.JSON(200, gin.H{
				"status": -1,
				"msg":    "获取日志列表失败3，已记录错误",
			})
			return
		}
		sentLogs = append(sentLogs, sentLog)
	}
	err = stmt1.QueryRow(count_args...).Scan(&num)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "获取日志列表失败4，已记录错误",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 1,
		"msg":    "获取成功",
		"data":   sentLogs,
		"num":    num,
	})
}

//func GetLogList(c *gin.Context) {
//	start := c.PostForm("start")
//	end := c.PostForm("end")
//	pages := c.PostForm("pages")
//	content := c.PostForm("content")
//	taskID := c.PostForm("task_id")
//	ty_pe := c.PostForm("type")
//	if start == "" || end == "" {
//		c.JSON(200, gin.H{
//			"status": -1,
//			"msg":    "传递参数错误",
//		})
//	}
//	db := common.GetDbInfo(c)
//	jwt := common.GetJwtInfo(c)
//	sql := "select id, sid, content, phone, time, task_id, receive, receive_time, decrease_num, status from sent_log where uid = ? and time between ? and ? order by id desc limit 50"
//	var args []any
//	args = append(args, jwt.Uid, start, end)
//	if pages != "" {
//		atoi, err := strconv.Atoi(pages)
//		if err != nil {
//			c.JSON(200, gin.H{
//				"status": -1,
//				"msg":    "传递参数错误",
//			})
//			return
//		}
//		sql += " offset ?"
//		args = append(args, (atoi-1)*50)
//	}
//
//	stmt, _ := db.Prepare(sql)
//	rows, err := stmt.Query(args...)
//	stmt1, _ := db.Prepare("select count(*) from sent_log where uid = ? and time between ? and ?")
//	if err != nil {
//		common.Log(c, err.Error())
//		c.JSON(200, gin.H{
//			"status": -1,
//			"msg":    "获取日志列表失败，已记录错误",
//		})
//		return
//	}
//	var (
//		sentLog  model.SentLogModel
//		sentLogs []model.SentLogModel
//		num      int
//	)
//	for rows.Next() {
//		err := rows.Scan(&sentLog.Id, &sentLog.Sid, &sentLog.Content, &sentLog.Phone, &sentLog.Time, &sentLog.TaskId, &sentLog.Receive, &sentLog.ReceiveTime, &sentLog.DecreaseNum, &sentLog.Status)
//		if err != nil {
//			common.Log(c, err.Error())
//			c.JSON(200, gin.H{
//				"status": -1,
//				"msg":    "获取日志列表失败，已记录错误",
//			})
//			return
//		}
//		sentLogs = append(sentLogs, sentLog)
//	}
//	err = stmt1.QueryRow(jwt.Uid, start, end).Scan(&num)
//	if err != nil {
//		common.Log(c, err.Error())
//		c.JSON(200, gin.H{
//			"status": -1,
//			"msg":    "获取日志列表失败，已记录错误",
//		})
//		return
//	}
//	c.JSON(200, gin.H{
//		"status": 1,
//		"msg":    "获取成功",
//		"data":   sentLogs,
//		"num":    num,
//	})
//}

//func FilterLogList(c *gin.Context) {
//	start := c.PostForm("start")
//	end := c.PostForm("end")
//	content := c.PostForm("content")
//	taskID := c.PostForm("task_id")
//	ty_pe := c.PostForm("type")
//	if start == "" || end == "" || ty_pe == "" {
//		c.JSON(200, gin.H{
//			"status": -1,
//			"msg":    "传递参数有空",
//		})
//		return
//	}
//	db := common.GetDbInfo(c)
//	jwt := common.GetJwtInfo(c)
//	var args []any
//	args = append(args, jwt.Uid, start, end)
//	sql := "select id, sid, content, phone, time, task_id, receive, receive_time, decrease_num, status from sent_log where uid = ? and time between ? and ? "
//	count := "select count(*) from sent_log where uid = ? and time between ? and ? "
//	if ty_pe != "all" {
//		sql += "and status = ?"
//		count += "and status = ?"
//		args = append(args, ty_pe)
//	}
//	if content != "" {
//		sql += " and `content` like CONCAT('%',?,'%')"
//		count += " and `content` like CONCAT('%',?,'%')"
//		args = append(args, content)
//	}
//	if taskID != "" {
//		sql += " and task_id = ?"
//		count += " and task_id = ?"
//		args = append(args, taskID)
//	}
//	sql += " order by id desc"
//	stmt, err := db.Prepare(sql)
//	if err != nil {
//		common.Log(c, err.Error())
//		c.JSON(200, gin.H{
//			"status": -1,
//			"msg":    "获取日志列表失败，已记录错误",
//		})
//		return
//	}
//	rows, err := stmt.Query(args...)
//	if err != nil {
//		common.Log(c, err.Error())
//		c.JSON(200, gin.H{
//			"status": -1,
//			"msg":    "获取日志列表失败，已记录错误",
//		})
//		return
//	}
//	stmt1, _ := db.Prepare(count)
//	var (
//		sentLog  model.SentLogModel
//		sentLogs []model.SentLogModel
//		num      int
//	)
//	for rows.Next() {
//		err := rows.Scan(&sentLog.Id, &sentLog.Sid, &sentLog.Content, &sentLog.Phone, &sentLog.Time, &sentLog.TaskId, &sentLog.Receive, &sentLog.ReceiveTime, &sentLog.DecreaseNum, &sentLog.Status)
//		if err != nil {
//			common.Log(c, err.Error())
//			c.JSON(200, gin.H{
//				"status": -1,
//				"msg":    "获取日志列表失败，已记录错误",
//			})
//			return
//		}
//		sentLogs = append(sentLogs, sentLog)
//	}
//	err = stmt1.QueryRow(args...).Scan(&num)
//	if err != nil {
//		common.Log(c, err.Error())
//		c.JSON(200, gin.H{
//			"status": -1,
//			"msg":    "获取日志列表失败，已记录错误",
//		})
//		return
//	}
//	c.JSON(200, gin.H{
//		"status": 1,
//		"msg":    "获取成功",
//		"data":   sentLogs,
//		"num":    num,
//	})
//}
