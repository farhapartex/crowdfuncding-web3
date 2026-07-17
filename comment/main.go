package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	commentpb "comment/gen"
	"comment/db"
	"comment/server"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on existing environment variables")
	}

	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDB := os.Getenv("POSTGRES_DB")
	postgresPort := os.Getenv("POSTGRES_PORT")
	if postgresUser == "" || postgresPassword == "" || postgresDB == "" || postgresPort == "" {
		log.Fatal("POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB, and POSTGRES_PORT must be set")
	}

	postgresHost := os.Getenv("POSTGRES_HOST")
	if postgresHost == "" {
		postgresHost = "localhost"
	}

	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		postgresUser, postgresPassword, postgresHost, postgresPort, postgresDB,
	)

	gormDB, err := db.Connect(databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.Migrate(gormDB); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()

	commentpb.RegisterCommentServiceServer(grpcServer, server.NewCommentServer(gormDB))

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	reflection.Register(grpcServer)

	log.Printf("comment service listening on :%s", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("grpc server failed: %v", err)
	}
}
