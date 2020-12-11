package main

import (
	"log"

	utils "labs.zoe.im/gowin-virus-demo"
)

func main() {
	app, err := utils.NewApp("Example 3: Embed DLL")
	if err != nil {
		log.Fatalln(err)
		return
	}

	app.Run()
}