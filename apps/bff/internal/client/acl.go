package client

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/its-me-mayday/sisyphus/bff/proto/acl"
)

func NewAclClient(addr string) (pb.AclServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("acl client: %w", err)
	}
	return pb.NewAclServiceClient(conn), conn, nil
}
