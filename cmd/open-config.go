package cmd

import (
	"log"
	"os"

	"github.com/okira-e/gotasks/internal/domain"
	"github.com/okira-e/gotasks/internal/utils"
	"github.com/okira-e/gotasks/internal/vars"
	"github.com/spf13/cobra"
)

var OpenConfig = &cobra.Command{
	Use:   "config",
	Short: "Opens the config file",
	Long:  `Opens the config file for all projects in Vi.`,
	Run: func(cmd *cobra.Command, args []string) {
		path, err := domain.GetConfigFilePathBasedOnOS()
		if err != nil {
			log.Fatalln("Failed to open up the config. ", err)
		}
		
		editor := os.Getenv(vars.EditorOfChoice)
		if editor == "" {
			editor = "vi"
		}
		
		err = utils.OpenInEditor(editor, path)
		if err != nil {
			log.Fatalf("Failed to open the config file in %s. %s\n", editor, err)
		}
	},
}
