package utils

import (
	"encoding/json"
	"os"

	"github.com/okira-e/gotasks/internal/vars"
	"github.com/sirupsen/logrus"
)

type Severity int

const (
	Info Severity = iota
	Debug
	Warn
	Error
	Fatal
)


func SaveLog(severity Severity, message string, context map[string]any) {
	contextJSON, err := json.Marshal(context)
	if err != nil {
		contextJSON = []byte("ERROR: Couldn't parse the data for this log.")
	}
	
	payload := logrus.WithFields(
		logrus.Fields{
			"context": string(contextJSON),
		},
	)
	
	switch severity {
		case Info:
		{
			payload.Info(message)
		}
		case Debug:
		{
			isDebugMode := os.Getenv(vars.DebugFlag)
			
			if isDebugMode == "true" {
				payload.Debug(message)
			}
		}
		case Warn:
		{
			payload.Warn(message)
		}
		case Error:
		{
			payload.Error(message)
		}
		case Fatal:
		{
			payload.Fatal(message)
		}
	}
}