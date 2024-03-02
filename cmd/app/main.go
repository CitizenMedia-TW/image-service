package main

import (
	"github.com/joho/godotenv"
	"image-service/internal/app"
	"log"
)

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file %s", err)
	}

	// Start the server
	app.StartServer()
}
