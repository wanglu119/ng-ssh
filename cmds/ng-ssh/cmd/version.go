package cmd

import (
	"fmt"
	
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command {
	Use: "version",
	Short: "Print the version number of auth_client",
	Long: "All software has version. This is auth_client's",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("0.1.0")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
