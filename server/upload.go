package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"yuuzin217/grpc-sample/pb"
)

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
