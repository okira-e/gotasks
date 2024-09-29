package domain

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"runtime"

	"github.com/okira-e/gotasks/internal/opt"
	"github.com/okira-e/gotasks/internal/utils"
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

// DoesUserConfigExist checks if a user config has already be generated for this user.
func DoesUserConfigExist() (bool, error) {
	filePath, err := GetConfigFilePathBasedOnOS()
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
	
	return config, nil
}

// GetUserConfig reads the user config file and returns a pointer 
// to a UserConfig object.
func GetUserConfig() (*UserConfig, error) {
	userConfig := new(UserConfig)

	filePath, err := GetConfigFilePathBasedOnOS()
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

	return userConfig, nil
}

func NewDefaultUserConfig() *UserConfig {
	ret := new(UserConfig)
	
	return ret
}

// AddBoard adds a new board to the config.
func (self *UserConfig) AddBoard(boardName string, dirPath string) error {
	board := new(Board)
	
	board.Name = boardName
	board.Dir = dirPath
	
	// Add the newly created board to the config.
	self.Boards = append(self.Boards, board)
	self.writeToDisk()
	
	return nil
}

// AddTask adds a new task to the left most column (idealy called Backlog).
func (self *UserConfig) AddTask(boardName string, task *Task) error {
	utils.SaveLog(utils.Debug, "Adding task", map[string]any{"task": task})
	
	boardOpt := self.GetBoard(boardName)
	if boardOpt.IsNone() {
		return errors.New("Couldn't find the board while trying to add a task")
	}

	board := boardOpt.Unwrap()
	
	if len(board.Columns) == 0 {
		return errors.New("No columns found to add this task to.")
	}
	
	columnName := board.Columns[0]
	
	board.Tasks[columnName] = append(board.Tasks[columnName], task)
	
	err := self.UpdateBoard(board)
	if err != nil {
		return err
	}
	
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
	filePath, err := GetConfigFilePathBasedOnOS()
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

type Board struct {
	Name    string   `json:"name"`
	Dir     string   `json:"dir"`
	Columns []string `json:"columns"`
	// Tasks are the individual cards on the board representing a task.
	Tasks map[string][]*Task `json:"tasks"`
}


// GetColumnForTask returns the name and the index of the column that this task belongs to. 
// It returns -1 as the index if didn't find the column.
func (self *Board) GetColumnForTask(task *Task) (string, int) {
	for i, columnName := range self.Columns {
		for _, it := range self.Tasks[columnName] {
			if it == task {
				return columnName, i
			}
		}
	}
	
	return "", -1
}

// GetConfigFilePathBasedOnOS returns the config folder path based on the OS.
func GetConfigDirPathBasedOnOS() (string, error) {
	var osUserName string

	if runtime.GOOS == "windows" {
		osUserName = os.Getenv("USERNAME")
		return "C:\\Users\\" + osUserName + "\\AppData\\Roaming\\gotasks", nil
	} else if runtime.GOOS == "darwin" {
		osUserName = os.Getenv("USER")
		return "/Users/" + osUserName + "/Library/Application Support/gotasks", nil
	} else if runtime.GOOS == "linux" {
		osHomeDir := os.Getenv("HOME")
		return osHomeDir + "/.config/gotasks", nil
	} else {
		err := errors.New("unsupported OS")
		return "", err
	}
}

// GetConfigFilePathBasedOnOS returns the config file path based on the OS.
func GetConfigFilePathBasedOnOS() (string, error) {
	configDirPath, err := GetConfigDirPathBasedOnOS()
	if err != nil {
		return "", err
	}
	
	if runtime.GOOS == "windows" {
		return configDirPath + "\\config.json", nil
	} else if runtime.GOOS == "darwin" {
		return configDirPath + "/config.json", nil
	} else if runtime.GOOS == "linux" {
		return configDirPath + "/config.json", nil
	} else {
		err := errors.New("unsupported OS")
		return "", err
	}
}

// doesConfigFileExists checks if the config file exists.
func doesConfigFileExists() (bool, error) {
	filePath, err := GetConfigFilePathBasedOnOS()
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
	filePath, err := GetConfigFilePathBasedOnOS()
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
