package storage

type ImageStorage interface {
	// Store
	// Stores the image in the storage and returns the URL
	// if image id is empty, a random one will be generated
	Store(fileSuffix string, imageData []byte, imageId string) (string, error)

	Delete(path string) error
}
