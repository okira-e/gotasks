package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "DEV"

var Version = &cobra.Command{
	Use:   "version",
	Short: "Displays the version of gotasks.",
	Long:  `Displays the version of gotasks.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gotasks version: %s\n", version)
	},
}
