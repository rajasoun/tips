// Licensed under the Creative Commons License.

package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LinuxCommand(t *testing.T) {
	t.Run("check get linux valid cmd", func(t *testing.T) {
		rootCmd.SetArgs([]string{"linux", "move"})
		writer := &bytes.Buffer{}
		err := Execute(writer)
		if err != nil {
			assert.Error(t, err)
		} else {
			gotWriter := writer.String()
			assert.Contains(t, gotWriter, "move")
		}
	})
}
