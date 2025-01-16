package main

import (
	"droneshield/config"
	"droneshield/internal/handlers"
	"droneshield/internal/middleware"
	"fmt"
	"net/http"

	ddhttp "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"go.uber.org/zap"
)

type contextKey string

const dbContextKey contextKey = "db"

func main() {
	tracer.Start(
		tracer.WithServiceName("droneshield-api"),
		tracer.WithEnv("production"),           // Replace with "staging" or "development" as needed
		tracer.WithAgentAddr("127.0.0.1:8126"), // Default Datadog Agent address
	)
	defer tracer.Stop()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	db, err := config.InitDB()
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to Droneshield API")
	})

	mux.HandleFunc("/api/v1/geofences", handlers.GeoFenceHandler)

	tracedMux := ddhttp.WrapHandler(mux, "droneshield-api", "production")

	wrappedMux := middleware.ErrorLoggingMiddleware(tracedMux)

	port := "8000"
	fmt.Printf("Server started on port %s\n", port)
	if err := http.ListenAndServe(":"+port, wrappedMux); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
