package main

import (
	"droneshield/config"
	"droneshield/internal/handlers"
	"droneshield/internal/middleware"
	"fmt"
	"log"
	"net/http"
)

func main() {
	db := config.InitDB()
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to Droneshield API")
	})

	mux.HandleFunc("/api/v1/geofences", handlers.GeoFenceHandler)

	wrappedMux := middleware.ErrorLoggingMiddleware(mux)

	port := "8000"
	fmt.Printf("Server started on port %s\n", port)
	if err := http.ListenAndServe(":"+port, wrappedMux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
