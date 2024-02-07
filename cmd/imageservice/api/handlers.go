package api

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"image-service/cmd/imageservice/database/entities"
	image_service2 "image-service/protobuffs/image-service"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type UploadImageResponse struct {
	Message string `json:"message"`
	Url     string `json:"url"`
	ImageId string `json:"imageId"`
}

// uploadImage is the handler for the upload route
func (a *App) uploadImage(w http.ResponseWriter, r *http.Request) {
	// Check if the request is a POST request
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get the collection name from the query string
	//collectionName := r.URL.Query().Get("collection")

	// Parse our multipart form, 10 << 20 specifies a maximum upload of 10 MB files.
	parseErr := r.ParseMultipartForm(10 << 20)

	if parseErr != nil {
		http.Error(w, "Could not parse multipart form", http.StatusBadRequest)
		return
	}

	usageId, _ := strconv.Atoi(r.FormValue("usage"))

	usage := image_service2.ParseUsage(usageId)
	if usage == image_service2.ImageUsage_Undefined {
		http.Error(w, "Invalid image usage.", http.StatusBadRequest)
		return
	}

	// Get the file from the form data
	file, handler, err := r.FormFile("image")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
	}
	defer file.Close()

	// Read the image data & type
	imageData, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error Reading the File")
		fmt.Println(err)
	}

	if !strings.Contains(http.DetectContentType(imageData), "image") {
		http.Error(w, "該檔案不是圖片", http.StatusBadRequest)
	}

	imgId := primitive.NewObjectID()

	url, err := a.storage.Store(extractFileNameSuffix(handler.Filename), imageData, imgId.Hex())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpImage := entities.NewTmpImageWithId(imgId, usage, url)

	_, err = a.database.StoreTmpImgInfo(context.TODO(), tmpImage)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success message
	res := UploadImageResponse{
		Message: "Image uploaded successfully",
		Url:     a.Config.Host + url,
		ImageId: imgId.Hex(),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)

	return
}

func extractFileNameSuffix(fileName string) string {
	split := strings.Split(fileName, ".")
	return split[len(split)-1]
}
