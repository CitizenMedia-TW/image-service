package main

import (
	"bytes"
	"errors"
	"fmt"
	"image-service/cmd/imageservice/api"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestUpload(t *testing.T) {
	go api.StartServer()
	time.Sleep(time.Second * 5)
	// Specify the file path of the image you want to upload
	imagePath := "../test_image.png"

	// Create a new HTTP request with a POST method
	req, err := http.NewRequest("POST", "http://localhost:80/upload", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Create a new buffer to store the request body
	var requestBody bytes.Buffer

	writer := multipart.NewWriter(&requestBody)

	// Open the image file
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a new form file field for the image
	fileField, err := writer.CreateFormFile("image", imagePath)
	if err != nil {
		fmt.Println("Error creating form file:", err)
		return
	}

	// Copy the image data to the form file field
	_, err = io.Copy(fileField, file)
	if err != nil {
		fmt.Println("Error copying file data:", err)
		return
	}

	err = errors.Join(
		writer.WriteField("uploader", "507f1f77bcf86cd799439011"),
		writer.WriteField("usage", "1"))

	if err != nil {
		fmt.Println("Error writing usage field:", err)
		return
	}

	// Close the multipart writer
	writer.Close()

	// Set the content type of the request to the multipart form data
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Set the request body to the buffer containing the multipart form data
	req.Body = ioutil.NopCloser(&requestBody)

	// Perform the HTTP request
	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer response.Body.Close()

	// Print the response status and body
	fmt.Println("Response Status:", response.Status)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	fmt.Println("Response Body:", string(body))
}
