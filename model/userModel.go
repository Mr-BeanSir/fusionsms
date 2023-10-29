package model

import "github.com/golang-jwt/jwt/v4"

type User struct {
	Uid             int     `json:"uid,omitempty"`
	Username        string  `json:"username,omitempty"`
	Password        string  `json:"password,omitempty"`
	Group           int     `json:"group"`
	Quota           float32 `json:"quota"`
	Jwt             string  `json:"jwt,omitempty"`
	Status          int     `json:"status"`
	Code            string  `json:"code,omitempty"`
	Phone           string  `json:"phone,omitempty"`
	LimitQuotaPhone int     `json:"limit_quota_phone,omitempty"`
	WhiteIP         string  `json:"white_ip,omitempty"`
	BlackPhone      string  `json:"black_phone,omitempty"`
}

type UserClaims struct {
	Uid      int     `json:"uid"`
	UserName string  `json:"username"`
	PassWord string  `json:"password"`
	Status   int     `json:"status"`
	Group    int     `json:"group"`
	Quota    float32 `json:"balance"`
	jwt.RegisteredClaims
}
