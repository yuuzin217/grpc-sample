package main

import (
	grpc_sample "yuuzin217/grpc-sample"
	"yuuzin217/grpc-sample/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

/*
gRPC Connection
*/
func newGRPCConnection(host string, insecure bool) (*grpc.ClientConn, pb.FileServiceClient, error) {
	creds, err := getTLSCreds()
	if err != nil {
		return nil, nil, err
	}
	var opts []grpc.DialOption
	if insecure {
		opts = append(opts, grpc.WithInsecure())
	} else {
		opts = append(opts, grpc.WithTransportCredentials(creds))
	}
	connection, err := grpc.Dial(host, opts...)
	if err != nil {
		return nil, nil, err
	}
	return connection, pb.NewFileServiceClient(connection), nil
}

func getTLSCreds() (credentials.TransportCredentials, error) {
	certFile := grpc_sample.RootCA_FilePath
	return credentials.NewClientTLSFromFile(certFile, "")
}
