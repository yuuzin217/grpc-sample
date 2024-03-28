package caller

import (
	"context"
	"yuuzin217/grpc-sample/pb"

	"google.golang.org/grpc/metadata"
)

func CallListFiles(client pb.FileServiceClient) ([]string, error) {
	md := metadata.New(map[string]string{"authorization": "bearer test-token"})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	result, err := client.ListFiles(ctx, &pb.ListFilesRequest{})
	if err != nil {
		return nil, err
	}
	return result.FileNames, nil
}
