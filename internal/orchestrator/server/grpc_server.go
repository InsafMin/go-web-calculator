package server

import (
	"context"
	"fmt"
	"log"

	"github.com/InsafMin/go-web-calculator/internal/db"
	"github.com/InsafMin/go-web-calculator/proto/taskpb"
)

type GRPCTaskServer struct {
	taskpb.UnimplementedTaskServiceServer
}

func (s *GRPCTaskServer) GetTask(ctx context.Context, req *taskpb.Empty) (*taskpb.TaskRequest, error) {
	expr := db.GetFirstPendingExpression()
	if expr == nil {
		return nil, fmt.Errorf("no pending expressions")
	}
	log.Printf("Sending task to worker: ID=%s, Expr=%s", expr.ID, expr.Expression)
	return &taskpb.TaskRequest{
		Id:         expr.ID,
		Expression: expr.Expression,
	}, nil
}

func (s *GRPCTaskServer) SendResult(ctx context.Context, req *taskpb.TaskResponse) (*taskpb.Empty, error) {
	log.Printf("Received result for %s: %.2f", req.Id, req.Result)
	err := db.UpdateExpressionStatus(req.Id, "done", req.Result)
	if err != nil {
		log.Printf("Failed to update expression status in DB: %v", err)
	}
	return &taskpb.Empty{}, nil
}
