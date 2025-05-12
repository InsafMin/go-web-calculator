package server

import (
	"log"
	"net"
	"net/http"

	"github.com/InsafMin/go-web-calculator/internal/auth"
	"github.com/InsafMin/go-web-calculator/internal/orchestrator/handlers"
	"github.com/InsafMin/go-web-calculator/proto/taskpb"
	"google.golang.org/grpc"
)

func StartHTTPServer() {
	http.HandleFunc("/api/v1/register", methodCheck(http.MethodPost, handlers.HandleRegister))
	http.HandleFunc("/api/v1/login", methodCheck(http.MethodPost, handlers.HandleLogin))
	http.HandleFunc("/api/v1/calculate", auth.AuthMiddleware(methodCheck(http.MethodPost, handlers.HandleCalculate)))
	http.HandleFunc("/api/v1/expressions", auth.AuthMiddleware(methodCheck(http.MethodGet, handlers.HandleGetExpressions)))
	http.HandleFunc("/api/v1/expressions/", auth.AuthMiddleware(methodCheck(http.MethodGet, handlers.HandleGetExpression)))

	log.Println("Starting HTTP server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func StartGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	taskpb.RegisterTaskServiceServer(s, &GRPCTaskServer{})

	log.Println("gRPC Server is running on :50051")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

func methodCheck(allowedMethod string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != allowedMethod {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		next(w, r)
	}
}
