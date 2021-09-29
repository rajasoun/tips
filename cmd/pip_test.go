package cmd

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PipPythonCommand(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
		flag  string
	}{
		//	{"Error Case", "help", "help", "--sudo"},
		//	{"Invalid Data", "j", "help", "sudo"},
		{"Success Case for sudo command ", "install", "install", "pip"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootCmd.SetArgs([]string{tt.flag, tt.input})
			writer := &bytes.Buffer{}
			err := Execute(writer)
			fmt.Print(writer)
			if err != nil {
				assert.Error(t, err)
			}
			//else {
			gotWriter := writer.String()
			fmt.Print(gotWriter)
			assert.Contains(t, gotWriter, tt.want)
			//	}
		})
	}
}
