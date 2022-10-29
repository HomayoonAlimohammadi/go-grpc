package api

import (
	"context"
	"io"
	"log"
	"pcstore/pb"
	"pcstore/sample"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func createLaptop(client pb.LaptopServiceClient) {
	laptop := sample.NewLaptop()
	req := pb.CreateLaptopRequest{
		Laptop: laptop,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := client.CreateLaptop(ctx, &req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			log.Println("laptop already exists")
		} else {
			log.Fatal("error creating the laptop", err)
		}
	}
	log.Println("created laptop with ID:", res.Id)
}

func searchLaptop(client pb.LaptopServiceClient, filter *pb.Filter) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.SearchLaptopRequest{Filter: filter}
	stream, err := client.SearchLaptop(ctx, req)
	if err != nil {
		log.Fatal("error searching for valid laptops:", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatal("error receiving search response from the server:", err)
		}
		laptop := res.GetLaptop()
		printLaptop(laptop)
	}
}

func printLaptop(laptop *pb.Laptop) {
	log.Printf("\033[1;33m%s\033[0m", ">>>>>>>>>>>>>>>>>>")
	log.Println("- found: ", laptop.GetId())
	log.Println("  + brand: ", laptop.GetBrand())
	log.Println("  + name: ", laptop.GetName())
	log.Println("  + cpu cores: ", laptop.GetCpu().GetNumberCores())
	log.Println("  + cpu min ghz: ", laptop.GetCpu().GetMinGhz())
	log.Println("  + ram: ", laptop.GetRam())
	log.Println("  + price: ", laptop.GetPriceUsd())
	log.Printf("\033[1;33m%s\033[0m", ">>>>>>>>>>>>>>>>>>\n")
}
