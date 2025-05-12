package main

import (
	"log"
	"os"

	"github.com/InsafMin/go-web-calculator/internal/agent/worker"
)

func main() {
	grpcURL := os.Getenv("ORCHESTRATOR_GRPC_URL")
	if grpcURL == "" {
		grpcURL = "localhost:50051"
	}

	log.Println("Starting worker...")
	worker.ConnectToOrchestrator(grpcURL)
	worker.StartWorker()
}
