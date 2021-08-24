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
}
type Tools struct {
	Git    []Tips `json:"git"`
	Docker []Tips `json:"docker"`
}

type confYml struct {
	Dir string `yaml:"tipsDataPath"`
}

const (
	defaultValue = "invalid command ,please pass valid tool command "
	emptyString  = " "
)

var (
	path     = os.Getenv("HOME")
	fileName = "/.tips.yml"
	fileRead = os.ReadFile
)

// GetTip returning Tip/Command to the controller
func GetTip(title string) string {
	data, _ := loadTipsFromJSON()
	commands := getAllCommands(data, title)
	for _, tip := range commands {
		return tip
	}
	return defaultValue
}

// getting all tips and titles
func getAllCommands(data Tools, title string) []string {
	title += emptyString
	cmdTool := strings.Split(title, emptyString)
	commands := make([]string, 0)
	if cmdTool[0] == "git" {
		for _, value := range data.Git {
			if strings.Contains(value.Tip, cmdTool[1]) || strings.Contains(value.Title, cmdTool[1]) {
				command := value.Title + " : " + value.Tip
				commands = append(commands, command)
			}
		}
	} else if cmdTool[0] == "docker" {
		for _, value := range data.Docker {
			if strings.Contains(value.Tip, cmdTool[1]) {
				command := value.Title + " : " + value.Tip
				commands = append(commands, command)
			}
		}
	}
	return commands
}

func loadTipsFromJSON() (Tools, error) {
	// run an app from main.go -> file path should be "data/tips.json" from developer side
	// if want to check all unit test cases ->file path should be "../data/tips.json"
	var path = getJSONFilePath()
	var data []byte
	data, _ = readJSONFile(path)
	var result Tools
	err := json.Unmarshal(data, &result)
	return result, err
}

// for main -- path should be $home+read from .tips.yml file from user side
// for testing -- ..//data/tips.json

func readfromYMLConfig(fileName string) (string, error) {
	config := confYml{}
	yamlFile, err := ioutil.ReadFile(path + fileName)
	if err != nil {
		return "", err
	}
	err = yaml.Unmarshal(yamlFile, &config)
	dirPath := path + config.Dir
	return dirPath, err
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
	// return currentDir + "/data/tips.json" // file path
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
