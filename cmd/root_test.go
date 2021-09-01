// Licensed under the Creative Commons License.

package cmd

import (
	"bytes"
	"fmt"
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
