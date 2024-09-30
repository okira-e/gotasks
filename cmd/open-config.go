package cmd

import (
	"log"
	"os"
	"os/exec"

	"github.com/okira-e/gotasks/internal/domain"
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
		
		err = openInEditor(editor, path)
		if err != nil {
			log.Fatalf("Failed to open the config file in %s. %s\n", editor, err)
		}
	},
}


func openInEditor(editor string, filePath string) error {
	
	// Execute the 'vi' command to open the config file
	cmd := exec.Command(editor, filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Run the command and return any error encountered
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}