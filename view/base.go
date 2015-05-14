package view

import (
	"io/ioutil"
	"log"

	"github.com/google/gxui"
	"github.com/google/gxui/math"
	"github.com/google/gxui/themes/dark"
)

var (
	windows gxui.Window
	theme   gxui.Theme
)

func MainWindow(driver gxui.Driver) {
	theme = dark.CreateTheme(driver)

	fontData, err := ioutil.ReadFile("./fonts/Microsoft Yahei.ttf")
	if err != nil {
		log.Fatalf("error reading font: %v", err)
	}
	font, err := driver.CreateFont(fontData, 20)
	if err != nil {
		panic(err)
	}
	theme.SetDefaultFont(font)

	headLayout := theme.CreateLinearLayout()
	headLayout.SetSizeMode(gxui.Fill)
	headLayout.SetHorizontalAlignment(gxui.AlignCenter)
	headLayout.SetDirection(gxui.TopToBottom)
	headLayout.SetSize(math.Size{W: 100, H: 50})

	label := theme.CreateLabel()
	label.SetMargin(math.CreateSpacing(10))
	label.SetText("豆瓣FM红心歌单下载")

	headLayout.AddChild(label)

	bodyLayout := theme.CreateLinearLayout()
	bodyLayout.SetSizeMode(gxui.Fill)
	bodyLayout.SetHorizontalAlignment(gxui.AlignCenter)
	bodyLayout.SetSize(math.Size{W: 100, H: 400})
	bodyLayout.SetDirection(gxui.TopToBottom)
	// bodyLayout.SetMargin(math.Spacing{T: 50})

	userHint := theme.CreateLabel()
	userHint.SetMargin(math.CreateSpacing(10))
	userHint.SetText("用户名")
	bodyLayout.AddChild(userHint)
	usernameField := theme.CreateTextBox()
	bodyLayout.AddChild(usernameField)

	passHint := theme.CreateLabel()
	passHint.SetMargin(math.CreateSpacing(10))
	passHint.SetText("密码")
	bodyLayout.AddChild(passHint)
	passwordField := theme.CreateTextBox()
	bodyLayout.AddChild(passwordField)

	footLayout := theme.CreateLinearLayout()
	footLayout.SetSizeMode(gxui.Fill)
	footLayout.SetHorizontalAlignment(gxui.AlignCenter)
	footLayout.SetDirection(gxui.TopToBottom)
	// footLayout.SetMargin(math.Spacing{T: 450})

	button := theme.CreateButton()
	button.SetText("登录并下载")
	button.OnClick(func(gxui.MouseEvent) {
		user := usernameField.Text()
		password := passwordField.Text()
		if user == "" || password == "" {
			return
		}
		label.SetText(user + "&" + password)
	})
	footLayout.AddChild(button)

	window := theme.CreateWindow(600, 800, "豆瓣FM下载器")
	window.SetBackgroundBrush(gxui.CreateBrush(gxui.White))

	window.AddChild(headLayout)
	window.AddChild(bodyLayout)
	window.AddChild(footLayout)

	window.OnClose(driver.Terminate)
	window.SetPadding(math.Spacing{L: 10, T: 10, R: 10, B: 10})
}

func CreateLayout() {

}

func CreateButton(text string) gxui.Button {
	button := theme.CreateButton()
	button.SetText(text)
	return button
}

func CreateLabel(text string) gxui.Label {
	label := theme.CreateLabel()
	label.SetText(text)
	return label
}

func CreateInput() gxui.TextBox {
	return theme.CreateTextBox()
}
