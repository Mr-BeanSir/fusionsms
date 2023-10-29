package common

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fusionsms/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"math/rand"
	"time"
)

var (
	JwtSecret = "test"
)

func Md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func RandomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func CreateJwt(claims *model.UserClaims) string {
	signingString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(JwtSecret))
	if err != nil {
		return "error"
	}
	return signingString
}

func ParseJwt(tokenString string) (*jwt.Token, error) {
	parse, err := jwt.ParseWithClaims(tokenString, &model.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtSecret), nil
	})
	if err != nil {
		return &jwt.Token{}, err
	}
	return parse, nil
}

func Log(c *gin.Context, msg string) error {
	jwts := GetJwtInfo(c)
	db := GetDbInfo(c)
	stmt, _ := db.Prepare("insert into log(uid,error,time) values(?,?,?)")
	_, err := stmt.Exec(jwts.Uid, msg, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return err
	}
	return nil
}

func GetJwtInfo(c *gin.Context) *model.UserClaims {
	tempjwta, _ := c.Get("claims")
	tempjwt, _ := tempjwta.(*jwt.Token)
	jwts, _ := tempjwt.Claims.(*model.UserClaims)
	return jwts
}

func GetDbInfo(c *gin.Context) *sql.DB {
	temp, _ := c.Get("db")
	db, _ := temp.(*sql.DB)
	return db
}

func SignIsOwn(db *sql.DB, uid, sid string) bool {
	stmt, _ := db.Prepare("select count(*) from sign where sid = ? and uid = ?")
	var num int
	err := stmt.QueryRow(sid, uid).Scan(&num)
	if err != nil {
		return false
	}
	if num == 1 {
		return true
	}
	return false
}

func SignIsOwnWithSignKey(db *sql.DB, uid, sign string) bool {
	stmt, _ := db.Prepare("select count(*) from sign where `key` = ? and uid = ?")
	var num int
	err := stmt.QueryRow(sign, uid).Scan(&num)
	if err != nil {
		return false
	}
	if num == 1 {
		return true
	}
	return false
}

func TemplateIsOwn(db *sql.DB, uid, tid string) bool {
	stmt, _ := db.Prepare("select count(*) from template where uid = ? and tid = ?")
	var num int
	err := stmt.QueryRow(uid, tid).Scan(&num)
	if err != nil {
		return false
	}
	if num == 1 {
		return true
	}
	return false
}

func LimitIsOwn(db *sql.DB, uid, id string) bool {
	stmt, _ := db.Prepare("select count(*) from `limit` where uid = ? and id = ?")
	var num int
	err := stmt.QueryRow(uid, id).Scan(&num)
	if err != nil {
		return false
	}
	if num == 1 {
		return true
	}
	return false
}

func GetBetweenTime(hour, minute, second string) string {
	return time.Now().Format("2006-01-02") + " " + hour + ":" + minute + ":" + second
}

func IsAdmin(db *sql.DB, uid int) bool {
	stmt2, _ := db.Prepare("select `group` from user where uid = ? ")
	var group int
	err := stmt2.QueryRow(uid).Scan(&group)
	if err != nil {
		return false
	}
	if group == 6 {
		return true
	}
	return false
}
