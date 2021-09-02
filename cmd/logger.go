// Licensed under the Creative Commons License.

package cmd

import (
	"io"

	"github.com/sirupsen/logrus"
)

// setting log level status for debugging
func setUpLogs(out io.Writer, level string) error {
	logrus.SetOutput(out)
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	logrus.SetLevel(logLevel)
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	return nil
}
