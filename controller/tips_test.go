// Licensed under the Creative Commons License.

package controller

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	os.Setenv("GO_ENV", "test")
}

func TestGetTipForTopicIntegration(t *testing.T) {
	outputBuffer := bytes.Buffer{}
	inputOuputData := []struct {
		name  string
		input string
		want  string
	}{
		{name: "Checking with Valid input", input: "git,delete", want: "git branch -d <local_branchname>    :    DELETE LOCAL BRANCH"},
		{name: "Checking with invalid input", input: "hello", want: "not available"},
	}
	for _, tt := range inputOuputData {
		t.Run(tt.name, func(t *testing.T) {
			GetTipForTopic(tt.input, &outputBuffer)
			got := outputBuffer.String()
			assert.Contains(t, got, tt.want)
		})
	}
}
