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
func isValidedArguments(args []string, toolName string, cmd *cobra.Command) (string, error) {
	var err error
	var input string
	if len(args) == 0 {
		err = cmd.Help()
		return input, err
	}
	input, err = getValidTopic(args)
	if err != nil {
		logrus.WithField("err", err).Debug("invalid user input")
		return input, err
	}
	input = toolName + " " + input
	logrus.WithField("userInput", input).Debug("successfully getting valid input")
	return input, nil
}

// getting topic with checking validation
func getValidTopic(args []string) (string, error) {
	var userInput = args[0]
	if checkInputLen(userInput) {
		logrus.WithField("topic", userInput).Debug("successfully validation checked")
		return userInput, nil
	} else if !isAlphabeticChar(userInput) {
		return emptyString, errors.New("does not be contains special char/digit")
	}
	return emptyString, errors.New("argument should not be empty/should be greater than 1")
}

//  checking  userinput validation
func checkInputLen(userInput string) bool {
	if len(userInput) > validLen && isAlphabeticChar(userInput) {
		return true
	}
	return false
}

// checking valid arguments
func suggestedArgument(writer io.Writer, args []string) error {
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

// checking input
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
