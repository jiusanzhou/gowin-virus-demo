age main

import (
	"log"

	widgets "github.com/lxn/walk/declarative"
	utils "labs.zoe.im/gowin-virus-demo"
)

func main() {
	app, err := utils.NewApp(
		"Example 2: Webview Form",
		widgets.WebView{
			Name:     "wv",
			URL:      "https://m.baidu.com",
		},
	)
	if err != nil {
		log.Fatalln(err)
		return
	}

	app.Run()
}
