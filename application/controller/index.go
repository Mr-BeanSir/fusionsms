package controller

import (
	"database/sql"
	"fusionsms/common"
	"github.com/gin-gonic/gin"
	"time"
)

func Logout(c *gin.Context) {
	db := common.GetDbInfo(c)
	jwt := common.GetJwtInfo(c)
	stmt, _ := db.Prepare("update user set jwt = '' where uid = ?")
	result, err := stmt.Exec(jwt.Uid)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "数据库操作错误，错误已记录",
		})
		return
	}
	if num, _ := result.RowsAffected(); num == 1 {
		c.JSON(200, gin.H{
			"status": 1,
			"msg":    "退出登录成功",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": -1,
		"msg":    "系统级错误，退出登录失败",
	})
}

func DaySendNum(c *gin.Context) {
	db := common.GetDbInfo(c)
	jwt := common.GetJwtInfo(c)
	num, err := getDaySendNum(db, jwt.Uid)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "获取日发信数错误，已记录",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 1,
		"msg":    "成功",
		"num":    num,
	})
}

func DaySendSuccessNum(c *gin.Context) {
	db := common.GetDbInfo(c)
	jwt := common.GetJwtInfo(c)
	num, err := getDaySendSuccessNum(db, jwt.Uid)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "获取错误，已记录",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 1,
		"msg":    "成功",
		"num":    num,
	})
}

func DaySendErrorNum(c *gin.Context) {
	db := common.GetDbInfo(c)
	jwt := common.GetJwtInfo(c)
	num, err := getDaySendErrorNum(db, jwt.Uid)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "获取错误，已记录",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 1,
		"msg":    "成功",
		"num":    num,
	})
}

func PaddingSignNum(c *gin.Context) {
	db := common.GetDbInfo(c)
	jwt := common.GetJwtInfo(c)
	num, err := getPaddingSignNum(db, jwt.Uid)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "获取错误，已记录",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 1,
		"msg":    "成功",
		"num":    num,
	})
}

func PaddingTemplateNum(c *gin.Context) {
	db := common.GetDbInfo(c)
	jwt := common.GetJwtInfo(c)
	num, err := getPaddingTemplateNum(db, jwt.Uid)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "获取错误，已记录",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 1,
		"msg":    "成功",
		"num":    num,
	})
}

func SurplusQuotaNum(c *gin.Context) {
	db := common.GetDbInfo(c)
	jwt := common.GetJwtInfo(c)
	num, err := getSurplusQuotaNum(db, jwt.Uid)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "获取错误，已记录",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 1,
		"msg":    "成功",
		"num":    num,
	})
}

func ControllerIndex(c *gin.Context) {
	db := common.GetDbInfo(c)
	jwt := common.GetJwtInfo(c)
	surplusQuotaNum, err := getSurplusQuotaNum(db, jwt.Uid)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "获取错误，已记录",
		})
		return
	}
	paddingTemplateNum, err := getPaddingTemplateNum(db, jwt.Uid)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "获取错误，已记录",
		})
		return
	}
	paddingSignNum, err := getPaddingSignNum(db, jwt.Uid)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "获取错误，已记录",
		})
		return
	}
	daySendErrorNum, err := getDaySendErrorNum(db, jwt.Uid)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "获取错误，已记录",
		})
		return
	}
	daySendSuccessNum, err := getDaySendSuccessNum(db, jwt.Uid)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "获取错误，已记录",
		})
		return
	}
	daySendNum, err := getDaySendNum(db, jwt.Uid)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "获取日发信数错误，已记录",
		})
		return
	}
	accountKey, err := getAccountKey(db, jwt.Uid)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "获取账号key错误，已记录",
		})
		return
	}
	c.JSON(200, gin.H{
		"status":             1,
		"daySendNum":         daySendNum,
		"daySendErrorNum":    daySendErrorNum,
		"daySendSuccessNum":  daySendSuccessNum,
		"paddingSignNum":     paddingSignNum,
		"paddingTemplateNum": paddingTemplateNum,
		"surplusQuotaNum":    surplusQuotaNum,
		"accountKey":         accountKey,
	})
}

