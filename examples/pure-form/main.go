package main

import (
	"log"

	utils "labs.zoe.im/gowin-virus-demo"
)

func main() {
	app, err := utils.NewApp("Example 1: Pure Form")
	if err != nil {
		log.Fatalln(err)
		return
	}

	app.Run()
}