package main

import (
	"github.com/indigowar/map-of-events/internal/app"
	"github.com/indigowar/map-of-events/internal/config"
)

func main() {
	cfg := config.LoadConfig()
	app.Run(cfg)
}
