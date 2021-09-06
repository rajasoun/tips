// Licensed under the Creative Commons License.

package cmd

import (
	"github.com/spf13/cobra"
)

var (
	gitCmd = GitCommand()
)

//  git command functionality

func GitCommand() *cobra.Command {
	var gitcmd = &cobra.Command{
		Use:   "git",
		Short: "Git is a DevOps tool used for source code management.",
		Long: ` "Git is used to tracking changes in the source code,
 enabling multiple developers to work together on non-linear development"`,
		Aliases: []string{},
		Version: "0.1v",
		Example: `tips git <command>

tips git stash
"Saving current state of unstaged changes to tracked files : git stash -k" `,
		Args: cobra.MaximumNArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			toolName = "git"
			err := isValidTopic(args, toolName, cmd)
			return err
		},
	}
	return gitcmd
}
