// Licensed under the Creative Commons License.

package cmd

import (
	"errors"
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.rajasoun/tips/controller"
)

const (
	validLen    int    = 1
	validArg    string = "git"
	emptyString string = ""
	firstLetter string = "g"
)

//  Checking argument and pass input to controller
func isValidTopic(args []string, toolName string, cmd *cobra.Command) error {
	switch {
	case len(args) == 0 && debug == emptyString:
		err := cmd.Help()
		return err
	case len(args) == 0 && debug != emptyString:
		return errors.New("please add an argument to debug")
	case args[0] != emptyString || debug != emptyString:
		if debug != emptyString {
			err := setUpLogs(cmd.OutOrStdout(), debug)
			if err != nil {
				return err
			}
		}
		logrus.WithField("loglevel", debug).Debug("successfully set logger level to debug ")
		input, err := getTopic(args)
		if err != nil {
			logrus.WithField("err", err).Debug("invalid user input")
			return err
		}
		input = toolName + " " + input
		logrus.WithField("userInput", input).Debug("successfully getting valid input ")
		controller.GetTipForTopic(input, cmd.OutOrStdout())
	}

	return nil
}

// getting topic with checking validation
func getTopic(args []string) (string, error) {
	var userInput = args[0]
	if isValidInput(userInput) {
		logrus.WithField("topic", userInput).Debug("successfully validation checked")
		return userInput, nil
	}
	return "", errors.New("argument should be greater than 1")
}

//  checking  userinput validation
func isValidInput(userInput string) bool {
	if len(userInput) > validLen && len(userInput) != 0 || userInput == validArg {
		return true
	}
	return false
}

var validArgDocker = "docker"
var dockerSuggest = "d"

// checking valid arguments
func isValidArguments(writer io.Writer, args []string) error {
	if args[0] != validArg || args[0] != validArgDocker {
		if string(args[0][0]) == firstLetter {
			fmt.Fprint(writer, "Did you mean this? \n git\n\n ")
		} else if string(args[0][0]) == dockerSuggest {
			fmt.Fprint(writer, "Did you mean this? \n docker\n\n ")
		}
	}
	fmt.Fprint(writer, "unknown command ", args[0], " for tips \n")
	logrus.WithField("command", args[0]).Debug("unknown command for tips ")
	return errors.New("invalid command for tips")
}
