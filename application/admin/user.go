package admin

import (
	"database/sql"
	"fusionsms/common"
	"fusionsms/model"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func List(c *gin.Context) {
	uid := c.PostForm("uid")
	username := c.PostForm("username")
	status := c.PostForm("status")
	pages := c.PostForm("pages")
	if status == "" {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "传递参数有空",
		})
		return
	}
	db := common.GetDbInfo(c)
	sql := "select uid, username, quota, `status` from user "
	sqlNum := "select count(*) from user "
	var args []any
	var argsNum []any
	if uid != "" || username != "" || status != "all" {
		sql += "where"
		sqlNum += "where"
	}
	if uid != "" {
		sql += "uid = ? "
		sqlNum += "uid = ? "
		if username != "" {
			sql += "and"
			sqlNum += "and"
		}
		args = append(args, uid)
		argsNum = append(argsNum, uid)
	}
	if username != "" {
		sql += "username like CONCAT('%',?,'%')"
		sqlNum += "username like CONCAT('%',?,'%')"
		if username != "" {
			sql += "and"
			sqlNum += "and"
		}
		args = append(args, username)
		argsNum = append(argsNum, username)
	}
	if status != "all" {
		sql += " status = ?"
		sqlNum += " status = ?"
		args = append(args, status)
		argsNum = append(argsNum, status)
	}
	sql += " order by uid desc limit 30"
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
		args = append(args, (atoi-1)*30)
	}

	stmt, _ := db.Prepare(sql)
	rows, err := stmt.Query(args...)
	stmt1, _ := db.Prepare(sqlNum)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "获取用户列表失败1，已记录错误",
		})
		return
	}
	var (
		user  model.User
		users []model.User
		num   int
	)
	for rows.Next() {
		err := rows.Scan(&user.Uid, &user.Username, &user.Quota, &user.Status)
		if err != nil {
			common.Log(c, err.Error())
			c.JSON(200, gin.H{
				"status": -1,
				"msg":    "获取用户列表失败2，已记录错误",
			})
			return
		}
		users = append(users, user)
	}
	err = stmt1.QueryRow(argsNum...).Scan(&num)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "获取用户列表失败3，已记录错误",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 1,
		"msg":    "获取成功",
		"data":   users,
		"num":    num,
	})
}

func GetDetail(c *gin.Context) {
	uid := c.Param("uid")
	if uid == "" {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "传参错误",
		})
		return
	}
	db := common.GetDbInfo(c)
	stmt, _ := db.Prepare("select username,white_ip,black_phone,`status` from user where uid = ?")
	var user model.User
	err := stmt.QueryRow(uid).Scan(&user.Username, &user.WhiteIP, &user.BlackPhone, &user.Status)
	if err != nil {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 1,
		"data":   user,
	})
}

func ChangeDetail(c *gin.Context) {
	uid := c.Param("uid")
	username := c.PostForm("username")
	password := c.PostForm("password")
	blackPhone := c.PostForm("black_phone")
	whiteIp := c.PostForm("white_ip")
	status := c.PostForm("status")
	if uid == "" || username == "" || status == "" {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "传参错误",
		})
		return
	}
	db := common.GetDbInfo(c)
	if password != "" {
		stmt, _ := db.Prepare("update user set username = ?,black_phone = ?,white_ip = ?,`status` = ?,`password` = ? where uid = ?")
		_, err := stmt.Exec(username, blackPhone, whiteIp, status, password, uid)
		if err != nil {
			c.JSON(200, gin.H{
				"status": -1,
				"msg":    err.Error(),
			})
			return
		}
	} else {
		stmt, _ := db.Prepare("update user set username = ?,black_phone = ?,white_ip = ?,`status` = ? where uid = ?")
		result, err := stmt.Exec(username, blackPhone, whiteIp, status, uid)
		if err != nil {
			c.JSON(200, gin.H{
				"status": -1,
				"msg":    err.Error(),
			})
			return
		}
		log.Println(result.RowsAffected())
	}
	c.JSON(200, gin.H{
		"status": 1,
		"msg":    "更新用户信息成功",
	})
}

func ChangeBalance(c *gin.Context) {
	uid := c.PostForm("uid")
	num := c.PostForm("num")
	if uid == "" || num == "" {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "传递参数有空",
		})
		return
	}
	db := common.GetDbInfo(c)
	if !changeBalance(db, uid, num) {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "修改UID：" + uid + "用户余额失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 1,
		"msg":    "修改UID：" + uid + "用户余额成功",
	})
}

func changeBalance(db *sql.DB, uid, num string) bool {
	stmt, _ := db.Prepare("update user set quota = quota + ? where uid = ?")
	result, err := stmt.Exec(num, uid)
	if err != nil {
		return false
	}
	_, err = result.RowsAffected()
	if err != nil {
		return false
	}
	return true
}