func ChartData(c *gin.Context) {
	db := common.GetDbInfo(c)
	jwt := common.GetJwtInfo(c)
	data, err := getUserSendDataChart(db, jwt.Uid)
	if err != nil {
		common.Log(c, err.Error())
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "获取错误，已记录",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 1,
		"data":   data,
	})
}

func getAccountKey(db *sql.DB, uid int) (string, error) {
	stmt, _ := db.Prepare("select `key` from user where uid = ?")
	var key string
	err := stmt.QueryRow(uid).
		Scan(&key)
	if err != nil {
		return "", err
	}
	return key, nil
}

func getDaySendNum(db *sql.DB, uid int) (int, error) {
	stmt, _ := db.Prepare("select count(*) from sent_log where uid = ? and time between ? and ?")
	var num int
	err := stmt.QueryRow(uid, common.GetBetweenTime("00", "00", "00"), common.GetBetweenTime("23", "59", "59")).
		Scan(&num)
	if err != nil {
		return -1, err
	}
	return num, nil
}

func getDaySendSuccessNum(db *sql.DB, uid int) (int, error) {
	stmt, _ := db.Prepare("select count(*) from sent_log where uid = ? and status = 1 and time between ? and ?")
	var num int
	err := stmt.QueryRow(uid, common.GetBetweenTime("00", "00", "00"), common.GetBetweenTime("23", "59", "59")).
		Scan(&num)
	//log.Println(common.GetBetweenTime("00", "00", "00"), common.GetBetweenTime("23", "59", "59"))
	if err != nil {
		return -1, err
	}
	return num, nil
}

func getDaySendErrorNum(db *sql.DB, uid int) (int, error) {
	stmt, _ := db.Prepare("select count(*) from sent_log where uid = ? and status = 2 and time between ? and ?")
	var num int
	err := stmt.QueryRow(uid, common.GetBetweenTime("00", "00", "00"), common.GetBetweenTime("23", "59", "59")).
		Scan(&num)
	if err != nil {
		return -1, err
	}
	return num, nil
}

func getPaddingSignNum(db *sql.DB, uid int) (int, error) {
	stmt, _ := db.Prepare("select count(*) from sign where uid = ? and status = 0")
	var num int
	err := stmt.QueryRow(uid).Scan(&num)
	if err != nil {
		return -1, err
	}
	return num, nil
}

func getPaddingTemplateNum(db *sql.DB, uid int) (int, error) {
	stmt, _ := db.Prepare("select count(*) from template where uid = ? and status = 0")
	var num int
	err := stmt.QueryRow(uid).Scan(&num)
	if err != nil {
		return -1, err
	}
	return num, nil
}

func getSurplusQuotaNum(db *sql.DB, uid int) (int, error) {
	stmt, _ := db.Prepare("select quota from user where uid = ?")
	var num int
	err := stmt.QueryRow(uid).Scan(&num)
	if err != nil {
		return -1, err
	}
	return num, nil
}

