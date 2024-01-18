package api

import (
	"image-service/cmd/imageservice/cnf"
	"image-service/cmd/imageservice/database"
	"image-service/cmd/imageservice/img_storage"
)

type App struct {
	storage  img_storage.ImageStorage
	database database.Db
	Config   cnf.ImageServiceConfig
}
