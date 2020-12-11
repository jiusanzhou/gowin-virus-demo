
package main

import (
	"log"
	"syscall"

	widgets "github.com/lxn/walk/declarative"
	"github.com/shirou/w32"
	utils "labs.zoe.im/gowin-virus-demo"
)

func main() {
	app, err := utils.NewApp(
		"Example 6: Query Reg",
		widgets.TextLabel{
			Text: "WeChat.exe =>"+GetWeChatExeLocation(),
		},
	)
	if err != nil {
		log.Fatalln(err)
		return
	}

	app.Run()
}

func GetWeChatExeLocation() string {
	defer func() {
		if recover() != nil {

		}
	}()
	key := w32.RegOpenKeyEx(w32.HKEY_CURRENT_USER, "Software\\Tencent\\WeChat", w32.KEY_ALL_ACCESS)

	ap := utils.RegQueryValueEx(syscall.Handle(key), "InstallPath")
	if ap == "" {
		return ""
	}
	return ap + "\\WeChat.exe"
}
