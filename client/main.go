package main

import (
	"fmt"
	"log"

	grpc_sample "yuuzin217/grpc-sample"
	caller "yuuzin217/grpc-sample/client/library"
	"yuuzin217/grpc-sample/pb"
)

func main() {

	connection, client, err := newGRPCConnection(grpc_sample.Host, false)
	if err != nil {
		log.Fatalf("Failed to gRPC connect: %v\n", err)
	}
	defer connection.Close()

	// ListFiles
	fileNames, err := CallListFiles(client)
	if err != nil {
		log.Fatalf("Failed to request ListFiles: %v", err)
	}
	fmt.Println(fileNames)

	// Download
	if err := CallDownload(client); err != nil {
		log.Fatalf("Failed to request Download: %v", err)
	}

	// Upload
	if err := CallUpload(client); err != nil {
		log.Fatalf("Failed to request Upload: %v", err)
	}

	// bidirectional
	if err := CallUploadAndNotifyProgress(client); err != nil {
		log.Fatalf("Failed to request two-way: %v", err)
	}

}

func CallListFiles(client pb.FileServiceClient) ([]string, error) {
	return caller.CallListFiles(client)
}

func CallDownload(client pb.FileServiceClient) error {
	return caller.CallDownload(client)
}

func CallUpload(client pb.FileServiceClient) error {
	return caller.CallUpload(client)
}

func CallUploadAndNotifyProgress(client pb.FileServiceClient) error {
	return caller.CallUploadAndNotifyProgress(client)
}
