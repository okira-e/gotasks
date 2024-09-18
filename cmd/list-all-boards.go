package cmd

import (
	"fmt"
	"log"

	"github.com/okira-e/gotasks/internal/domain"
	"github.com/okira-e/gotasks/internal/utils"
	"github.com/okira-e/gotasks/internal/vars"
	"github.com/spf13/cobra"
)

var ListAllBoards = &cobra.Command{
	Use:   "list",
	Short: "List all the boards.",
	Long:  `List all boards saved in the user config.`,
	Run: func(cmd *cobra.Command, args []string) {
		// @Todo: Make this list every board in the config including the directory of the board that it was created at.
		
		userConfig, err := domain.GetUserConfig()
		if err != nil {
			log.Fatalf("Failed to get the user config. %s", err.Error())
		}
		
		if len(userConfig.Boards) == 0 {
			utils.PrintInColor(vars.GreenColor, "No boards were found.", false)
			return
		}
		
		utils.PrintInColor(vars.GreenColor, "Boards:", false)
		for _, board := range userConfig.Boards {
			fmt.Println("- ", board)
		}
	},
}
