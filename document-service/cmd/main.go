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
	httpHandler "github.com/sergeimurashev/hospital-system-api/document-service/internal/delivery/http"
	"github.com/sergeimurashev/hospital-system-api/document-service/internal/domain"
	"github.com/sergeimurashev/hospital-system-api/document-service/internal/repository"
	"github.com/sergeimurashev/hospital-system-api/document-service/internal/service"
	"github.com/sergeimurashev/hospital-system-api/document-service/pkg/auth"
	"github.com/sergeimurashev/hospital-system-api/document-service/pkg/elasticsearch"
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

	if err := db.AutoMigrate(&domain.Document{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	documentRepo := repository.NewDocumentRepository(db)

	authClient := auth.NewClient()

	esClient, err := elasticsearch.NewClient()
	if err != nil {
		log.Fatalf("Failed to create Elasticsearch client: %v", err)
	}

	documentService := service.NewDocumentService(documentRepo, esClient, authClient)

	handler := httpHandler.NewHandler(documentService, authClient)

	router := gin.Default()
	handler.RegisterRoutes(router)

	srv := &http.Server{
		Addr:    ":8004",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
