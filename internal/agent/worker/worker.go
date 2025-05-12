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

		result, err := calculator.Calc(resp.Expression)
		if err != nil {
			log.Printf("Error calculating: %v", err)
			continue
		}

		_, sendErr := client.SendResult(context.Background(), &taskpb.TaskResponse{
			Id:     resp.Id,
			Result: result,
		})
		if sendErr != nil {
			log.Printf("Failed to send result: %v", sendErr)
		} else {
			log.Printf("Sent result for %s: %.2f", resp.Id, result)
		}
	}
}
