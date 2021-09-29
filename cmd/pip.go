// Licensed under the Creative Commons License.

package cmd

import (
	"github.com/spf13/cobra"
	"github.rajasoun/tips/controller"
)

var (
	pipCmd = PipCommand()
)

//  sudo command functionality
func PipCommand() *cobra.Command {
	var pipCmd = &cobra.Command{
		Use:     "pip",
		Short:   "PIP is a package manager for Python packages",
		Long:    `Pip is a tool to create Python virtual environment`,
		Aliases: []string{},
		Version: "0.1v",
		Example: `tips pip <command>
tips pip install`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = checklogger()
			toolName = "pip"
			input, err := isValidedArguments(args, toolName, cmd)
			if err == nil && input != "" {
				controller.GetTipForTopic(input, cmd.OutOrStdout())
			}
			return err
		},
	}
	return pipCmd
}
