package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/matrixyuri/DouMusic/model"
)

func main() {
	if len(os.Args) < 4 {
		showHelp()
		return
	}
	mode := os.Args[1]

	imgName := getCaptcha()
	user := model.User{Name: "rubbish990@foxmail.com", Pass: "qwe123"}
	// user := model.User{Name: os.Args[2], Pass: os.Args[3]}
	var captcha string
	fmt.Print("请输入验证码: ")
	fmt.Scanf("%s\n", &captcha)
	user.Login(imgName, captcha)
	songs := model.GetList(user.Ck)
	if len(songs.Songs) == 0 {
		log.Fatal("\033[0;31m列表获取失败，请重试\033[0m")
	}
	fmt.Println("红心列表获取成功")
	fmt.Printf("总计: [%d]\n", len(songs.Songs))

	switch mode {
	case "-l":
		model.SaveList(user.NickName+"'s Love.txt", songs)
		fmt.Println("保存完毕")
	case "-d":
		songs.Download()
	}
	return
}

func showHelp() {
	fmt.Println("使用方法:")
	fmt.Println("\t-l [用户名] [密码]\t登录豆瓣帐号并自动保存红心歌曲列表")
	fmt.Println("\t-d [用户名] [密码]\t登录豆瓣帐号并自动下载红心歌曲")
}

func getCaptcha() string {
	imgName := model.GetCaptcha()
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", "./captcha/"+imgName+".jpeg")
	default:
		cmd = exec.Command("./captcha/" + imgName + ".jpeg")
	}
	cmd.Start()
	return imgName
}
