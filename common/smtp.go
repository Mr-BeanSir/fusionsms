package common

import (
	"crypto/tls"
	"fmt"
	"fusionsms/model"
	"github.com/gin-gonic/gin"
	"github.com/jordan-wright/email"
	"net/smtp"
	"strings"
)

func SendWithContext(c *gin.Context, title, content, to string) error {
	db := GetDbInfo(c)
	var system model.System
	err := db.QueryRow("select smtp_server,smtp_password,smtp_username,smtp_nickname,smtp_ssl from system").
		Scan(&system.SmtpServer, &system.SmtpPassword, &system.SmtpUsername, &system.SmtpNickname, &system.SmtpSSL)
	if err != nil {
		return err
	}
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", system.SmtpNickname, system.SmtpUsername)
	e.To = []string{to}
	e.Subject = title
	e.Text = []byte(content)
	if system.SmtpSSL == 1 {
		err = e.SendWithTLS(system.SmtpServer, smtp.PlainAuth("", system.SmtpUsername, system.SmtpPassword, system.SmtpServer[:strings.Index(system.SmtpServer, ":")]), &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         system.SmtpServer[:strings.Index(system.SmtpServer, ":")],
		})
	} else {
		err = e.Send(system.SmtpServer, smtp.PlainAuth("", system.SmtpUsername, system.SmtpPassword, system.SmtpServer[:strings.Index(system.SmtpServer, ":")]))
	}
	if err != nil {
		return err
	}
	return nil
}

func SendEmail(SmtpServer, SmtpUsername, SmtpPassword, SmtpNickname, title, content, to, ssl string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", SmtpNickname, SmtpUsername)
	e.To = []string{to}
	e.Subject = title
	e.Text = []byte(content)
	var err error
	if ssl == "1" {
		err = e.SendWithTLS(SmtpServer, smtp.PlainAuth("", SmtpUsername, SmtpPassword, SmtpServer[:strings.Index(SmtpServer, ":")]), &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         SmtpServer[:strings.Index(SmtpServer, ":")],
		})
	} else {
		err = e.Send(SmtpServer, smtp.PlainAuth("", SmtpUsername, SmtpPassword, SmtpServer[:strings.Index(SmtpServer, ":")]))
	}
	if err != nil {
		return err
	}
	return nil
}
