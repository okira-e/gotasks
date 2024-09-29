package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/okira-e/gotasks/internal/domain"
	"github.com/spf13/cobra"
)

var ListAllBoards = &cobra.Command{
	Use:   "list",
	Short: "List all the boards.",
	Long:  `List all boards saved in the user config.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := domain.GetUserConfig()
		if err != nil {
			log.Fatalf("Failed to get the user config. %s", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"#", "Board", "Number of tasks", "Progress", "Path"})
		
		for i, board := range config.Boards {
			totalNumberOfTasks := 0
			numberOfCompletedTasks := 0
			for i, columnName := range board.Columns {
				totalNumberOfTasks += len(board.Tasks[columnName])
				
				if i == len(board.Columns) - 1 { // Check if we are on the last column to the right.
					numberOfCompletedTasks += len(board.Tasks[columnName])					
				}
			}
			
			progress := float32(0)
			if totalNumberOfTasks != 0 {
				progress = float32(numberOfCompletedTasks) / float32(totalNumberOfTasks)
			}
			
			t.AppendRow([]any{
				i + 1, 
				board.Name, 
				totalNumberOfTasks,
				fmt.Sprint("%", int(progress * 100)),
				board.Dir,
			})
		}
		t.AppendSeparator()
		
		t.Render()
	},
}
