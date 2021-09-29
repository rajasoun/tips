// Licensed under the Creative Commons License.

// +build !integration

package model

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func init() {
	os.Setenv("GO_ENV", "test")
}

func TestGetTip(t *testing.T) {
	inputOuputData := []struct {
		name  string
		input string
		want  string
	}{
		{name: "Get Tip for valid Topic - rebase", input: "git,rebase", want: "git rebase master feature && git checkout master && git merge -    :    REBASES 'FEATURE' TO 'MASTER' AND MERGES IT IN TO MASTER "},
		{name: "Get Tip for invalid Topic - dummy", input: "dummy", want: "Tip is not available for this input,please pass valid input"},
		{name: "Get Tip for valid Topic - log", input: "docker,log", want: "docker log -S'<a term in the source>'    :    SEARCH CHANGE BY CONTENT"},
		{name: "Get Tip for valid Topic - move", input: "linux,move", want: "mv [Source] [Destination]    :    MOVE A FILE/DIRECTORY FROM ONE LOCATION TO ANOTHER."},
		{name: "Get Tip for valid Topic - state", input: "git,state", want: "git status    :    DISPLAYS THE STATE OF THE WORKING DIRECTORY"},
		{name: "Get Tip for valid Topic - poweroff", input: "sudo,poweroff", want: "sudo poweroff    :    POWEROFF DIRECTLY FROM YOUR TERMINAL"},
		{name: "Get Tip for valid Topic - install", input: "pip,install", want: "pip install <package>    :    TO INSTALL PACKAGE"},
	}
	for _, tt := range inputOuputData {
		t.Run(tt.name, func(t *testing.T) {
			got := GetTip(tt.input)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestLoadTipsFromJSON(t *testing.T) {
	t.Run("Load Tips From Json File and check if there are tips ", func(t *testing.T) {
		_, err := loadTipsFromJSON()
		assert.NoError(t, err)
	})
}

func TestGetTipJsonFilePath(t *testing.T) {
	t.Run("Check Getting Tips Json File Path Dynalically", func(t *testing.T) {
		got := getJSONFilePath()
		want := "/data/tips.json"
		assert.Contains(t, got, want)
	})
}

func TestGetCurrentWorkingDir(t *testing.T) {
	t.Run("checking current dir error", func(t *testing.T) {
		// Mocked function for os.Getwd
		myGetWd := func() (string, error) {
			myErr := errors.New("Simulated error")
			return "", myErr
		}
		// Update the var to this mocked function
		osGetWd = myGetWd
		// This will return error
		_, err := getCurrentWorkingDir()
		assert.Error(t, err)
	})
	t.Run("Checking Current Working directory path", func(t *testing.T) {
		// Mocked function for os.Getwd
		myGetWd := func() (string, error) {
			return "/gophers/workspace/tips", nil
		}
		osGetWd = myGetWd
		got, _ := getCurrentWorkingDir()
		want := "/gophers/workspace/tips"
		assert.Equal(t, got, want)
	})
}
func TestReadJsonFile(t *testing.T) {
	t.Run("Loading invalid Json File should fail", func(t *testing.T) {
		// Mocked function for os.ReadFile
		fileReading := func(string) ([]byte, error) {
			myErr := errors.New("Simulated error")
			return nil, myErr
		}
		fileRead = fileReading
		_, err := readJSONFile("/data")
		assert.Error(t, err)
	})
	t.Run("Unit Testing readjson file data", func(t *testing.T) {
		fileReading := func(string) ([]byte, error) {
			var data = []byte(`[{
				"title":"Rebases 'feature' to 'master' and merges it in to master ",
				"tip":"git rebase master feature && git checkout master && git merge -"
			 }]`)
			return data, nil
		}
		fileRead = fileReading
		got, _ := readJSONFile("/gophers/workspace//data/tips.json")
		want := "Rebases 'feature' to 'master'"
		assert.Contains(t, string(got), want)
	})
}

func TestReadfromYMLConfig(t *testing.T) {
	t.Run("Checking Error on not found file", func(t *testing.T) {
		_, err := readfromYMLConfig("/dummy/.json")
		assert.Error(t, err)
	})
	t.Run("checking the user file path ", func(t *testing.T) {
		err := mockcreatetestfile("testfile.yml", "$HOME/dummy/dummy.txt")
		if err != nil {
			t.Fatal(err)
		}
		path = ""
		got, err := readfromYMLConfig("testfile.yml")
		want := "dummy/dummy.txt"
		assert.Contains(t, got, want)
		assert.NoError(t, err)
		os.Remove("testfile.yml")
	})
	t.Run("checking the user file path ", func(t *testing.T) {
		err := mockcreatetestfile("testfile.yml", ".tips/tips.json")
		if err != nil {
			t.Fatal(err)
		}
		path = ""
		got, err := readfromYMLConfig("testfile.yml")
		want := ".tips/tips.json"
		assert.Contains(t, got, want)
		assert.NoError(t, err)
		os.Remove("testfile.yml")
	})
}

func mockcreatetestfile(testfile string, data string) error {
	_, err := os.Create(testfile)
	if err != nil {
		return err
	}
	filedata := map[string]string{
		"tipsDataLocalPath":  data,
		"tipsDataRemotePath": "",
	}
	dataa, _ := yaml.Marshal(&filedata)
	err = ioutil.WriteFile(testfile, dataa, 0)
	return err
}
