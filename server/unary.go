package main

import (
	"context"
	"fmt"
	"os"
	"yuuzin217/grpc-sample/pb"
)

func (*server) ListFiles(ctx context.Context, req *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
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
