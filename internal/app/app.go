package app

// import (
// 	"image-service/cmd/imageservice/cnf"
// 	"image-service/cmd/imageservice/database"
// 	"image-service/cmd/imageservice/img_storage"
// )

import (
	"github.com/joho/godotenv"

	config2 "image-service/internal/config"
	"image-service/internal/database"
	"image-service/internal/grpcapp"
	"image-service/internal/restapp"
	"image-service/internal/storage"
	"log"
	"net/http"
)

// type App struct {
// 	storage  img_storage.ImageStorage
// 	database database.Db
// 	Config   cnf.ImageServiceConfig
// }

const portNumber = ":80"

func StartServer() {
	log.Println("Starting image server.")
	// Load the .env file
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file %s", err)
	}
	config := config2.NewConfig()
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
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
