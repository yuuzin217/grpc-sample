package main

import (
	"bytes"
	"context"
	"fmt"
	grpcsample "grpc-sample"
	"grpc-sample/pb"
	"io"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedFileServiceServer
}

const (
	address = grpcsample.Address

	dir_storage        = grpcsample.Dir_storage
	dir_storage_remote = grpcsample.Dir_storage_remote
	dir_storage_local  = grpcsample.Dir_storage_local

	uploaded_text = "uploaded_text.txt"
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

func (*server) Download(request *pb.DownloadRequest, stream pb.FileService_DownloadServer) error {
	fmt.Println("Download was invoked.")
	file, err := os.Open(fmt.Sprint(dir_storage, request.FileName))
	if err != nil {
		return err
	}
	defer file.Close()
	buf := make([]byte, 5)
	for {
		n, err := file.Read(buf)
		// データが何も読み込まれなかった or ファイルの終端まで到達した
		if n == 0 || err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if err := stream.Send(&pb.DownloadResponse{Data: buf[:n]}); err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}
}

func (*server) Upload(stream pb.FileService_UploadServer) error {
	fmt.Println("Upload was invoked.")
	var buf bytes.Buffer
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			file, err := os.Create(fmt.Sprint(dir_storage_remote, uploaded_text))
			if err != nil {
				return err
			}
			defer file.Close()
			if _, err := file.Write(buf.Bytes()); err != nil {
				return err
			}
			return stream.SendAndClose(&pb.UploadResponse{Size: int32(buf.Len())})
		}
		if err != nil {
			return err
		}
		buf.Write(request.Data)
	}
}

func main() {
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterFileServiceServer(s, &server{})
	fmt.Println("server is running...")
	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to Serve: %v", err)
	}
}