func getUserSendDataChart(db *sql.DB, uid int) (map[string]map[string]int, error) {
	var (
		num     int
		success int
		errors  int
	)
	nums := make(map[string]map[string]int)
	day1 := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	day2 := time.Now().AddDate(0, 0, -2).Format("2006-01-02")
	day3 := time.Now().AddDate(0, 0, -3).Format("2006-01-02")
	day4 := time.Now().AddDate(0, 0, -4).Format("2006-01-02")
	day5 := time.Now().AddDate(0, 0, -5).Format("2006-01-02")
	day6 := time.Now().AddDate(0, 0, -6).Format("2006-01-02")
	day := time.Now().Format("2006-01-02")
	sqlAll := "SELECT COUNT(*) from sent_log where uid = ? and DATE_FORMAT(time,'%y-%m-%d') = DATE_FORMAT(?,'%y-%m-%d')"
	sqlSuccess := "SELECT COUNT(*) from sent_log where uid = ? and status = 1 and DATE_FORMAT(time,'%y-%m-%d') = DATE_FORMAT(?,'%y-%m-%d')"
	sqlError := "SELECT COUNT(*) from sent_log where uid = ? and status = 2 and DATE_FORMAT(time,'%y-%m-%d') = DATE_FORMAT(?,'%y-%m-%d')"
	stmtAll, _ := db.Prepare(sqlAll)
	stmtSuccess, _ := db.Prepare(sqlSuccess)
	stmtError, _ := db.Prepare(sqlError)
	err := stmtAll.QueryRow(uid, day).Scan(&num)
	if err != nil {
		return nums, err
	}
	err = stmtSuccess.QueryRow(uid, day).Scan(&success)
	if err != nil {
		return nums, err
	}
	err = stmtError.QueryRow(uid, day).Scan(&errors)
	if err != nil {
		return nums, err
	}
	nums = setChartMap(nums, day, num, success, errors)
	{
		err := stmtAll.QueryRow(uid, day1).Scan(&num)
		if err != nil {
			return nums, err
		}
		err = stmtSuccess.QueryRow(uid, day1).Scan(&success)
		if err != nil {
			return nums, err
		}
		err = stmtError.QueryRow(uid, day1).Scan(&errors)
		if err != nil {
			return nums, err
		}
		nums = setChartMap(nums, day1, num, success, errors)
	}
	{
		err = stmtAll.QueryRow(uid, day2).Scan(&num)
		if err != nil {
			return nums, err
		}
		err = stmtSuccess.QueryRow(uid, day2).Scan(&success)
		if err != nil {
			return nums, err
		}
		err = stmtError.QueryRow(uid, day2).Scan(&errors)
		if err != nil {
			return nums, err
		}
		nums = setChartMap(nums, day2, num, success, errors)
	}
	{
		err = stmtAll.QueryRow(uid, day3).Scan(&num)
		if err != nil {
			return nums, err
		}
		err = stmtSuccess.QueryRow(uid, day3).Scan(&success)
		if err != nil {
			return nums, err
		}
		err = stmtError.QueryRow(uid, day3).Scan(&errors)
		if err != nil {
			return nums, err
		}
		nums = setChartMap(nums, day3, num, success, errors)
	}
	{
		err := stmtAll.QueryRow(uid, day4).Scan(&num)
		if err != nil {
			return nums, err
		}
		err = stmtSuccess.QueryRow(uid, day4).Scan(&success)
		if err != nil {
			return nums, err
		}
		err = stmtError.QueryRow(uid, day4).Scan(&errors)
		if err != nil {
			return nums, err
		}
		nums = setChartMap(nums, day4, num, success, errors)
	}
	{
		err = stmtAll.QueryRow(uid, day5).Scan(&num)
		if err != nil {
			return nums, err
		}
		err = stmtSuccess.QueryRow(uid, day5).Scan(&success)
		if err != nil {
			return nums, err
		}
		err = stmtError.QueryRow(uid, day5).Scan(&errors)
		if err != nil {
			return nums, err
		}
		nums = setChartMap(nums, day5, num, success, errors)
	}
	{
		err = stmtAll.QueryRow(uid, day6).Scan(&num)
		if err != nil {
			return nums, err
		}
		err = stmtSuccess.QueryRow(uid, day6).Scan(&success)
		if err != nil {
			return nums, err
		}
		err = stmtError.QueryRow(uid, day6).Scan(&errors)
		if err != nil {
			return nums, err
		}
		nums = setChartMap(nums, day6, num, success, errors)
	}
	return nums, nil
}

func setChartMap(data map[string]map[string]int, date string, all, success, errors int) map[string]map[string]int {
	datas := map[string]int{
		"all":     all,
		"success": success,
		"error":   errors,
	}
	data[date] = datas
	return data
}
