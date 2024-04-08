package main

import (
	"context"
	"fmt"
	"log"
	"net"
	grpc_sample "yuuzin217/grpc-sample"
	"yuuzin217/grpc-sample/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
)

type server struct {
	pb.UnimplementedFileServiceServer
}

func myLogging() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		log.Printf("request data: %+v", req)

		resp, err = handler(ctx, req)
		if err != nil {
			return nil, err
		}

		log.Printf("response data: %+v", resp)

		return resp, nil
	}
}

func authorize(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}
	if token != "test-token" {
		// return nil, errors.New("bad token")
		return nil, status.Error(codes.Unauthenticated, "token is invalid.")
	}
	return ctx, nil
}

func main() {
	listen, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	creds, err := credentials.NewServerTLSFromFile(
		grpc_sample.CertFile_filePath,
		grpc_sample.CertKey_filePath,
	)
	if err != nil {
		log.Fatalln(err)
	}

	var interceptor []grpc.UnaryServerInterceptor
	interceptor = append(interceptor,
		myLogging(),
		grpc_auth.UnaryServerInterceptor(authorize),
	)

	var opts_grpc []grpc.ServerOption
	opts_grpc = append(opts_grpc,
		grpc.Creds(creds),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(interceptor...)),
	)

	s := grpc.NewServer(opts_grpc...)
	pb.RegisterFileServiceServer(s, &server{})
	fmt.Println("server is running...")
	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to Serve: %v", err)
	}
}
