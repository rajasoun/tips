// Licensed under the Creative Commons License.

package cmd

import "github.com/spf13/cobra"

var (
	linuxCmd = LinuxCommand()
)

//  linux command functionality

func LinuxCommand() *cobra.Command {
	var linuxcmd = &cobra.Command{
		Use:   "linux",
		Short: "Linux is an open source operating system (OS)",
		Long: ` "Linux facilitates with powerful support for networking.
The client-server systems can be easily set to a Linux system."`,
		Aliases: []string{},
		Version: "0.1v",
		Example: `tips linux <command>`,
		Args:    cobra.MaximumNArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			toolName = "linux"
			err := isValidTopic(args, toolName, cmd)
			return err
		},
	}
	return linuxcmd
}
