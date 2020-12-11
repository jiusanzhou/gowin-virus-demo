package main

import (
	"log"
	"syscall"

	utils "labs.zoe.im/gowin-virus-demo"
)

func main() {
	app, err := utils.NewApp("Example 5: Attach Window")
	if err != nil {
		log.Fatalln(err)
		return
	}

	mw := utils.NewWindow(syscall.Handle(app.MainWindow().Handle()))

	WindowAttachToWeChat(mw, "微信测试版")

	app.Run()
}