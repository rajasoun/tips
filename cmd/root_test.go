// Licensed under the Creative Commons License.

package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	os.Setenv("GO_ENV", "test")
}

func Test_NewRootCmd(t *testing.T) {
	t.Run("checking valid inputs", func(t *testing.T) {
		outputBuffer := bytes.NewBufferString("")
		rootCmd.SetOut(outputBuffer)
		rootCmd.SetArgs([]string{})
		err := rootCmd.Execute()
		if err != nil {
			t.Fatal(err)
		}
		out, err := ioutil.ReadAll(outputBuffer)
		if err != nil {
			t.Fatal(err)
		}
		got := string(out)
		want := "help"
		assert.Contains(t, got, want, "expected \"%s\" got \"%s\"", want, got)
	})
	t.Run("checking invalid user inputs", func(t *testing.T) {
		inputBuffer := "dummy"
		rootCmd.SetArgs([]string{inputBuffer})
		err := rootCmd.Execute()
		if err != nil {
			assert.Error(t, err)
		}
		assert.Error(t, err)
	})
	t.Run("checking config valid path for data file from user", func(t *testing.T) {
		outputBuffer := bytes.NewBufferString("")
		rootCmd.SetOut(outputBuffer)
		rootCmd.SetArgs([]string{"--cfgFile", "/dummy/./dummy.json"})
		err := rootCmd.Execute()
		if err != nil {
			assert.Error(t, err)
		}
		out, err := ioutil.ReadAll(outputBuffer)
		if err != nil {
			t.Fatal(err)
		}
		got := string(out)
		want := "help"
		assert.Contains(t, got, want, "expected \"%s\" got \"%s\"", want, got)
	})
}

func TestExecute(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
		flag  string
	}{
		{"Success Case for git command", "stash", "stash", "git"},
		{"Error Case", "help", "help", "--tips"},
		{"Invalid Data", "j", "help", "docker"},
		{"Success Case for docker command ", "log", "log", "docker"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootCmd.SetArgs([]string{tt.flag, tt.input})
			writer := &bytes.Buffer{}
			err := Execute(writer)
			if err != nil {
				assert.Error(t, err)
			} else {
				gotWriter := writer.String()
				fmt.Print(gotWriter)
				assert.Contains(t, gotWriter, tt.want)
			}
		})
	}
}

func Test_SetLogger(t *testing.T) {
	tests := []struct {
		name  string
		want  string
		level string
	}{
		{"Checking set level logger ", "", "debug"},
		{"invalid level logger", "not a valid logrus Level", "dummy"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := bytes.Buffer{}
			err := setUpLogs(&output, tt.level)
			if err != nil {
				assert.Error(t, err)
			} else {
				got := output.String()
				fmt.Print(got)
				want := tt.want
				assert.Equal(t, got, want)
			}
		})
	}
}
func Test_GitCommand(t *testing.T) {
	t.Run("checking help command", func(t *testing.T) {
		outputBuffer := bytes.NewBufferString("")
		rootCmd.SetOut(outputBuffer)
		rootCmd.SetArgs([]string{"git"})
		err := gitCmd.Execute()
		if err != nil {
			t.Fatal(err)
		}
		out, err := ioutil.ReadAll(outputBuffer)
		if err != nil {
			t.Fatal(err)
		}
		got := string(out)
		want := "help"
		assert.Contains(t, got, want, "want \"%s\" got \"%s\"", want, got)
	})

	t.Run("Checking valid data", func(t *testing.T) {
		outputBuffer := bytes.NewBufferString("")
		rootCmd.SetOut(outputBuffer)
		expected := "checkout"
		rootCmd.SetArgs([]string{"git", expected})
		err := gitCmd.Execute()
		if err != nil {
			t.Fatal(err)
		}
		out, err := ioutil.ReadAll(outputBuffer)
		if err != nil {
			t.Fatal(err)
		}
		got := string(out)
		assert.Contains(t, got, expected, "expected \"%s\" got \"%s\"", expected, got)
	})
	t.Run("Checking valid data", func(t *testing.T) {
		outputBuffer := bytes.NewBufferString("")
		rootCmd.SetOut(outputBuffer)
		expected := "d"
		rootCmd.SetArgs([]string{"git", expected})
		err := gitCmd.Execute()
		if err != nil {
			assert.Error(t, err)
		}
		out, err := ioutil.ReadAll(outputBuffer)
		if err != nil {
			assert.Error(t, err)
		}
		got := string(out)
		assert.Contains(t, got, "help", "expected \"%s\" got \"%s\"", "help", got)
	})
	t.Run("checking valid command", func(t *testing.T) {
		outputBuffer := bytes.NewBufferString("")
		rootCmd.SetOut(outputBuffer)
		rootCmd.SetArgs([]string{"gi"})
		err := gitCmd.Execute()
		assert.Error(t, err)
	})
	t.Run("checking valid command", func(t *testing.T) {
		outputBuffer := bytes.NewBufferString("")
		rootCmd.SetOut(outputBuffer)
		rootCmd.SetArgs([]string{"git", "--debug", "debug"})
		err := gitCmd.Execute()
		assert.Error(t, err)
	})
	t.Run("checking  valid logger status", func(t *testing.T) {
		outputBuffer := bytes.NewBufferString("")
		rootCmd.SetOut(outputBuffer)
		rootCmd.SetArgs([]string{"git", "push", "--debug", "dummy"})
		err := gitCmd.Execute()
		assert.Error(t, err)
	})
}
func TestDownloadFileFromURL(t *testing.T) {
	file := "dummytest.json"
	t.Run("Getting error when url is not correct", func(t *testing.T) {
		got := downloadFileFromURL("", file)
		assert.Error(t, got)
	})
	t.Run("getting error when saving file path is not correct", func(t *testing.T) {
		got := downloadFileFromURL("", "hello/dummy.json")
		assert.Error(t, got)
	})
	t.Run("checking get noerror on downloading the json file from url", func(t *testing.T) {
		got := downloadFileFromURL(dataLink, file)
		assert.NoError(t, got)
	})
	t.Run("creating issue on copy the json file in dir path", func(t *testing.T) {
		copyFunc := func(io.Writer, io.Reader) (int64, error) {
			var a int64
			return a, errors.New("error")
		}
		copyData = copyFunc
		got := downloadFileFromURL(dataLink, file)
		assert.Error(t, got)
	})
	os.Remove(file)
}

