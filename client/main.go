package main

import (
	"context"
	"fmt"
	grpcsample "grpc-sample"
	"grpc-sample/pb"
	"log"

	"google.golang.org/grpc"
)

const (
	address = grpcsample.Address
)

func main() {
	connection, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer connection.Close()
	fileNames, err := callListFiles(pb.NewFileServiceClient(connection))
	if err != nil {
		log.Fatalf("Failed to Request ListFiles: %v", err)
	}
	fmt.Println(fileNames)
}

func callListFiles(client pb.FileServiceClient) ([]string, error) {
	result, err := client.ListFiles(context.Background(), &pb.ListFilesRequest{})
	if err != nil {
		return nil, err
	}
	return result.FileNames, nil
}
