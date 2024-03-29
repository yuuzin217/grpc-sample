package main

import (
	"fmt"
	"log"

	grpc_sample "yuuzin217/grpc-sample"
	"yuuzin217/grpc-sample/client/caller"
	"yuuzin217/grpc-sample/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	address = grpc_sample.Address
)

func main() {

	connection, client, err := newGRPCConnection()
	if err != nil {
		log.Fatalf("Failed to connect: %v\n", err)
	}
	defer connection.Close()

	// // ListFiles
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

func newGRPCConnection() (*grpc.ClientConn, pb.FileServiceClient, error) {

	certFile := "../ssl/rootCA.pem"
	creds, err := credentials.NewClientTLSFromFile(certFile, "")
	if err != nil {
		return nil, nil, err
	}

	// connection, err := grpc.Dial(address, grpc.WithInsecure())
	connection, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, nil, err
	}

	return connection, pb.NewFileServiceClient(connection), nil
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
