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
	if req.Error != "" {
		db.UpdateExpressionError(req.Id, req.Error)
		log.Printf("Expression %s failed: %v", req.Id, req.Error)
		return &taskpb.Empty{}, nil
	}

	db.UpdateExpressionStatus(req.Id, "done", req.Result)
	log.Printf("Expression %s calculated: %.2f", req.Id, req.Result)

	return &taskpb.Empty{}, nil
}
