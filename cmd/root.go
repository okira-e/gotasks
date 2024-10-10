package cmd

import (
	"log"
	"os"
	"runtime"

	"github.com/okira-e/gotasks/cmd/board"
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
		
		userConfig, err := domain.GetUserConfig()
		if err != nil {
			log.Fatalf("Failed to get the user config. %s", err)
		}
		
		// Assume this is not a new board/board.
		// Traverse the current file path until you find a directory with the name of
		// board saved in the config. If not found, create one at the current directory.
		
		originalPwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get the current directory. %s", err)
		}
		
		os := runtime.GOOS
		
		var boardName string
		yeildedPwd := originalPwd
		for yeildedPwd != "" {
			var currentDirName string

			yeildedPwd, currentDirName = getLastDirName(yeildedPwd, byte(utils.Cond(os == "windows", '\\', '/')))
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
			_, boardName = getLastDirName(originalPwd, byte(utils.Cond(os == "windows", '\\', '/')))
			userConfig.CreateBoard(boardName, originalPwd)
		}
		
		app, err := ui.NewApp(userConfig, boardName)
		if err != nil {
			log.Fatalf("Failed to initialize app. %v", err)
		}
		app.Run()
	},
}

func Execute() {
	rootCmd.AddCommand(ListAllBoards)
	rootCmd.AddCommand(OpenConfig)
	rootCmd.AddCommand(board.BoardCmd)
	
	board.BoardCmd.AddCommand(board.OpenBoardByName)

	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("Error executing root command. %s", err.Error())
	}
}


// getLastDirName takes in "/Users/You/Projects/Todo" and returns ("/Users/You/Projects", "Todo").
// Returns an empty string if no "/" was found in the path.
func getLastDirName(path string, pathSeparator byte) (string, string) {
	// Windows
	if path[1:3] == ":\\" {
		path = path[2:]
	}
	
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == pathSeparator {
			return path[:i], path[i+1:]
		}
	}
	
	return "", ""
}