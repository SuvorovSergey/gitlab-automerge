package main

import (
	"log"

	"github.com/SuvorovSergey/gitlab-automerge/internal/app"
	"github.com/SuvorovSergey/gitlab-automerge/internal/config"
)

func main() {
	cfg, err := config.NewConfig("config.yml")

	if err != nil {
		log.Fatalf("Error occured while loading configuration: %v", err)
	}

	app.Run(cfg)
}
