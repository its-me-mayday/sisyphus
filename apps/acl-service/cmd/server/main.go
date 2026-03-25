package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/its-me-mayday/sisyphus/acl-service/internal/handler"
	"github.com/its-me-mayday/sisyphus/acl-service/internal/repository"

	pb "github.com/its-me-mayday/sisyphus/acl-service/proto/acl"
)

func main() {
	dsn := getEnv("DATABASE_URL", "postgres://sisyphus:sisyphus@localhost:5432/sisyphus?sslmode=disable")
	port := getEnv("GRPC_PORT", "50052")

	repo, err := repository.New(dsn)
	if err != nil {
		log.Fatalf("repository: %v", err)
	}

	if err := repo.Migrate(context.Background()); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	log.Println("✓ migrations ok")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterAclServiceServer(srv, handler.New(repo))
	reflection.Register(srv)

	log.Printf("✓ acl-service listening on :%s", port)
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
