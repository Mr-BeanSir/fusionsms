package model

type System struct {
	EmailCheck   int    `json:"email_check,omitempty"`
	SmtpServer   string `json:"smtp_server,omitempty"`
	SmtpUsername string `json:"smtp_username,omitempty"`
	SmtpPassword string `json:"smtp_password,omitempty"`
	SmtpNickname string `json:"smtp_nickname,omitempty"`
	SmtpSSL      int    `json:"smtp_ssl"`
}
