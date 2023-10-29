package model

import "time"

type SentLogModel struct {
	Id          int        `json:"id"`
	Uid         int        `json:"uid"`
	Sid         int        `json:"sid"`
	Content     string     `json:"content"`
	Phone       string     `json:"phone"`
	Time        time.Time  `json:"time"`
	TaskId      *int       `json:"task_id"`
	Receive     string     `json:"receive"`
	ReceiveTime *time.Time `json:"receive_time"`
	DecreaseNum int        `json:"decrease_num"`
	Status      int        `json:"status"`
}
