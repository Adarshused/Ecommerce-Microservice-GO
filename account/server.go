package account

import (
	"context"
	"fmt"
	"net"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service Service
}

func ListenGRPC (s Service, port int) error {

   lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

   if err != nil {
	return err
   }

   serv := grpc.NewServer()


}