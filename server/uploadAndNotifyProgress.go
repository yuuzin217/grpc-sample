package main

import (
	"fmt"
	"io"
	"yuuzin217/grpc-sample/pb"
)

func (*server) UploadAndNotifyProgress(stream pb.FileService_UploadAndNotifyProgressServer) error {
	fmt.Println("UploadAndNotifyProgress was invoked.")
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("UploadAndNotifyProgress finished.")
			return nil
		}
		if err != nil {
			return err
		}
		if err := sendSize(stream, request.Data); err != nil {
			return err
		}
	}
}

func sendSize(stream pb.FileService_UploadAndNotifyProgressServer, data []byte) error {
	fmt.Printf("received data: %v\n", data)
	size := 0
	size += len(data)
	return stream.Send(
		&pb.UploadAndNotifyProgressResponse{
			Msg: fmt.Sprintf("received %v bytes", size),
		},
	)
}
