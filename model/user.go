package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type User struct {
	Ck         string `json:"ck"`   //不知道干嘛用
	Id         string `json:"id"`   //用户ID
	Name       string `json:"-"`    //用户名
	Pass       string `json:"-"`    //密码
	NickName   string `json:"name"` //昵称
	PlayRecord struct {
		Liked  int `json:"liked"`  //红心歌曲数目
		Played int `json:"played"` //播放过歌曲数目
		Banned int `json:"banned"` //不听歌曲数目
	} `json:"play_record"`
}

func GetCaptcha() string {
	captcha := httpDo(captchaApi, "", "GET")
	captcha = captcha[1 : len(captcha)-1]

	imgFile, _ := os.Create("./captcha/" + captcha + ".jpeg")
	defer imgFile.Close()

	url := fmt.Sprintf(captchaImgApi, captcha)
	img := httpDo(url, "", "GET")
	io.Copy(imgFile, bytes.NewReader([]byte(img)))

	return captcha
}

func (u *User) Login(captchaId string, captchaSolution string) {
	args := "source=radio&alias=%s&form_password=%s&captcha_id=%s&captcha_solution=%s&remember=on"
	args = fmt.Sprintf(args, u.Name, u.Pass, captchaId, captchaSolution)
	ret := httpDo(loginApi, args, "POST")

	type loginInfo struct {
		User   *User  `json:"user_info"`
		R      int    `json:"r"`
		ErrNo  int    `json:"err_no"`
		ErrMsg string `json:"err_msg"`
	}
	info := loginInfo{User: u}
	err := json.NewDecoder(strings.NewReader(ret)).Decode(&info)
	if err != nil {
		fmt.Println(err)
		return
	}
	if info.R != 0 {
		fmt.Println(info.ErrMsg)
		return
	}
}
