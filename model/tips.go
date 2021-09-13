// Licensed under the Creative Commons License.

package model

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Tips struct {
	Title string `json:"title"`
	Tip   string `json:"tip"`
	// Alternatives []string `json:"alternatives"`
}
type Tools struct {
	Git    []Tips `json:"git"`
	Docker []Tips `json:"docker"`
	Linux  []Tips `json:"linux"`
}

type confYml struct {
	TipsDataLocalPath  string `yaml:"tipsDataLocalPath"`
	TipsDataRemotePath string `yaml:"tipsDataRemotePath"`
}

const (
	defaultValue = "Tip is not available for this input,please pass valid input"
	// emptyString   = " "
	specialLetter = ","
)

var (
	path     = os.Getenv("HOME") + dir
	fileName = "/.tips.yml"
	dir      = "/.tips"
	fileRead = os.ReadFile
)

// GetTip returning Tip/Command to the controller
func GetTip(title string) string {
	data, _ := loadTipsFromJSON()
	commands := getAllCommands(data, title)
	for _, tip := range commands {
		logrus.WithField("value", tip).Debug("getting valid tips from json, cmd exist")
		return tip
	}
	logrus.WithField("value", defaultValue).Debug("getting invalid tips from json, cmd not exist")
	return defaultValue
}

// getting all tips and titles
func getAllCommands(data Tools, title string) []string {
	title += specialLetter
	cmdTool := strings.Split(title, specialLetter)
	commands := make([]string, 0)
	switch {
	case cmdTool[0] == "git":
		commands = gettingToolcmd(data.Git, cmdTool[1])
	case cmdTool[0] == "docker":
		commands = gettingToolcmd(data.Docker, cmdTool[1])
	case cmdTool[0] == "linux":
		commands = gettingToolcmd(data.Linux, cmdTool[1])
	}
	return commands
}
func gettingToolcmd(tool []Tips, input string) []string {
	commands := make([]string, 0)
	for _, value := range tool {
		if strings.Contains(strings.ToLower(value.Tip), input) || strings.Contains(strings.ToLower(value.Title), input) {
			command := value.Tip + "    :    " + strings.ToUpper(value.Title)
			commands = append(commands, command)
		}
	}
	return commands
}
func loadTipsFromJSON() (Tools, error) {
	var path = getJSONFilePath()
	var data []byte
	data, _ = readJSONFile(path)
	var result Tools
	err := json.Unmarshal(data, &result)
	logrus.WithField("error", err).Debug("loading the data from json")
	return result, err
}

// Run from main -- path should be $home+/.tips/.tips.yml file from user side
// and  --  jsonfile path should be $home+/.tips/data.json file from user side
// Run for testing -- ..//data/tips.json

func readfromYMLConfig(fileName string) (string, error) {
	config := confYml{}
	yamlFile, err := ioutil.ReadFile(path + fileName)
	if err != nil {
		logrus.WithField("error", err).Debug("facing issue on reading .tips.yml file")
		return "", err
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if strings.Contains(config.TipsDataLocalPath, "$HOME") {
		value := strings.ReplaceAll(config.TipsDataLocalPath, "$HOME", "")
		logrus.WithField("config data", config).Debug("successfully get json data file path")
		return os.Getenv("HOME") + value, err
	}
	logrus.WithField("config data", config).Debug("successfully get json data file path")
	return config.TipsDataLocalPath, err
}

// getting file path for main file and testing function
func getJSONFilePath() string {
	currentDir, _ := readfromYMLConfig(fileName)
	isInTest := os.Getenv("GO_ENV") == "test"
	if isInTest {
		currentDir, _ = getCurrentWorkingDir()
		baseDir := filepath.Base(currentDir)
		currentDir = strings.ReplaceAll(currentDir, baseDir, "")
		currentDir += "/data/tips.json"
	}
	return currentDir
}

// reading data from json file
func readJSONFile(path string) ([]byte, error) {
	data, err := fileRead(path)
	if err != nil {
		logrus.WithField("file path ", path).Debug("unsuccessfully reading the file path ")
		return nil, err
	}
	logrus.WithField("file path ", path).Debug("successfully reading the file path ")
	return data, nil
}

// Get Working directory function
var osGetWd = os.Getwd

// getting current working dir.
func getCurrentWorkingDir() (string, error) {
	workingDir, err := osGetWd()
	if err != nil {
		logrus.WithField("working dir", workingDir).Debug("unsuccessfully reading the working dir path ")
		return "", errors.New("could not get current working directory")
	}
	logrus.WithField("working dir", workingDir).Debug("successfully reading the working dir path ")
	return workingDir, nil
}
