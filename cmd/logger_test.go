// Licensed under the Creative Commons License.

package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
				want := tt.want
				assert.Equal(t, got, want)
			}
		})
	}
}
