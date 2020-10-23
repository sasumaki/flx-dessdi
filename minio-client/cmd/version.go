package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of AIGA",
	Long:  `All software has versions. This is AIGA`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("AIGA version v0.1alpha")
	},
}
