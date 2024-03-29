package caller

import (
	"context"
)

func newContext() context.Context {
	return context.Background()
}

// func commonError(err error) {
// 	log.Fatalln("Error: ", err)
// }
