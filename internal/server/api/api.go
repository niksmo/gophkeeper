package api

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrInternal = status.Error(codes.Internal, "internal error")
