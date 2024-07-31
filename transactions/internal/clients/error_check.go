package clients

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func isTransportError(err error) bool {
	return status.Code(err) == codes.Unavailable
}
