package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	httpHandler "github.com/yourusername/hospital-system-api/document-service/internal/delivery/http"
	"github.com/yourusername/hospital-system-api/document-service/internal/domain"
	"github.com/yourusername/hospital-system-api/document-service/internal/repository"
	"github.com/yourusername/hospital-system-api/document-service/internal/service"
	"github.com/yourusername/hospital-system-api/document-service/pkg/auth"
	"github.com/yourusername/hospital-system-api/document-service/pkg/elasticsearch"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	dsn := "host=" + os.Getenv("DB_HOST") +
		" port=" + os.Getenv("DB_PORT") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate the schema
	if err := db.AutoMigrate(&domain.Document{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize repository
	documentRepo := repository.NewDocumentRepository(db)

	// Initialize auth client
	authClient := auth.NewClient()

	// Initialize Elasticsearch client
	esClient, err := elasticsearch.NewClient()
	if err != nil {
		log.Fatalf("Failed to create Elasticsearch client: %v", err)
	}

	// Initialize service
	documentService := service.NewDocumentService(documentRepo, esClient, authClient)

	// Initialize HTTP handler
	handler := httpHandler.NewHandler(documentService, authClient)

	// Initialize router
	router := gin.Default()
	handler.RegisterRoutes(router)

	// Create server
	srv := &http.Server{
		Addr:    ":8004",
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
