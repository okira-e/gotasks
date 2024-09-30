package board

import (
	"fmt"
	"log"

	"github.com/okira-e/gotasks/internal/domain"
	"github.com/okira-e/gotasks/internal/ui"
	"github.com/spf13/cobra"
)

var OpenBoardByName= &cobra.Command{
	Use:   "open",
	Short: "Open a specific board by name",
	Long:  `Open a specific board by name.`,
	Run: func(cmd *cobra.Command, args []string) {
		
		if len(args) == 0 {
			fmt.Println("Please provide a board name to open.")
			return
		}
		
		boardName := args[0]
		
		config, err := domain.GetUserConfig()
		if err != nil {
			log.Fatalf("Failed to get a userConfig instance. ", err)
		}
		
		for _, board := range config.Boards {
			if board.Name == boardName {
				userConfig, err := domain.GetUserConfig()
				if err != nil {
					log.Fatalf("Failed to setup the user config. %s", err)
				}
				
				app, err := ui.NewApp(userConfig, boardName)
				if err != nil {
					log.Fatalf("Failed to initialize app. %v", err)
				}
				app.Run()
				break
			}
		}
		
		fmt.Println("Couldn't find the board you tried to open.")
		fmt.Println("Run \"gotasks list\" to view all available boards.")
		return
	},
}
