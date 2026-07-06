package main

import (
	"log"
	"os"

	"github.com/nhassl3/hairdress_arz/internal/app"
	"github.com/nhassl3/hairdress_arz/internal/config"
)

func main() {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		env := os.Getenv("ENVIRONMENT")
		switch env {
		case "prod":
			configFile = "config/prod.yaml"
		default:
			configFile = "config/local.yaml"
		}
	}

	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}

	cfg, err := config.LoadConfig(configFile, envFile)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	if err := app.Run(cfg); err != nil {
		log.Fatalf("application error: %v", err)
	}
}
