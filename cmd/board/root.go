package board

import (
	"github.com/spf13/cobra"
)

var BoardCmd = &cobra.Command{
	Use:   "board",
	Short: "Perform an operation on a specific board",
	Long:  `Perform an operation on a specific board.`,
}
