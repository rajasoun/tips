// Licensed under the Creative Commons License.

package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var (
	configLink       = "https://raw.githubusercontent.com/rajasoun/tips/main/data/tips.yml"
	dir              = "/.tips"
	path             = os.Getenv("HOME")
	getHTTPRequest   = http.Get
	readResponseBody = io.ReadAll
)

type confYml struct {
	TipsDataLocalPath  string `yaml:"tipsDataLocalPath"`
	TipsDataRemotePath string `yaml:"tipsDataRemotePath"`
}

func isExist(name string) bool {
	if fileInfo, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) && fileInfo == nil {
			logrus.WithField("file", name).Debug("the file is not exist")
			return false
		}
	}
	logrus.WithField("file", name).Debug("the file is exist at location")
	return true
}

func httPDownload(url string) ([]byte, error) {
	res, err := getHTTPRequest(url)
	if err != nil {
		logrus.WithField("url ", url).Debug("not getting response from the url, Error found:", err)
		return nil, err
	}
	defer res.Body.Close()
	bytesData, err := readResponseBody(res.Body)
	if err != nil {
		logrus.WithField("Error", err).Debug("getting error on reading the data from url")
		return nil, err
	}
	logrus.WithField("error", err).Debug("successfully read the data and get response from url")
	return bytesData, err
}

func writeFile(dst string, data []byte) error {
	err := ioutil.WriteFile(dst, data, 0444) // only read permission
	if err != nil {
		logrus.WithField("found error", err).Debug("getting error not able to write data into file/ file not exist")
		return err
	}
	logrus.WithField("error", err).Debug("successfully write data into the file")
	return nil
}

func downloadToFile(url string, dst string) error {
	var err error
	var bytesData []byte
	if bytesData, err = httPDownload(url); err == nil {
		if writeFile(dst, bytesData) == nil {
			logrus.WithField("file path", dst).Debug("successfully download the file from url")
			return nil
		}
	}
	logrus.WithField("file path", dst).Debug("getting error on download the file from url: ", url)
	return err
}

// reading yml file data
func readfromYMLConfig(filePath string) (confYml, error) {
	config := confYml{}
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		logrus.WithField("error", err).Debug("getting error on read the yml file from local path ,filepath:", filePath)
		return config, err
	}
	err = yaml.Unmarshal(yamlFile, &config)
	logrus.WithField("file path ", filePath).Debug("successfully read the .tips.yml file from filepath : ", filePath)
	return config, err
}

// add setconfigfunc
func TipsConfigurationSetting(fileName string) error {
	var err error
	switch {
	case !isExist(path + dir):
		err = os.Mkdir(path+dir, 0700)
		if err != nil {
			logrus.WithField("error", err).Debug("facing issue on making directory at path:", path+dir)
			return err
		}
		err = downloadToFile(configLink, path+dir+fileName)
		if err != nil {
			logrus.WithField("error", err).Debug("facing issue on download file from url at path:", path+dir+fileName)
			return err
		}
		pathValue, err := readfromYMLConfig(path + dir + fileName)
		if err != nil {
			logrus.WithField("error", err).Debug("facing issue on download file from url at path:", path+dir+fileName)
			return err
		}
		localPath := getPath(pathValue.TipsDataLocalPath)
		_ = downloadToFile(pathValue.TipsDataRemotePath, localPath)
		fmt.Print("downloaded successfully")
		logrus.WithField("Path location", localPath).Debug("successfully download the file at location from", pathValue.TipsDataRemotePath)
	case isExist(path + dir):
		// if possible then convert into switch
		if !isExist(path + dir + fileName) {
			err = downloadToFile(configLink, path+dir+fileName)
			if !isExist(path + dir + "/" + "data.json") {
				pathValue, _ := readfromYMLConfig(path + dir + fileName)
				localPath := getPath(pathValue.TipsDataLocalPath)
				err = downloadToFile(pathValue.TipsDataRemotePath, localPath)
			}
			// one more condition for check data.json file
			logrus.WithField("error", err).Debug("downloaded the files if .tips dir is exist")
			return err
		}
	}
	return nil
}

func getPath(pathValue string) string {
	var homeVar = "$HOME"
	if strings.Contains(pathValue, homeVar) {
		homeDir := os.Getenv("HOME")
		pathValue = strings.ReplaceAll(pathValue, homeVar, "")
		logrus.WithField("file path ", path).Debug("successfully updated the file path")
		return homeDir + pathValue
	}
	return pathValue
}
