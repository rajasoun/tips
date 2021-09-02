// Licensed under the Creative Commons License.

package cmd

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationTipsUtil(t *testing.T) {
	t.Run("check valid data", func(t *testing.T) {
		input := []string{"check"}
		got, err := getTopic(input)
		if err != nil {
			assert.Error(t, err)
		} else {
			assert.Equal(t, input[0], got)
			assert.Nil(t, err)
		}
	})
	t.Run("check related commands and suggest", func(t *testing.T) {
		output := bytes.Buffer{}
		err := isValidArguments(&output, []string{"g"})
		got := output.String()
		want := "Did you mean this?"
		assert.Contains(t, got, want)
		assert.Error(t, err)
	})
	t.Run("check related commands and suggest", func(t *testing.T) {
		output := bytes.Buffer{}
		err := isValidArguments(&output, []string{"d"})
		got := output.String()
		want := "Did you mean this?"
		assert.Contains(t, got, want)
		assert.Error(t, err)
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
}
