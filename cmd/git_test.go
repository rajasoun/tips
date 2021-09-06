// Licensed under the Creative Commons License.

package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func Test_Cases(t *testing.T) {
	t.Run("empty with ", func(t *testing.T) {
		outputBuffer := bytes.NewBufferString("")
		rootCmd.SetOut(outputBuffer)
		rootCmd.SetArgs([]string{"git", ""})
		err := gitCmd.Execute()
		fmt.Print(gitCmd.Execute())
		assert.Error(t, err)
	})
}
