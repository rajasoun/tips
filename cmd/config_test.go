// Licensed under the Creative Commons License.

package cmd

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileConfigurationSetting(t *testing.T) {
	_, err := os.Create("test.yml")
	if err != nil {
		t.Fatal(err)
	}
	testsIsExit := []struct {
		testDetails string
		input       string
		want        bool
	}{
		{testDetails: "check if file not exist", input: "dummy.yml", want: false},
		{testDetails: "check if file exist", input: "test.yml", want: true},
	}
	for _, tt := range testsIsExit {
		t.Run(tt.testDetails, func(t *testing.T) {
			got := isExist(tt.input)
			assert.Equal(t, got, tt.want)
		})
	}
	path = ""
	testsReadfromYMLConfig := []struct {
		testDetails string
		input       string
		want        []byte
	}{
		{testDetails: "check getting error on reading data from yml file", input: "/dummy/.json", want: nil},
		{testDetails: "check read config from yml file", input: "test.yml", want: nil},
	}
	for _, tt := range testsReadfromYMLConfig {
		t.Run(tt.testDetails, func(t *testing.T) {
			got, err := readfromYMLConfig(tt.input)
			if err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Nil(t, err)
				assert.NotEqual(t, tt.want, got)
			}
		})
	}
	testDownloadtoFile := []struct {
		testDetail    string
		link          string
		filepathInput string
		output        []byte
	}{
		{testDetail: "check if file is exist and link is correct ", link: "https://raw.githubusercontent.com/rajasoun/tips/main/data/tips.yml", filepathInput: "test.yml", output: nil},
		{testDetail: "check getting error when link is invalid", link: "", filepathInput: "", output: nil},
	}
	for _, tt := range testDownloadtoFile {
		t.Run(tt.testDetail, func(t *testing.T) {
			got := downloadToFile(tt.link, tt.filepathInput)
			if got != nil {
				assert.Error(t, got)
			} else {
				assert.Nil(t, got)
			}
		})
	}
	testsGetPath := []struct {
		testdetail string
		input      string
		want       string
	}{
		{testdetail: "check if Home dir is present", input: "$HOME/.tips/data.json", want: "/home/vscode/.tips/data.json"},
		{testdetail: "check if Home dir is not present", input: "/.dummy/data.txt", want: "/.dummy/data.txt"},
	}

	for _, tt := range testsGetPath {
		t.Run(tt.testdetail, func(t *testing.T) {
			got := getPath(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
	_, err = os.Create("test.txt")
	if err != nil {
		t.Fatal(err)
	}
	testsWriteFile := []struct {
		testDetails  string
		testfilePath string
		input        string
	}{
		{testDetails: "check write data into file", testfilePath: "test.txt", input: "hello"},
		{testDetails: "check getting an error on writing data into file", testfilePath: "", input: "hello"},
	}
	for _, tt := range testsWriteFile {
		t.Run(tt.testDetails, func(t *testing.T) {
			err := writeFile(tt.testfilePath, []byte(tt.input))
			if err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
	os.Remove("test.yml")
	os.Remove("test.txt")
}

func Test_tipsConfigurationSetting(t *testing.T) {
	t.Run("Checking data is loading or not from yml file", func(t *testing.T) {
		path = os.Getenv("HOME")
		if isExist(path + dir) {
			os.RemoveAll(path + dir)
		}
		got := TipsConfigurationSetting("/.tips.yml")
		assert.NoError(t, got)
	})
	t.Run("getting error when creating dir", func(t *testing.T) {
		path = ""
		dir = "/.dummy"
		got := TipsConfigurationSetting("/.tips.txt")
		assert.Error(t, got)
	})
	t.Run("getting error when downloading file from url", func(t *testing.T) {
		if isExist(os.Getenv("HOME") + "/.testdummy") {
			os.Remove(os.Getenv("HOME") + "/.testdummy")
		}
		configLink = "dummy.dummy.com"
		path = os.Getenv("HOME")
		dir = "/.testdummy"
		got := TipsConfigurationSetting("/.tips.txt")
		assert.Error(t, got)
	})
	t.Run("getting error when reading the file", func(t *testing.T) {
		if isExist("dummy") {
			os.RemoveAll("dummy")
		}
		configLink = "https://raw.githubusercontent.com/rajasoun/tips/main/data/tips.yml"
		path = ""
		dir = "dummy"
		got := TipsConfigurationSetting("")
		assert.Error(t, got)
	})
	t.Run("checking yml / json file exist or not,if .tips dir is exist", func(t *testing.T) {
		if !isExist(os.Getenv("HOME") + "/.testdir") {
			err := os.Mkdir(os.Getenv("HOME")+"/.testdir", 0700)
			if err != nil {
				t.Fatal(err)
			}
		}
		path = os.Getenv("HOME")
		dir = "/.testdir"
		got := TipsConfigurationSetting("/.tips.yml")
		assert.NoError(t, got)
		os.RemoveAll(os.Getenv("HOME") + "/" + ".testdir")
	})
}

func Test_httpDownload(t *testing.T) {
	tests := []struct {
		testdetail string
		input      string
		want       string
	}{
		{testdetail: "check get http requrest successfully", input: "https://raw.githubusercontent.com/rajasoun/tips/main/data/tips.yml", want: ""},
		{testdetail: "check error in http get request", input: "hello/hello", want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.testdetail, func(t *testing.T) {
			got, err := httPDownload(tt.input)
			if err != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NotNil(t, got)
				assert.NoError(t, err)
			}
		})
	}
	t.Run("check error on read the data from uri", func(t *testing.T) {
		test := func(io.Reader) ([]byte, error) {
			return nil, errors.New("getting simulation error")
		}
		readResponseBody = test
		got, err := httPDownload(configLink)
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}
