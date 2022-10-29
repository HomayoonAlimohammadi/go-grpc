package server

import (
	"context"
	"errors"
	"log"
	"pcstore/pb"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LaptopServer struct {
	pb.UnimplementedLaptopServiceServer
	store LaptopStore
}

func NewLaptopServer(store LaptopStore) *LaptopServer {
	return &LaptopServer{
		store: store,
	}
}

func (server LaptopServer) CreateLaptop(ctx context.Context, request *pb.CreateLaptopRequest) (*pb.CreateLaptopResponse, error) {
	laptop := request.GetLaptop()
	log.Printf("receive a create-laptop request with id: %s", laptop.Id)

	if len(laptop.Id) > 0 {
		// check if it's a valid UUID
		_, err := uuid.Parse(laptop.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "laptop ID is not a valid UUID: %v", err)
		}
	} else {
		// generate new UUID if not already exists
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot generate a new laptop ID: %v", err)
		}
		laptop.Id = id.String()
	}

	// check for context cancelation or deadline exceeded
	if isContextCanceled(ctx) {
		return nil, ErrorContextCanceled
	}
	if isContextDeadlineExceeded(ctx) {
		return nil, ErrorContextDeadlineExceeded
	}

	// try saving the laptop
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

func (server *LaptopServer) SearchLaptop(req *pb.SearchLaptopRequest, stream pb.LaptopService_SearchLaptopServer) error {
	filter := req.GetFilter()
	log.Printf("recieve a SearchLaptop request with filter: %v", filter)

	err := server.store.Search(
		stream.Context(),
		filter,
		func(laptop *pb.Laptop) error {
			res := &pb.SearchLaptopResponse{Laptop: laptop}
			err := stream.Send(res)
			if err != nil {
				return err
			}

			log.Println("sent laptop with id:", laptop.Id)
			return nil
		},
	)
	if err != nil {
		return status.Errorf(codes.Internal, "unexpected error: %v", err)
	}
	return nil
}
