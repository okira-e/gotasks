package utils

import (
	"log"

	"github.com/gizak/termui/v3"
)

// ExitApp closes termui and logs using Fatalf.
func ExitApp(message string) {
	termui.Close()
	
	log.Fatalln(message)
}