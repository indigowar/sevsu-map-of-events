package main

import (
	"github.com/indigowar/map-of-events/internal/app"
	"github.com/indigowar/map-of-events/internal/config"
	"log"
)

func main() {
	cfg, err := config.Init("config")
	if err != nil {
		log.Fatalln("Wrong configuration of application ", err)
	}
	app.Run(cfg)
}
