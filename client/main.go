package main

import (
	"bytes"
	"context"
	"fmt"
	grpcsample "grpc-sample"
	"grpc-sample/pb"
	"io"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
)

const (
	address = grpcsample.Address

	dir_storage_local  = grpcsample.Dir_storage_local
	dir_storage_remote = grpcsample.Dir_storage_remote

	file_local      = "local_text.txt"
	file_remote     = "remote_text.txt"
	downloaded_text = "downloaded_text.txt"

	kb = grpcsample.KB
)

func main() {
	connection, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer connection.Close()

	client := pb.NewFileServiceClient(connection)

	// ListFiles
	fileNames, err := callListFiles(client)
	if err != nil {
		log.Fatalf("Failed to request ListFiles: %v", err)
	}
	fmt.Println(fileNames)

	// Download
	if err := callDownload(client); err != nil {
		log.Fatalf("Failed to request Download: %v", err)
	}

	// Upload
	if err := callUpload(client); err != nil {
		log.Fatalf("Failed to request Upload: %v", err)
	}
}

func callListFiles(client pb.FileServiceClient) ([]string, error) {
	result, err := client.ListFiles(context.Background(), &pb.ListFilesRequest{})
	if err != nil {
		return nil, err
	}
	return result.FileNames, nil
}

func callDownload(client pb.FileServiceClient) error {
	stream, err := client.Download(context.Background(), &pb.DownloadRequest{FileName: fmt.Sprint(dir_storage_remote, file_remote)})
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			file, err := os.Create(fmt.Sprint(dir_storage_local, downloaded_text))
			if err != nil {
				return err
			}
			defer file.Close()
			if _, err := file.Write(buf.Bytes()); err != nil {
				return err
			}
			log.Println("download finished.")
			return nil
		}
		if err != nil {
			return err
		}
		buf.Write(response.Data)
	}
}

func callUpload(client pb.FileServiceClient) error {
	file, err := os.Open(fmt.Sprint(dir_storage_local, file_local))
	if err != nil {
		return err
	}
	defer file.Close()
	stream, err := client.Upload(context.Background())
	if err != nil {
		return err
	}
	buf := make([]byte, 1*kb)
	for {
		n, err := file.Read(buf)
		if n == 0 || err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if err := stream.Send(&pb.UploadRequest{Data: buf[:n]}); err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}
	response, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	log.Printf("received data size: %v", response.Size)
	return err
}
