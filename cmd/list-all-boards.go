package cmd

import (
	"github.com/spf13/cobra"
)

var ListAllBoards = &cobra.Command{
	Use:   "list",
	Short: "List all the boards.",
	Long:  `List all boards saved in the user config.`,
	Run: func(cmd *cobra.Command, args []string) {
		// @Todo: Make this list every board in the config including the directory of the board that it was created at.
	},
}
