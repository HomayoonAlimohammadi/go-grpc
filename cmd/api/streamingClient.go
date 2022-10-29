package api

import (
	"fmt"
	"log"
	"pcstore/pb"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var streamingClientCmd = &cobra.Command{
	Use:   "sclient",
	Short: "starts a streaming client",
	Long: `This command starts a streaming client.
	Pretend like this is a long description!`,
	Run: func(cmd *cobra.Command, args []string) {
		runStreamingClient()
	},
}

func init() {
	rootCmd.AddCommand(streamingClientCmd)
}

func runStreamingClient() {
	conn, err := grpc.Dial(fmt.Sprintf(":%s", webPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("can not dial on port:", webPort, err)
	}

	client := pb.NewLaptopServiceClient(conn)
	for i := 0; i < 10; i++ {
		createLaptop(client)
	}

	filter := &pb.Filter{
		MaxPriceUsd: 3000,
		MinCpuCores: 4,
		MinCpuGhz:   2.5,
		MinRam:      &pb.Memory{Value: 8, Unit: pb.Memory_GIGABYTE},
	}
	searchLaptop(client, filter)
}
