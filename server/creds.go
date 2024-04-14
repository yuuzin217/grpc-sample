package main

import (
	grpc_sample "yuuzin217/grpc-sample"

	"google.golang.org/grpc/credentials"
)

func loadCreds() (credentials.TransportCredentials, error) {
	return credentials.NewServerTLSFromFile(
		grpc_sample.CertFile_filePath,
		grpc_sample.CertKey_filePath,
	)
}
