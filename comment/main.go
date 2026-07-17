package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	commentpb "comment/gen"
	"comment/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()

	commentpb.RegisterCommentServiceServer(grpcServer, server.NewCommentServer())

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	reflection.Register(grpcServer)

	log.Printf("comment service listening on :%s", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("grpc server failed: %v", err)
	}
}
