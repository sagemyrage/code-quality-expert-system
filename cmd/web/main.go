package main

import (
	"log"

	"github.com/sagemyrage/code-quality-expert-system/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load error: %v", err)
	}
	log.Printf("starting server on port %s", cfg.App.Port)
}
