package cmd

import (
	"log"
	"os"

	"github.com/okira-e/gotasks/internal/domain"
	"github.com/okira-e/gotasks/internal/ui"
	"github.com/okira-e/gotasks/internal/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gotasks",
	Short: "This is a tool for migrating database environments.",
	Long: `
gotasks is a kanban board in the TUI. Think of it as a tool that sits
between a to-do list and a Jira board that is accessible from the terminal.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		var userConfig *domain.UserConfig
		
		doesUserConfigExists, err := domain.DoesUserConfigExist()
		if err != nil {
			log.Fatalf("Failed to check if user config exists. %v", err)
		}
		
		if !doesUserConfigExists {
			userConfig, err = domain.SetupUserConfig()
			if err != nil {
				log.Fatalf("Failed to setup the user config. %s", err)
			}
		} else {
			userConfig, err = domain.GetUserConfig()
			if err != nil {
				log.Fatalf("Failed to setup the user config. %s", err)
			}
		}
		
		// Assume this is not a new board/board.
		// Traverse the current file path until you find a directory with the name of
		// board saved in the config. If not found, create one at the current directory.
		
		originalPwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get the current directory. %s", err)
		}
		
		var boardName string
		yeildedPwd := originalPwd
		for yeildedPwd != "" {
			var currentDirName string

			yeildedPwd, currentDirName = getLastDirName(yeildedPwd)
			if currentDirName == "" {
				log.Fatalf("Failed to get the directory name.")
			}
			
			boardOpt := userConfig.GetBoard(currentDirName)
			if boardOpt.IsSome() {
				boardName = currentDirName
				break
			} else {
				continue
			}
		}
		
		if boardName == "" {
			_, boardName = getLastDirName(originalPwd)
			userConfig.AddBoard(boardName, originalPwd)
		}
		
		// Initialize the App
		app, err := ui.NewApp(userConfig, boardName)
		if err != nil {
			log.Fatalf("Failed to initialize app. %v", err)
		}
		
		utils.SaveLog(utils.Info, "Initialized App", nil)

		app.Run()
	},
}

func Execute() {
	rootCmd.AddCommand(ListAllBoards)
	rootCmd.AddCommand(OpenConfig)

	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("Error executing root command. %s", err.Error())
	}
}


// getLastDirName takes in "/Users/You/Projects/Todo" and returns ("/Users/You/Projects", "Todo").
// Returns an empty string if no "/" was found in the path.
func getLastDirName(path string) (string, string) {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			return path[:i], path[i+1:]
		}
	}
	
	return "", ""
}