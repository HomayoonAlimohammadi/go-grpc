package api

import (
	"fmt"

	"github.com/spf13/cobra"
)

var serveStreamCmd = &cobra.Command{
	Use:   "stream",
	Short: "start streaming grpc server",
	Long: `Start streaming grpc server.
	Pretend that this is a long description.`,
	Run: func(cmd *cobra.Command, args []string) {
		runStreamingServer()
	},
}

func init() {
	rootCmd.AddCommand(serveStreamCmd)
}

func runStreamingServer() {
	fmt.Println("starting streaming server on port:", webPort)
}
