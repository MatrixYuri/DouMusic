package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/matrixyuri/DouMusic/model"
)

func main() {
	if len(os.Args) <= 1 {
		showHelp()
		return
	}
	mode := os.Args[1]
	switch mode {
	case "-list":
		if len(os.Args) < 4 {
			showHelp()
			return
		}
		imgName := getCaptcha()
		user := model.User{Name: "rubbish990@foxmail.com", Pass: "qwe123"}
		// user := model.User{Name: os.Args[2], Pass: os.Args[3]}
		var captcha string
		fmt.Print("请输入验证码: ")
		fmt.Scanf("%s\n", &captcha)
		user.Login(imgName, captcha)
		songs := model.GetList(user.Ck)
		fmt.Printf("红心歌曲总计: [%d]\n", len(songs.Songs))
		model.SaveList(user.NickName, songs)
	case "-download":

	}
	// gl.StartDriver(view.MainWindow)
}

func showHelp() {
	fmt.Println("使用方法:")
	fmt.Println("\t-list [用户名] [密码]\t登录豆瓣帐号并自动保存红心歌曲列表")
	fmt.Println("\t-download [用户名] [密码]\t登录豆瓣帐号并自动下载红心歌曲")
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
