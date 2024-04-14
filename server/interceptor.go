package main

import (
	"context"
	"log"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func myLogging() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		log.Printf("request data: %+v", req)

		resp, err = handler(ctx, req)
		if err != nil {
			return nil, err
		}

		log.Printf("response data: %+v", resp)

		return resp, nil
	}
}

func authorize(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}
	if token != "test-token" {
		// return nil, errors.New("bad token")
		return nil, status.Error(codes.Unauthenticated, "token is invalid.")
	}
	return ctx, nil
}

func loadUnaryInterceptor() []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		myLogging(),
		grpc_auth.UnaryServerInterceptor(authorize),
	}
}

func loadStreamInterceptor() []grpc.StreamServerInterceptor {
	return []grpc.StreamServerInterceptor{}
}
