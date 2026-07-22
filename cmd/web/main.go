package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"golang-clean-architecture/internal/config"
	"os"
	"strconv"
)

func main() {
	godotenv.Load() // ponytail: silent fail if .env missing
	log := config.NewLogger()
	db := config.NewDatabase(log)
	validate := config.NewValidator()
	app := config.NewFiber()
	producer := config.NewKafkaProducer(log)

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
		Producer: producer,
	})

	webPort, _ := strconv.Atoi(os.Getenv("WEB_PORT"))
	err := app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
