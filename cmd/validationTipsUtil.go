// Licensed under the Creative Commons License.

package cmd

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"unicode"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.rajasoun/tips/controller"
)

const (
	validLen    int    = 1
	emptyString string = ""
)

var (
	validArgDocker = "docker"
	dockerSuggest  = "d"
	gitSuggest     = "g"
	validArgGit    = "git"
)

//  Checking argument and pass input to controller
func isValidTopic(args []string, toolName string, cmd *cobra.Command) error {
	switch {
	case len(args) == 0 && debug == emptyString:
		err := cmd.Help()
		return err
	case len(args) == 0 && debug != emptyString:
		return errors.New("please add an argument to debug")
	case args[0] != emptyString || debug != emptyString || args[0] == "":
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
	} else if !isAlphabeticChar(userInput) {
		return "", errors.New("does not be contains special char/digit")
	}
	return "", errors.New("argument should not be empty/should be greater than 1")
}

//  checking  userinput validation
func isValidInput(userInput string) bool {
	if len(userInput) > validLen && isAlphabeticChar(userInput) {
		return true
	}
	return false
}

// checking valid arguments
func isValidArguments(writer io.Writer, args []string) error {
	if args[0] != validArgGit || args[0] != validArgDocker {
		if string(args[0][0]) == gitSuggest {
			fmt.Fprint(writer, "Did you mean this? \n git\n\n ")
		} else if string(args[0][0]) == dockerSuggest {
			fmt.Fprint(writer, "Did you mean this? \n docker\n\n ")
		}
	}
	fmt.Fprint(writer, "unknown command ", args[0], " for tips \n")
	logrus.WithField("command", args[0]).Debug("unknown command for tips ")
	return errors.New("invalid command for tips")
}

func isAlphabeticChar(input string) bool {
	isAlpha := regexp.MustCompile(`^[A-Za-z]+$`).MatchString
	return isAlpha(input) && !hasSymbol(input)
}
func hasSymbol(input string) bool {
	for _, letter := range input {
		if unicode.IsSymbol(letter) {
			return true
		}
	}
	return false
}
