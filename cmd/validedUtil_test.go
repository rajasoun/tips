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
		got, err := getValidTopic(input)
		if err != nil {
			assert.Error(t, err)
		} else {
			assert.Equal(t, input[0], got)
			assert.Nil(t, err)
		}
	})
	t.Run("check related commands and suggest", func(t *testing.T) {
		output := bytes.Buffer{}
		err := suggestedArgument(&output, []string{"g"})
		got := output.String()
		want := "Did you mean this?"
		assert.Contains(t, got, want)
		assert.Error(t, err)
	})
	t.Run("check related commands and suggest", func(t *testing.T) {
		output := bytes.Buffer{}
		err := suggestedArgument(&output, []string{"d"})
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
func TestCases(t *testing.T) {
	tests := []struct {
		testDetails string
		input       string
		want        bool
	}{
		{testDetails: "check input having digit values", input: "123a", want: false},
		{testDetails: "check input having special letter", input: "-!$$$$", want: false},
		{testDetails: "check input having only alphabets", input: "Abc", want: true},
		{testDetails: "check input having string with more than one words", input: "copy the dir", want: true},
		{testDetails: "check input having string with more than one words with unneed spaces", input: " copy the dir  ", want: false},
		{testDetails: "check input having string with more than one words", input: "copy the $8dir", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.testDetails, func(t *testing.T) {
			got := isAlphabeticChar(tt.input)
			assert.Equal(t, got, tt.want)
		})
	}
	t.Run("", func(t *testing.T) {
		got1 := hasSymbol("$$$$")
		assert.Equal(t, got1, true)
	})
	t.Run("", func(t *testing.T) {
		got1, err := getValidTopic([]string{"$$$$"})
		assert.Error(t, err)
		assert.Equal(t, got1, "")
	})
}
