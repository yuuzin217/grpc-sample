package main

import (
	"fmt"
	"io"
	"yuuzin217/grpc-sample/pb"
)

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
