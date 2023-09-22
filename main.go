package main

import (
	"context"
	"github.com/auyeongweibin/auth-server/auth"
	"github.com/auyeongweibin/auth-server/configs"
	"google.golang.org/grpc"
	"log"
	"net"
)

type myRegistrationServer struct {
	auth.UnimplementedRegistrationServer
}

func (s myRegistrationServer) Create(context.Context, *auth.RegistrationRequest) (*auth.RegistrationResponse, error) {
	return &auth.RegistrationResponse{
		Email: "sample@email.com", 
		Message: "registration successful!",
	}, nil
}

func main() {
	configs.ConnectDB()

	lis, err := net.Listen("tcp", ":8089")
	if err != nil {
		log.Fatalf("cannot create listener: %s", err)
	}

	serverRegistrar := grpc.NewServer()
	service := &myRegistrationServer{}
	auth.RegisterRegistrationServer(serverRegistrar, service)

	err = serverRegistrar.Serve(lis)
	if err != nil {
		log.Fatalf("impossible to serve: %s", err)
	}
}