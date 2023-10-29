package model

import (
	"database/sql"
	"time"
)

type ExchangeLog struct {
	Id      int        `json:"id,omitempty"`
	Uid     int        `json:"uid,omitempty"`
	Content string     `json:"content,omitempty"`
	Time    *time.Time `json:"time,omitempty"`
}

type ExchangeCode struct {
	Id         int           `json:"id,omitempty"`
	Code       string        `json:"code,omitempty"`
	Quota      int           `json:"quota,omitempty"`
	Status     int           `json:"status"`
	CreateTime *time.Time    `json:"create_time"`
	UseTime    *time.Time    `json:"use_time"`
	UseUid     sql.NullInt16 `json:"use_uid"`
}
