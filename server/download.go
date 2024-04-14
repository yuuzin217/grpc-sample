package main

import (
	"fmt"
	"io"
	"os"
	"time"
	"yuuzin217/grpc-sample/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
		if n == 0 || err == io.EOF {
			// データが何も読み込まれなかった or ファイルの終端まで到達した
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
