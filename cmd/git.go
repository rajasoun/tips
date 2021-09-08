// Licensed under the Creative Commons License.

package cmd

import (
	"github.com/spf13/cobra"
	"github.rajasoun/tips/controller"
)

var (
	gitCmd = GitCommand()
)

//  git command functionality

func GitCommand() *cobra.Command {
	var gitcmd = &cobra.Command{
		Use:   "git",
		Short: "Git is a DevOps tool used for source code management.",
		Long: `"Git is used to tracking changes in the source code,
enabling multiple developers to work together on non-linear development"`,
		Aliases: []string{},
		Version: "0.1v",
		Example: `tips git <command>
E.g:
"tips git stash"
"git stash -k    :    SAVING CURRENT STATE OF UNSTAGED CHANGES TO TRACKED FILES" `,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = checklogger()
			toolName = "git"
			input, err := isValidedArguments(args, toolName, cmd)
			if err == nil && input != "" {
				controller.GetTipForTopic(input, cmd.OutOrStdout())
			}
			return err
		},
	}
	return gitcmd
}
