package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"log"
	"net/http"
)

type UploadImageResponse struct {
	Message string `json:"message"`
	Url     string `json:"url"`
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

	// Get the file from the formdata
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

	url, err := a.storage.Store(handler.Filename, imageData)
	if err != nil {
		http.Error(w, url, http.StatusInternalServerError)
		return
	}

	// Respond with success message
	res := UploadImageResponse{
		Message: "Image uploaded successfully",
		Url:     url,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)

	return
}

// displayImage is the handler for the display route
// should only be used when the image storage is mongodb
func (a *App) displayImage(w http.ResponseWriter, r *http.Request) {
	mongoStorage := a.storage.(*MongoImageStorage)
	if mongoStorage == nil {
		http.Error(w, "Server implementation error, contact admin", http.StatusInternalServerError)
		return
	}
	// Check if the request is a GET request
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the query string parameter _id and collection
	queryValues := r.URL.Query()
	targetId := queryValues.Get("_id")
	collectionName := queryValues.Get("collection")

	// Convert the _id to an ObjectID
	docId, err := primitive.ObjectIDFromHex(targetId)
	if err != nil {
		log.Print(err)
		http.Error(w, "Error converting _id to ObjectID", http.StatusInternalServerError)
		return
	}
	log.Println(`bson.M{"_id": docID}:`, bson.M{"_id": docId})

	// Retrieve the image from MongoDB by _id
	var result MongoFields
	err = mongoStorage.database.Collection(collectionName).FindOne(context.Background(), bson.M{"_id": docId}).Decode(&result)
	if err != nil {
		log.Print(err)
		http.Error(w, "Error retrieving image from MongoDB", http.StatusInternalServerError)
		return
	}

	// Set appropriate headers
	w.Header().Set("Content-Type", result.Type)
	w.Header().Set("Content-Disposition", "inline; filename="+result.Name)

	// Write the image binary data directly to the response
	imageData := result.Data
	_, err = w.Write(imageData)
	if err != nil {
		http.Error(w, "Error writing image data to response", http.StatusInternalServerError)
	}
}

/**
// deleteImage is the handler for the delete route
func (a *App) deleteImage(w http.ResponseWriter, r *http.Request) {
	// Check if the request is a DELETE request
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the query string parameter _id and collection
	queryValues := r.URL.Query()
	targetId := queryValues.Get("_id")
	collectionName := queryValues.Get("collection")

	// Convert the _id to an ObjectID
	docId, err := primitive.ObjectIDFromHex(targetId)
	if err != nil {
		log.Print(err)
		http.Error(w, "Error converting _id to ObjectID", http.StatusInternalServerError)
		return
	}

	// Delete the image from MongoDB by _id
	result, err := a.database.Collection(collectionName).DeleteOne(context.Background(), bson.M{"_id": docId})
	if err != nil {
		log.Print(err)
		http.Error(w, "Error deleting image from MongoDB", http.StatusInternalServerError)
		return
	}
	log.Println("Deleted", result.DeletedCount, "documents")

	// Respond with success Message
	res := UploadImageResponse{
		Message: fmt.Sprintf("Deleted %d documents", result.DeletedCount),
		Id:      targetId,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

*/
