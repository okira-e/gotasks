package domain

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"runtime"

	"github.com/okira-e/gotasks/internal/opt"
)

/*
{
	"boards": [
		{
			"name": "masa",
			"dir": "/Users/omarrafat/Boards/masa",
			"board": {
				"columns": ["Todo", "Open", "Closed"],
				"tasks": {
					"Todo": [
						{
							"title": "Lorem",
							"description": "Some optional Lorem Epison"
						},
						{
							"title": "Lorem",
							"description": "Some optional Lorem Epison"
						}
					],
					"Open": [
						{
							"title": "Lorem",
							"description": "Some optional Lorem Epison"
						}
					],
					"Closed": [
						{
							"title": "Lorem",
							"description": "Some optional Lorem Epison"
						}
					]
				}
			}
		}
	]
}
*/


type UserConfig struct {
	Boards []*Board `json:"boards"`
}

type Board struct {
	Name  string `json:"name"`
	Dir   string `json:"dir"`
	Columns []string          `json:"columns"`
	// Tasks are the individual cards on the board representing a task.
	Tasks map[string][]Task `json:"tasks"`
}

// DoesUserConfigExist checks if a user config has already be generated for this user.
func DoesUserConfigExist() (bool, error) {
	filePath, err := getConfigFilePathBasedOnOS()
	if err != nil {
		return false, err
	}
	
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false, nil
	}

	return true, nil
}

// SetupUserConfig creates a new default config and writes it to disk.
// It returns a pointer to the new config.
func SetupUserConfig() (*UserConfig, error) {
	config := NewDefaultUserConfig()
	
	err := config.writeToDisk()
	if err != nil {
		return nil, err
	}
	
	return &config, nil
}

// GetUserConfig reads the user config file and returns a pointer 
// to a UserConfig object.
func GetUserConfig() (*UserConfig, error) {
	var userConfig UserConfig

	filePath, err := getConfigFilePathBasedOnOS()
	if err != nil {
		return nil, err
	}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	
	err = json.Unmarshal(fileContent, &userConfig)
	if err != nil {
		return nil, err
	}

	return &userConfig, nil
}

func NewDefaultUserConfig() UserConfig {
	return UserConfig{
		Boards: []*Board{},
	}
}

// AddBoard adds a new board to the config.
func (self *UserConfig) AddBoard(boardName string, dirPath string) error {
	board := Board {
		Name: boardName,
		Dir: dirPath,
		Columns: []string{},
		Tasks: make(map[string][]Task),
	}
	
	// Add the newly created board to the config.
	self.Boards = append(self.Boards, &board)
	self.writeToDisk()
	
	return nil
}

// GetBoard searches the config for a board with the given name.
func (self *UserConfig) GetBoard(boardName string) opt.Option[*Board] {
	for _, it := range self.Boards {
		if it.Name == boardName {
			return opt.Some(it)
		}
	}
	
	return opt.None[*Board]()
}

// UpdateBoard finds and updates the given board in the user config.
func (self *UserConfig) UpdateBoard(board *Board) error {
	for i, it := range self.Boards {
		if it.Name == board.Name {
			self.Boards[i] = board
		}
	}
	
	err := self.writeToDisk()
	if err != nil {
		log.Fatalf("Failed to write to the user config on board update. %s", err)
	}
	
	return nil
}

// AddColumnToBoard adds a column to the board of the board with the given name.
func (self *UserConfig) AddColumnToBoard(boardName string, columnName string) error {
	var board *Board = nil
	
	boardOpt := self.GetBoard(boardName)
	board = boardOpt.Expect("Failed to find board whole adding a column.")

	board.Columns = append(board.Columns, columnName)
	
	err := self.UpdateBoard(board)
	if err != nil {
		return err
	}
	
	return nil
}

// writeToDisk writes or creates the config files with the provided user config.
func (self UserConfig) writeToDisk() error { 
	filePath, err := getConfigFilePathBasedOnOS()
	if err != nil {
		return err
	}

	file, _ := os.Create(filePath)
	defer file.Close()

	fileContent, err := json.MarshalIndent(self, "", "\t")
	if err != nil {
		return err
	}

	_, err = file.Write(fileContent)
	if err != nil {
		return err
	}
	
	return nil
}

// getCOnfigFilePathBasedOnOS returns the config file path based on the OS.
func getConfigFilePathBasedOnOS() (string, error) {
	var osUserName string

	if runtime.GOOS == "windows" {
		osUserName = os.Getenv("USERNAME")
		return "C:\\Users\\" + osUserName + "\\AppData\\Roaming\\gotasks\\config.json", nil
	} else if runtime.GOOS == "darwin" {
		osUserName = os.Getenv("USER")
		return "/Users/" + osUserName + "/Library/Application Support/gotasks/config.json", nil
	} else if runtime.GOOS == "linux" {
		osHomeDir := os.Getenv("HOME")
		return osHomeDir + "/.config/gotasks/config.json", nil
	} else {
		err := errors.New("unsupported OS")
		return "", err
	}
}

// doesConfigFileExists checks if the config file exists.
func doesConfigFileExists() (bool, error) {
	filePath, err := getConfigFilePathBasedOnOS()
	if err != nil {
		return false, err
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false, nil
	}

	return true, nil
}

// createDefaultConfigFiles creates the config file.
func createDefaultConfigFiles() error {
	filePath, err := getConfigFilePathBasedOnOS()
	if err != nil {
		return err
	}

	// Create the directory.
	dirPath := filePath[:len(filePath)-len("/config.json")]
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return err
	}

	// Create the file inside the directory.
	file, err := os.Create(filePath)
	defer file.Close()

	defaultConfig := NewDefaultUserConfig()
	
	defaultConfigJSON, err := json.Marshal(defaultConfig)
	if err != nil {
		return err
	}
	
	_, err = file.Write([]byte(defaultConfigJSON))
	if err != nil {
		return err
	}

	return nil
}
