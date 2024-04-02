package caller

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"
	"yuuzin217/grpc-sample/pb"
)

func CallUpload(client pb.FileServiceClient) error {
	filePath := fmt.Sprint(dir_storage_local, file_local)
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	ctx := context.Background()
	stream, err := client.Upload(ctx)
	if err != nil {
		return err
	}
	buf := make([]byte, 1*kb)
	for {
		n, err := file.Read(buf)
		if n == 0 || err == io.EOF {
			return closeAndRecv(stream)
		}
		if err != nil {
			return err
		}
		if err := stream.Send(newUploadRequest(buf, n)); err != nil {
			return err
		}
		log.Println("uploading...")
		time.Sleep(1 * time.Second)
	}
}

func newUploadRequest(buf []byte, cap int) *pb.UploadRequest {
	return &pb.UploadRequest{
		Data: buf[:cap],
	}
}

func closeAndRecv(stream pb.FileService_UploadClient) error {
	res, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	log.Printf("received data size: %v\n", res.Size)
	return nil
}
