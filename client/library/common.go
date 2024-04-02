package caller

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/metadata"
)

func newMDData(mdKey, mdVal string) map[string]string {
	return map[string]string{mdKey: mdVal}
}

func newBearerAuthContext(token string) context.Context {
	mdKey := "authorization"
	mdVal := fmt.Sprint("bearer ", token)
	md := metadata.New(newMDData(mdKey, mdVal))
	return metadata.NewOutgoingContext(context.Background(), md)
}

func newContextWithTimeout(second int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(second)*time.Second)
}
