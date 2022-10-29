package api

import (
	"fmt"
	"log"
	"pcstore/pb"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	unaryClientCmd = &cobra.Command{
		Use:   "uclient",
		Short: "runs the pcstore client",
		Long: `This command runs the pcstore client.
		Pretend that this is a long description.`,
		Run: func(cmd *cobra.Command, args []string) {
			runUnaryClient()
		},
	}
)

func init() {
	rootCmd.AddCommand(unaryClientCmd)
}

func runUnaryClient() {
	conn, err := grpc.Dial(fmt.Sprintf(":%s", webPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("can not dial on port:", webPort, err)
	}

	client := pb.NewLaptopServiceClient(conn)

	createLaptop(client)
}
