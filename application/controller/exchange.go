package controller

import (
	"fusionsms/common"
	"fusionsms/model"
	"github.com/gin-gonic/gin"
	"reflect"
	"strconv"
	"time"
)

func GetExchangeList(c *gin.Context) {
	db := common.GetDbInfo(c)
	jwt := common.GetJwtInfo(c)
	var exchange model.ExchangeLog
	var exchanges []model.ExchangeLog
	rows, err := db.Query("select id, uid, content, `time` from exchange_log where uid = ? order by id desc ", jwt.Uid)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "获取兑换记录失败1",
		})
		return
	}
	for rows.Next() {
		err := rows.Scan(&exchange.Id, &exchange.Uid, &exchange.Content, &exchange.Time)
		if err != nil {
			common.Log(c, err.Error())
			c.JSON(200, gin.H{
				"status": -1,
				"msg":    "获取兑换记录失败2",
			})
			return
		}
		exchanges = append(exchanges, exchange)
	}
	c.JSON(200, gin.H{
		"status": 1,
		"msg":    "获取成功",
		"data":   exchanges,
	})
}

func ExchangeCode(c *gin.Context) {
	code := c.PostForm("code")
	if code == "" {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "空传参",
		})
		return
	}
	db := common.GetDbInfo(c)
	jwt := common.GetJwtInfo(c)
	var exchange model.ExchangeCode
	err := db.QueryRow("select quota,`status` from exchange_code where code=?", code).Scan(&exchange.Quota, &exchange.Status)
	if reflect.DeepEqual(exchange, model.ExchangeCode{}) || err != nil {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "不存在的兑换码或数据库有误1",
		})
		return
	}
	if exchange.Status != 0 {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "兑换码已使用",
		})
		return
	}
	_, err = db.Exec("update exchange_code set `status` = 1,use_time = NOW(),use_uid = ? where code = ?;", jwt.Uid, code)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "兑换码使用失败，错误已记录",
		})
		return
	}
	_, err = db.Exec("insert into exchange_log(uid, content, `time`) value (?,?,?)", jwt.Uid, "兑换额度:"+strconv.Itoa(exchange.Quota), time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "插入记录失败，错误已记录",
		})
	}
	_, err = db.Exec("update user set quota = quota + ? where uid = ?;", exchange.Quota, jwt.Uid)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "兑换码使用失败，错误已记录",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 1,
		"msg":    "兑换成功",
	})
}
