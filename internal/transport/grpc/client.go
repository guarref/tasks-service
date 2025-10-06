package grpc

import (
	userpb "github.com/guarref/project-protos/proto/user"
	"google.golang.org/grpc"
)

func NewUserClient(addr string) (userpb.UserServiceClient, *grpc.ClientConn, error) {
	// 1. grpc.Dial(addr, grpc.WithInsecure())
	// 2. userpb.NewUserServiceClient(conn)
	// 3. вернуть client, conn, err

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	defer conn.Close()

	client := userpb.NewUserServiceClient(conn)

	return client, conn, nil
}