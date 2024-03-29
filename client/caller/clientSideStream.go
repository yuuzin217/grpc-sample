package caller

import (
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
	ctx := newContext()
	stream, err := client.Upload(ctx)
	if err != nil {
		return err
	}
	buf := make([]byte, 1*kb)
	for {
		n, err := file.Read(buf)
		if n == 0 || err == io.EOF {
			response, err := stream.CloseAndRecv()
			if err != nil {
				return err
			}
			log.Printf("received data size: %v", response.Size)
			return nil
		}
		if err != nil {
			return err
		}
		if err := stream.Send(&pb.UploadRequest{Data: buf[:n]}); err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}
}
