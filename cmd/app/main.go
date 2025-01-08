package main

import (
	"log"

	"yalp_ulab/config"
	"yalp_ulab/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
