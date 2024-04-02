package caller

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
	"yuuzin217/grpc-sample/pb"
)

func CallUploadAndNotifyProgress(client pb.FileServiceClient) error {
	filePath := fmt.Sprint(dir_storage_local, file_local)
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	ctx := context.Background()
	stream, err := client.UploadAndNotifyProgress(ctx)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go request(&wg, stream, file)
	go response(&wg, stream)
	wg.Wait()
	return nil
}

func request(wg *sync.WaitGroup, stream pb.FileService_UploadAndNotifyProgressClient, file *os.File) {
	defer wg.Done()
	defer stream.CloseSend()
	buf := make([]byte, 1*kb)
	for {
		n, err := file.Read(buf)
		if n == 0 || err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		if err := stream.Send(&pb.UploadAndNotifyProgressRequest{Data: buf[:n]}); err != nil {
			log.Fatalln(err)
		}
		time.Sleep(1 * time.Second)
	}
}

func response(wg *sync.WaitGroup, stream pb.FileService_UploadAndNotifyProgressClient) {
	defer wg.Done()
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("received message: %v", response.Msg)
	}
}
