package main

import (
	"log"
	"net"

	"github.com/H3199/doggodb/internal/api"
	"github.com/H3199/doggodb/internal/api/generated/db"
	"github.com/H3199/doggodb/internal/data"
	"github.com/H3199/doggodb/internal/query"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Setup storage and executor
	storage := data.NewInMemoryStorage()
	executor := query.NewExecutor(*storage)

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Create service implementation
	dbServer := &api.Server{
		Executor: executor,
	}

	// Register service with gRPC
	db.RegisterDatabaseServiceServer(grpcServer, dbServer)

	// Enable reflection
	reflection.Register(grpcServer)

	// Listen on a port
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Server listening on :50051")

	// Start serving
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
