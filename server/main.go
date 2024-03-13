package main

import (
	"context"
	"fmt"
	grpcsample "grpc-sample"
	"grpc-sample/pb"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

// var _ *pb.FileServiceServer = (*server)(nil)

type Server struct {
	pb.UnimplementedFileServiceServer
}

const (
	dir_storage = "../storage"
	address     = grpcsample.Address
)

func (*Server) ListFiles(ctx context.Context, req *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {

	fmt.Println("ListFiles was invoked.")

	paths, err := os.ReadDir(dir_storage)
	if err != nil {
		return nil, err
	}

	var fileNames []string
	for _, path := range paths {
		if !path.IsDir() {
			fileNames = append(fileNames, path.Name())
		}
	}

	return &pb.ListFilesResponse{
		FileNames: fileNames,
	}, nil

}

func main() {
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterFileServiceServer(s, &Server{})

	fmt.Println("server is running...")
	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to Serve: %v", err)
	}
}
