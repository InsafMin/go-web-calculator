package worker

import (
	"context"
	"log"
	"time"

	"github.com/InsafMin/go-web-calculator/pkg/calculator"
	"github.com/InsafMin/go-web-calculator/proto/taskpb"
	"google.golang.org/grpc"
)

var client taskpb.TaskServiceClient

func ConnectToOrchestrator(grpcURL string) {
	conn, err := grpc.Dial(grpcURL, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Could not connect to gRPC server: %v\n", err)
	}
	client = taskpb.NewTaskServiceClient(conn)
	log.Println("Connected to orchestrator via gRPC")
}

func StartWorker() {
	for {
		resp, err := client.GetTask(context.Background(), &taskpb.Empty{})
		if err != nil {
			log.Printf("No tasks available, retrying... (%v)", err)
			time.Sleep(1 * time.Second)
			continue
		}

		log.Printf("Got task: ID=%s, Expr=%s", resp.Id, resp.Expression)

		result, calcErr := calculator.Calc(resp.Expression)
		errorMsg := ""
		if calcErr != nil {
			errorMsg = calcErr.Error()
			log.Printf("Error calculating %s: %v", resp.Id, calcErr)
		}

		_, sendErr := client.SendResult(context.Background(), &taskpb.TaskResponse{
			Id:     resp.Id,
			Result: result,
			Error:  errorMsg,
		})
		if sendErr != nil {
			log.Printf("Failed to send result for %s: %v\n", resp.Id, sendErr)
		}
	}
}
