package main

import (
	"log"
	"qrstreamer/internal/app"
	"qrstreamer/util"
)

func main() {
	cfg, err := util.LoadConfig("./")
	if err != nil {
		log.Fatal(err)
	}
	app.Run(cfg)
}
