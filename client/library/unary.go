package caller

import (
	"yuuzin217/grpc-sample/pb"
)

func CallListFiles(client pb.FileServiceClient) ([]string, error) {
	ctx := newBearerAuthContext("test-token")
	result, err := client.ListFiles(ctx, &pb.ListFilesRequest{})
	if err != nil {
		return nil, err
	}
	return result.FileNames, nil
}
