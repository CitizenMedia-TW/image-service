package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

// portNumber is the port number that the server will listen on
const portNumber = ":80"

type App struct {
	storage ImageStorage
}

// main is the entry point for the application
func main() {

	log.Println("Starting server on port", portNumber)

	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file %s", err)
	}
	config := NewConfig()
	var imageStorage ImageStorage
	switch config.storageType {
	case "mongo":
		imageStorage = NewMongoImageStorage(config)
		break
	case "s3":
		imageStorage = NewS3ImageStorage(config)
		break
	default:
		log.Fatalf("Cannot find image storage type %s", config.storageType)
		return
	}

	var app = App{
		imageStorage,
	}
	// Declare the server
	srv := &http.Server{
		Addr:    portNumber,
		Handler: app.routes(),
	}

	// Start the server
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
