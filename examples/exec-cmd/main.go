package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"

	utils "labs.zoe.im/gowin-virus-demo"
)

func main() {
	app, err := utils.NewApp(
		"Example 6: Exec Command",
	)
	if err != nil {
		log.Fatalln(err)
		return
	}

	OpenURL("https://www.baidu.com")

	app.Run()
}

// OpenURL ...
func OpenURL(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	return err
}