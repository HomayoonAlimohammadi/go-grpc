package api

import (
	"fmt"
	"log"
	"net"
	"pcbook/pb"
	"pcbook/server"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var (
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "run the grpc server",
		Long: `This command sets up the server.
	Pretend that this is a long description.`,
		Run: func(cmd *cobra.Command, args []string) {
			serve()
		},
	}
)

func serve() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", webPort))
	if err != nil {
		log.Fatal("error listening on port:", webPort, err)
	}

	grpcServer := grpc.NewServer()
	laptopStore := server.NewInMemoryLaptopStore()
	laptopServer := server.NewLaptopServer(laptopStore)

	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	log.Println("starting the server on port:", webPort)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal("error starting the server on port:", webPort, err)
	}
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
