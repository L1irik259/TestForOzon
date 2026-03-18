package server

import (
	"log"
	"net"
	"os"
	"path/filepath"

	pb "github.com/L1irik259/TestForOzon/internal/transport/proto/github.com/L1irik259/TestForOzon/transport/genetation/go/v1"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	adapter "github.com/L1irik259/TestForOzon/internal/adapter"
	service "github.com/L1irik259/TestForOzon/internal/service"
	transport "github.com/L1irik259/TestForOzon/internal/transport/service"
)

func StartServer() {
	exePath, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get working directory:", err)
	}

	envPath := filepath.Join(exePath, ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Println("Warning: .env file not found at", envPath, ", fallback to system env")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	itemAdapter := adapter.NewItemAdapter(db)
	itemService := service.NewItemService(itemAdapter)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	server := transport.NewServer(*itemService)
	pb.RegisterOzonServiceServer(grpcServer, server)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
