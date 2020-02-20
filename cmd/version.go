package cmd

import (
	"fmt"

	"github.com/shin1x1/docker-base-chaser/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	Long:  "Show version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.Text(CommandName))
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
