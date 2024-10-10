package main

import (
	"fmt"
	"log"
	"os"

	"github.com/okira-e/gotasks/cmd"
	"github.com/okira-e/gotasks/internal/domain"
	"github.com/sirupsen/logrus"
)

type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format("2006-01-02T15:04:05.000")
	
	level := fmt.Sprintf("[%s]", entry.Level.String())
	level = level[0:1] + level[1:] // Ensure it's uppercase

	msg := fmt.Sprintf("%s %s %s - %s\n", timestamp, level, entry.Message, entry.Data)

	return []byte(msg), nil
}

func init() {
	doesUserConfigExists, err := domain.DoesUserConfigExist()
	if err != nil {
		log.Fatalf("Failed to check if user config exists. %v", err)
	}
	
	if !doesUserConfigExists {
		_, err = domain.SetupUserConfig()
		if err != nil {
			log.Fatalf("Failed to setup the user config. %s", err)
		}
	}
	
	configPath, err := domain.GetConfigDirPathBasedOnOS()
	if err != nil {
		log.Fatalf("Couldn't get the config path to write logs at. %s", err)
	}
	
	file, err := os.OpenFile(configPath + "/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
	    logrus.SetOutput(file)
	} else {
	    logrus.Info("Failed to log to file, using default stderr")
	}

	logrus.SetLevel(logrus.DebugLevel)

	logrus.SetFormatter(&CustomFormatter{})
}

func main() {
	cmd.Execute()
}