func Test_checkTipsData(t *testing.T) {
	t.Run("checking home dir  when file is present", func(t *testing.T) {
		file := "/dummy.txt"
		_, err := os.Create(path + file)
		if err != nil {
			t.Fatal(err)
		}
		got := checkTipsData(file)
		want := true
		assert.Equal(t, got, want)
		os.Remove(path + file)
	})
	t.Run("if file is not present", func(t *testing.T) {
		got := checkTipsData("xyz.txt")
		want := false
		assert.Equal(t, got, want)
	})
}

func Test_createDir(t *testing.T) {
	t.Run("creating new dir ", func(t *testing.T) {
		file := "xyz.yml"
		got := createDir(file)
		assert.NoError(t, got)
		os.Remove(file)
	})
	t.Run("If dir path is not exist to create a file", func(t *testing.T) {
		createfileError := func(string) (*os.File, error) {
			return nil, errors.New("simulation error")
		}
		createFile = createfileError
		got := createDir(".data.yml")
		assert.Error(t, got)
	})
}

func TestReadfromYMLConfig(t *testing.T) {
	t.Run("Reading data from yml file", func(t *testing.T) {
		_, err := readfromYMLConfig("/dummy/.json")
		assert.Error(t, err)
	})
	t.Run(" getting noerror on reading data from yml file", func(t *testing.T) {
		// refactor create a mock file
		got, err := readfromYMLConfig(fileName)
		assert.NoError(t, err)
		assert.Nil(t, err)
		assert.NotNil(t, got)
	})
}

func Test_initilizeTipsTool(t *testing.T) {
	t.Run("Checking data is loading or not from yml file", func(t *testing.T) {
		rootCmd.SetArgs([]string{})
		got := InitializeTipsTool("abctest.txt")
		assert.Error(t, got)
		os.Remove(path + "abctest.txt")
	})
	t.Run("if file is not present", func(t *testing.T) {
		createfileError := func(string) (*os.File, error) {
			return nil, errors.New("simulation error")
		}
		createFile = createfileError
		got := InitializeTipsTool("abcDummy")
		assert.Error(t, got)
		os.Remove(path + "abcDummy")
	})
}

func TestIsExit(t *testing.T) {
	t.Run("file is exist or not", func(t *testing.T) {
		got := isExist("dummy.yml")
		want := false
		assert.Equal(t, got, want)
	})
	mockremove()
}

func mockremove() {
	var err error
	err = os.Remove(path + "/.tips.yml")
	if err != nil {
		fmt.Print(err)
	}
	err = os.Remove(path + "/.tips")
	if err != nil {
		fmt.Print(err)
	}
}
