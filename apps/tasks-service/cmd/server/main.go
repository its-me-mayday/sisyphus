package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/its-me-mayday/sisyphus/tasks-service/internal/handler"
	"github.com/its-me-mayday/sisyphus/tasks-service/internal/repository"

	pb "github.com/its-me-mayday/sisyphus/tasks-service/proto/tasks"
)

func main() {
	dsn := getEnv("DATABASE_URL", "postgres://sisyphus:sisyphus@localhost:5432/sisyphus?sslmode=disable")
	port := getEnv("GRPC_PORT", "50051")

	// Repository
	repo, err := repository.New(dsn)
	if err != nil {
		log.Fatalf("repository: %v", err)
	}

	// Migrate
	if err := repo.Migrate(context.Background()); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	log.Println("✓ migrations ok")

	// gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterTasksServiceServer(srv, handler.New(repo))
	reflection.Register(srv) // utile per grpcurl in dev

	log.Printf("✓ tasks-service listening on :%s", port)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
