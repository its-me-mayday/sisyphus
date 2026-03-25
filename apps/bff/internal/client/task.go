package client

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/its-me-mayday/sisyphus/bff/proto/tasks"
)

func NewTasksClient(addr string) (pb.TasksServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("tasks client: %w", err)
	}
	return pb.NewTasksServiceClient(conn), conn, nil
}
