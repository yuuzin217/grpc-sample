package main

import (
	"fmt"
	"log"
	"net"
	"yuuzin217/grpc-sample/pb"

	"google.golang.org/grpc"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
)

type server struct {
	pb.UnimplementedFileServiceServer
}

func loadGRPCOpts() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(loadUnaryInterceptor()...)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(loadStreamInterceptor()...)),
	}
}

func main() {
	listen, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpc_opts := loadGRPCOpts()

	creds, err := loadCreds()
	if err != nil {
		log.Fatalln(err)
	}
	grpc_opts = append(grpc_opts, grpc.Creds(creds))

	s := grpc.NewServer(grpc_opts...)
	pb.RegisterFileServiceServer(s, &server{})
	fmt.Println("server is running...")
	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to Serve: %v", err)
	}
}
