package restapp

import (
	"log"
	"net/http"
)

import (
	"image-service/internal/config"
	"image-service/internal/database"
	"image-service/internal/storage"
)

type RestImageService struct {
	storage  storage.ImageStorage
	database database.Db
	Config   config.ImageServiceConfig
}

func New(storage storage.ImageStorage, database database.Db, config config.ImageServiceConfig) RestImageService {
	return RestImageService{
		storage:  storage,
		database: database,
		Config:   config,
	}
}

// routes defines the routes for the application
func (a *RestImageService) Routes() http.Handler {
	// Declare a new router
	mux := http.NewServeMux()

	// Use middleware for all routes
	mux.Handle("/", a.middlewares(http.HandlerFunc(a.home)))
	mux.Handle("/upload", a.middlewares(http.HandlerFunc(a.uploadImage)))
	//mux.Handle("/display", a.middlewares(http.HandlerFunc(a.displayImage)))
	//mux.Handle("/delete", a.middlewares(http.HandlerFunc(a.deleteImage)))

	// mux.HandleFunc("/", a.home)
	// mux.HandleFunc("/upload", a.uploadImage)
	// mux.HandleFunc("/display", a.displayImage)
	// mux.HandleFunc("/delete", a.deleteImage)

	return mux
}

// Define the middlewares for the application
func (a *RestImageService) middlewares(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("Executing middleware")

		//Allow CORS here By * or specific origin
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
		// log.Print("Executing middlewareTwo again")
	})
}

// home is the handler for the home page
func (a *RestImageService) home(w http.ResponseWriter, r *http.Request) {
	log.Println("Home page accessed")
	w.Write([]byte("Hello from a HandleFunc #1"))
}
