package api

import (
	"github.com/joho/godotenv"

	config2 "image-service/cmd/imageservice/cnf"
	"image-service/cmd/imageservice/database"
	"image-service/cmd/imageservice/img_storage"
	"log"
	"net/http"
)

// portNumber is the port number that the server will listen on
const portNumber = ":80"

func StartServer() {
	log.Println("Starting image server.")
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file %s", err)
	}
	config := config2.NewConfig()
	var imageStorage img_storage.ImageStorage = img_storage.NewS3ImageStorage(config)
	var db database.Db = database.NewMongoDB(config)

	var app = App{
		storage:  imageStorage,
		database: db,
		Config:   config,
	}
	// Declare the server
	srv := &http.Server{
		Addr:    portNumber,
		Handler: app.routes(),
	}

	log.Println("Image server started.")
	// Start the server
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
