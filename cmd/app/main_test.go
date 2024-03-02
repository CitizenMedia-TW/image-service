package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"image-service/internal/app"
	"image-service/internal/restapp"
	"image-service/protobuffs/image-service"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func loadEnv() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file %s", err)
	}
}

func TestMain(m *testing.M) {
	loadEnv()
	go app.StartServer()
	time.Sleep(time.Second * 5) // Wait for the server to start
	os.Exit(m.Run())
}

func TestUploadAndDelete(t *testing.T) {
	time.Sleep(time.Second * 5)     // Wait for the server to start
	imagePath := "./test_image.png" // The test image file path

	// 1. Upload image via REST API
	resUpload := postImage(imagePath)
	assert.NotEmpty(t, resUpload.Url)
	assert.NotEmpty(t, resUpload.ImageId)
	assert.True(t, strings.Contains(resUpload.Url, resUpload.ImageId))
	r, err := http.Get(resUpload.Url) // Test the given URL
	assert.NoError(t, err)
	assert.Equal(t, 200, r.StatusCode)

	// Connect to the gRPC server
	conn, err := grpc.Dial("localhost:1111", grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)
	client := image_service.NewImageServiceClient(conn)

	// 2. Confirm usage
	resConfirm, err := client.ConfirmImageUse(context.TODO(), &image_service.ConfirmImageUsedRequest{
		ImageId: resUpload.ImageId,
		Strict:  false,
	})
	assert.NoError(t, err)
	assert.Equal(t, resConfirm.Success, true)

	// 3. Delete permanent image just uploaded
	resDelete, err := client.DeleteImage(context.TODO(), &image_service.DeleteImageRequest{ImageId: resUpload.ImageId})
	assert.NoError(t, err)
	assert.Equal(t, true, resDelete.Success)
	r, err = http.Get(resUpload.Url)
	assert.Equal(t, 403, r.StatusCode, "Should be deleted")
}

func postImage(imagePath string) restapp.UploadImageResponse {
	uResponse := restapp.UploadImageResponse{}

	// Create a new HTTP request with a POST method
	req, err := http.NewRequest("POST", "http://localhost:80/upload", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return uResponse
	}

	// Create a new buffer to store the request body
	var requestBody bytes.Buffer

	writer := multipart.NewWriter(&requestBody)

	// Open the image file
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return uResponse
	}
	defer file.Close()

	// Create a new form file field for the image
	fileField, err := writer.CreateFormFile("image", imagePath)
	if err != nil {
		fmt.Println("Error creating form file:", err)
		return uResponse
	}

	// Copy the image data to the form file field
	_, err = io.Copy(fileField, file)
	if err != nil {
		fmt.Println("Error copying file data:", err)
		return uResponse
	}

	err = errors.Join(
		writer.WriteField("uploader", "507f1f77bcf86cd799439011"),
		writer.WriteField("usage", "1"))

	if err != nil {
		fmt.Println("Error writing usage field:", err)
		return uResponse
	}

	// Close the multipart writer
	writer.Close()

	// Set the content type of the request to the multipart form data
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Set the request body to the buffer containing the multipart form data
	req.Body = io.NopCloser(&requestBody)

	// Perform the HTTP request
	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return uResponse
	}
	defer response.Body.Close()

	// Print the response status and body
	fmt.Println("Response Status:", response.Status)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return uResponse
	}

	json.Unmarshal(body, &uResponse)
	return uResponse
}
