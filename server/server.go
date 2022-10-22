package server

import (
	"context"
	"errors"
	"log"
	"pcbook/pb"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LaptopServer struct {
	store LaptopStore
}

func NewLaptopServer(store LaptopStore) *LaptopServer {
	return &LaptopServer{
		store: store,
	}
}

func (server *LaptopServer) CreateLaptop(ctx context.Context, request *pb.CreateLaptopRequest) (*pb.CreateLaptopResponse, error) {
	laptop := request.GetLaptop()
	log.Printf("receive a create-laptop request with id: %s", laptop.Id)

	if len(laptop.Id) > 0 {
		// check if it's a valid UUID
		_, err := uuid.Parse(laptop.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "laptop ID is not a valid UUID: %v", err)
		}
	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot generate a new laptop ID: %v", err)
		}
		laptop.Id = id.String()
	}
	if ctx.Err() == context.Canceled {
		log.Print("request is canceled")
		return nil, status.Error(codes.Canceled, "request is canceled")
	}

	if ctx.Err() == context.DeadlineExceeded {
		log.Print("deadline is exceeded")
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}

	err := server.store.Save(laptop)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, ErrorLaptopAlreadyExists) {
			code = codes.AlreadyExists
		}
		return nil, status.Errorf(code, "can not save laptop in the store: %v", err)
	}

	log.Printf("Saved laptop %s in the store.\n", laptop.Id)

	res := &pb.CreateLaptopResponse{
		Id: laptop.Id,
	}
	return res, nil
}
