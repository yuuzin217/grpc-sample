package caller

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"
	"yuuzin217/grpc-sample/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CallDownload(client pb.FileServiceClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	stream, err := client.Download(ctx, &pb.DownloadRequest{FileName: file_remote})
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			if err := createDownloadedFile(buf); err != nil {
				return err
			}
			log.Println("download finished.")
			return nil
		}
		if grpcErr, ok := status.FromError(err); ok {
			if err := gRPCErrorCheck(grpcErr); err != nil {
				return err
			}
		}
		if err != nil {
			return err
		}
		log.Println("loading...")
		buf.Write(response.Data)
	}
}

func createDownloadedFile(buf bytes.Buffer) error {
	filePath := fmt.Sprint(dir_storage_local, downloaded_text)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.Write(buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func gRPCErrorCheck(grpcErr *status.Status) error {
	switch grpcErr.Code() {
	case codes.NotFound:
		return fmt.Errorf("error code: %v, error message: %v", grpcErr.Code(), grpcErr.Message())
	case codes.DeadlineExceeded:
		return errors.New("deadline exceeded")
	default:
		// 未ハンドリングの gRPC のエラー
		return nil
	}
}
