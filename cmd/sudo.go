// Licensed under the Creative Commons License.

package cmd

import (
	"github.com/spf13/cobra"
	"github.rajasoun/tips/controller"
)

var (
	sudoCmd = SudoCommand()
)

//  sudo command functionality
func SudoCommand() *cobra.Command {
	var sudocmd = &cobra.Command{
		Use:     "sudo",
		Short:   "Sudo stands for SuperUser DO and is used to access restricted files and operations",
		Long:    `The sudo command temporarily elevates privileges allowing users to complete sensitive tasks without logging in as the root user.`,
		Aliases: []string{},
		Version: "0.1v",
		Example: `tips sudo <command>
tips sudo reboot`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = checklogger()
			toolName = "sudo"
			input, err := isValidedArguments(args, toolName, cmd)
			if err == nil && input != "" {
				controller.GetTipForTopic(input, cmd.OutOrStdout())
			}
			return err
		},
	}
	return sudocmd
}
