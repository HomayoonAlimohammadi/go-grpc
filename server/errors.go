package server

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrorContextCanceled         = status.Error(codes.Canceled, "request is canceled")
	ErrorContextDeadlineExceeded = status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	ErrorLaptopAlreadyExists     = errors.New("laptop already exists")
	ErrorCanNotCopyLaptop        = errors.New("can not copy laptop")
)
