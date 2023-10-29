package model

type Sign struct {
	Sid     int    `json:"id"`
	Uid     int    `json:"uid"`
	Content string `json:"sign"`
	Key     string `json:"key"`
	Status  int    `json:"status"`
	SuperId int    `json:"superId"`
}
