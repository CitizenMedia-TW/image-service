package main

type ImageStorage interface {
	// Store
	// Stores the image in the storage and returns the URL
	Store(fileSuffix string, imageData []byte) (string, error)
}
