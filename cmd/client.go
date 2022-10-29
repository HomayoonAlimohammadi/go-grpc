package cmd

import (
	"context"
	"fmt"
	"log"
	"pcbook/pb"
	"pcbook/sample"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var (
	clientCmd = &cobra.Command{
		Use:   "client",
		Short: "runs the pcbook client",
		Long: `This command runs the pcbook client.
		Pretend that this is a long description.`,
		Run: func(cmd *cobra.Command, args []string) {
			runClient()
		},
	}
)

func init() {
	rootCmd.AddCommand(clientCmd)
}

func runClient() {
	conn, err := grpc.Dial(fmt.Sprintf(":%s", webPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("can not dial on port:", webPort, err)
	}

	client := pb.NewLaptopServiceClient(conn)

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
			log.Fatal("error creating the laptop")
		}
	}
	log.Println("created laptop with ID:", res.Id)
}
