package app

import (
	serverConfig "image-service/internal/config"
	"image-service/internal/database"
	"image-service/internal/grpcapp"
	"image-service/internal/restapp"
	"image-service/internal/storage"
	"log"
	"net/http"
)

const portNumber = ":80"

func StartServer() {
	log.Println("Starting image server.")
	config := serverConfig.NewConfig()
	var imageStorage storage.ImageStorage = storage.NewS3ImageStorage(config)
	var db database.Db = database.NewMongoDB(config)

	go grpcapp.StartGrpc(imageStorage, db)

	var app = restapp.New(
		imageStorage,
		db,
		config,
	)
	// Declare the server
	srv := &http.Server{
		Addr:    portNumber,
		Handler: app.Routes(),
	}

	log.Println("Image server started.")
	// Start the server
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
