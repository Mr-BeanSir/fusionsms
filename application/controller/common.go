package controller

import (
	"errors"
	"fusionsms/config"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"net/http"
	"net/url"
)

var (
	BaseUrl    = "http://api.sms.douboke.cn"
	BaseData   = "token=" + config.Key + "&"
	BaseValues = url.Values{
		"token": {config.Key},
	}
)

func SetBaseValues(key string) {
	BaseValues = url.Values{
		"token": {key},
	}
}

func Smsdou(to, content string) error {
	BaseValues.Add("to", to)
	BaseValues.Add("content", content)
	resp, err := http.PostForm("/Api/Sent", BaseValues)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	re := gjson.Get(string(all), "bool")
	if !re.Bool() {
		return errors.New(gjson.Get(string(all), "msg").String())
	}
	return nil
}

func ApiAddSign(content string) (string, string, string, error) {
	BaseValues.Add("content", content)
	log.Println(BaseValues)
	resp, err := http.PostForm(BaseUrl+"/Api/addSign", BaseValues)
	if err != nil {
		return "", "", "", errors.New("接口错误，请重试")
	}
	defer resp.Body.Close()
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", err
	}
	re := gjson.Get(string(all), "bool")
	if !re.Bool() {
		return "", "", "", errors.New(gjson.Get(string(all), "msg").String())
	}
	return gjson.Get(string(all), "sign_id").String(),
		gjson.Get(string(all), "sign_key").String(),
		gjson.Get(string(all), "md5").String(), nil
}

func ApiResetSign(sign string) (string, string, error) {
	BaseValues.Add("sign_id", sign)
	resp, err := http.PostForm(BaseUrl+"/Api/resetKey", BaseValues)
	if err != nil {
		return "", "", errors.New("接口错误，请重试")
	}
	defer resp.Body.Close()
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	re := gjson.Get(string(all), "bool")
	if !re.Bool() {
		return "", "", errors.New(gjson.Get(string(all), "msg").String())
	}
	return gjson.Get(string(all), "key").String(),
		gjson.Get(string(all), "md5").String(), nil
}

func ApiAddTemplate(sign_id, content string) ([]gjson.Result, error) {
	BaseValues.Add("sign_id", sign_id)
	BaseValues.Add("content", content)
	resp, err := http.PostForm(BaseUrl+"/Api/addTemplates", BaseValues)
	if err != nil {
		return []gjson.Result{}, errors.New("接口错误，请重试")
	}
	defer resp.Body.Close()
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return []gjson.Result{}, err
	}
	log.Println(string(all))
	re := gjson.Get(string(all), "bool")
	if !re.Bool() {
		return []gjson.Result{}, errors.New(gjson.Get(string(all), "msg").String())
	}
	return gjson.Get(string(all), "template_ids").Array(), nil
}

func ApiDeleteTemplate(sid, tid string) error {
	BaseValues.Add("sign_id", sid)
	BaseValues.Add("template_id", tid)
	resp, err := http.PostForm(BaseUrl+"/Api/deleteTemplate", BaseValues)
	if err != nil {
		return errors.New("接口错误，请重试")
	}
	defer resp.Body.Close()
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	re := gjson.Get(string(all), "bool")
	if !re.Bool() {
		return errors.New(gjson.Get(string(all), "msg").String())
	}
	return nil
}
