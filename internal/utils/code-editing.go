package utils

import (
	"os"
	"os/exec"
)

func OpenInEditor(editor string, filePath string) error {
	
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