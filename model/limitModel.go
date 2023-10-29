package model

type Limit struct {
	ID   int    `json:"id"`
	Uid  int    `json:"uid"`
	Sid  int    `json:"sid"`
	Num  int    `json:"num"`
	Time string `json:"time"`
}
