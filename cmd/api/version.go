package api

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "1.0.0"
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "display the version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("PC-Store Version %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
