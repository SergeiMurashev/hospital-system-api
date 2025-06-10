package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/sergeimurashev/hospital-system-api/hospital-service/internal/delivery/grpc"
	"github.com/sergeimurashev/hospital-system-api/hospital-service/internal/domain"
	"github.com/sergeimurashev/hospital-system-api/hospital-service/internal/repository"
	"github.com/sergeimurashev/hospital-system-api/hospital-service/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	hospitalRepo := repository.NewHospitalRepository(db)
	roomRepo := repository.NewRoomRepository(db)

	hospitalService := service.NewHospitalService(hospitalRepo, roomRepo)

	grpcServer := grpc.NewServer()
	hospitalServer := grpc.NewServer(hospitalService)
	proto.RegisterHospitalServiceServer(grpcServer, hospitalServer)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("GRPC_PORT")))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Starting gRPC server on port %s", os.Getenv("GRPC_PORT"))
	if err := grpcServer.Server(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func initDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&domain.Hospital{}, &domain.Room{}); err != nil {
		return nil, err
	}

	return db, nil
}
