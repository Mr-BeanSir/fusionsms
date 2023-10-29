package admin

import (
	"fmt"
	"fusionsms/common"
	"fusionsms/model"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetCodeList(c *gin.Context) {
	pages := c.PostForm("pages")
	db := common.GetDbInfo(c)
	sql := "select id,`code`,quota,`status`,create_time,use_time,use_uid from exchange_code order by id desc limit 30"
	sqlNum := "select count(*) from exchange_code"
	if pages != "" {
		atoi, err := strconv.Atoi(pages)
		if err != nil {
			c.JSON(200, gin.H{
				"status": -1,
				"msg":    "错误传参",
			})
			return
		}
		sql += fmt.Sprintf(" offset %d", (atoi-1)*30)
	}
	var (
		exchange  model.ExchangeCode
		exchanges []model.ExchangeCode
		num       int
	)

	rows, err := db.Query(sql)
	if err != nil {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    err.Error(),
		})
		return
	}
	for rows.Next() {
		err := rows.Scan(&exchange.Id, &exchange.Code, &exchange.Quota, &exchange.Status, &exchange.CreateTime, &exchange.UseTime, &exchange.UseUid)
		if err != nil {
			c.JSON(200, gin.H{
				"status": -1,
				"msg":    err.Error(),
			})
			return
		}
		exchanges = append(exchanges, exchange)
	}
	_ = rows.Close()
	err = db.QueryRow(sqlNum).Scan(&num)
	if err != nil {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 1,
		"msg":    "获取成功",
		"data":   exchanges,
		"num":    num,
	})
}

func AddExchangeCode(c *gin.Context) {
	num := c.PostForm("num")
	quota := c.PostForm("quota")
	if num == "" || quota == "" {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "空传参",
		})
		return
	}
	db := common.GetDbInfo(c)
	atoi, _ := strconv.Atoi(num)
	for i := 0; i < atoi; i++ {
		_, err := db.Exec("insert into exchange_code(`code`, quota, create_time) value (?,?,now())", common.RandomString(16), quota)
		if err != nil {
			c.JSON(200, gin.H{
				"status": -1,
				"msg":    err.Error(),
			})
			return
		}
	}
	c.JSON(200, gin.H{
		"status": 1,
		"msg":    "添加成功",
	})
}
