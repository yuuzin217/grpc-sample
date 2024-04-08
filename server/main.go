package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
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
	path := fmt.Sprint(dir_storage_remote, request.FileName)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return status.Error(codes.NotFound, "file was not found.")
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	buf := make([]byte, 1*kb)
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

func (*server) UploadAndNotifyProgress(stream pb.FileService_UploadAndNotifyProgressServer) error {
	fmt.Println("UploadAndNotifyProgress was invoked.")
	var size int
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("UploadAndNotifyProgress finished.")
			return nil
		}
		if err != nil {
			return err
		}
		data := request.Data
		fmt.Printf("received data: %v\n", data)
		size += len(data)
		if err := stream.Send(
			&pb.UploadAndNotifyProgressResponse{
				Msg: fmt.Sprintf("received %v bytes", size),
			},
		); err != nil {
			return err
		}
	}
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
