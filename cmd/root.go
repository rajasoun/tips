// Licensed under the Creative Commons License.

package cmd

import (
	"errors"
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	// "github.rajasoun/tips/config"

	"github.rajasoun/tips/controller"
)

var (
	gitCmd          = GitCommand()
	rootCmd         = NewRootCmd()
	dockerCmd       = DockerCommand()
	cmd             *cobra.Command
	debug, toolName string
	// cfgFile
	fileName = "/.tips.yml"
)

const (
	validLen    int    = 1
	validArg    string = "git"
	emptyString string = ""
	firstLetter string = "g"
)

func NewRootCmd() *cobra.Command {
	cmd = &cobra.Command{
		Use:     "tips",
		Long:    "tips provides help for docker and git cli commands ",
		Short:   "tips for command line interface function",
		Aliases: []string{},
		Version: "0.1v",
		Args:    cobra.MaximumNArgs(1),
		Example: `-> tips <tool_name> <command>

tips git push
tips docker ps`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			_ = TipsConfigurationSetting(fileName)
			if len(args) == 0 {
				err = cmd.Help()
			} else if err := isValidArguments(cmd.OutOrStdout(), args); err != nil {
				return err
			}
			return err
		},
	}
	return cmd
}

// Execute executes the root command.
func Execute(writer io.Writer) error {
	rootCmd.SetOutput(writer)
	return rootCmd.Execute()
}

//  put into another file , file name like validationUtil.go
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

// checking valid arguments
func isValidArguments(writer io.Writer, args []string) error {
	if args[0] != validArg {
		if string(args[0][0]) == firstLetter {
			fmt.Fprint(writer, "Did you mean this? \n git\n\n ")
		}
	}
	fmt.Fprint(writer, "unknown command ", args[0], " for tips \n")
	logrus.WithField("command", args[0]).Debug("unknown command for tips ")
	return errors.New("invalid command for tips")
}

func init() {
	cmd.PersistentFlags().StringVarP(&debug, "debug", "", "", "verbose logging")
	_ = cmd.PersistentFlags().MarkHidden("debug")
	rootCmd.AddCommand(gitCmd)
	rootCmd.AddCommand(dockerCmd)
	// cmd.PersistentFlags().StringVarP(&cfgFile, "cfgFile", "", "", "config file (default is $HOME/.tips.yaml or $HOME/.tips/tips.json)")
}
