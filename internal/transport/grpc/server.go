package grpc

import (
	"fmt"
	"net"

	taskpb "github.com/guarref/project-protos/proto/task"
	userpb "github.com/guarref/project-protos/proto/user"
	"github.com/guarref/tasks-service/internal/task"
	"google.golang.org/grpc"
)

func RunGRPC(svc *task.TaskService, uc userpb.UserServiceClient) error {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		return fmt.Errorf("gRPC fail to listen on port 50052: %v", err)
	}

	grpcSrv := grpc.NewServer()

	handler := NewHandler(svc, uc)
	taskpb.RegisterTaskServiceServer(grpcSrv, handler)

	if err := grpcSrv.Serve(lis); err != nil {
		return fmt.Errorf("Error of running server: %v", err)
	}
	return nil
}