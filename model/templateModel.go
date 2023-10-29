package model

type Template struct {
	Tid      int    `json:"tid"`
	Sid      int    `json:"sid"`
	Uid      int    `json:"uid"`
	Status   int    `json:"status"`
	Content  string `json:"content"`
	Reason   string `json:"reason"`
	SuperTid string `json:"super_tid"`
}
