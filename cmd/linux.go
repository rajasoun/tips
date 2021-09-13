// Licensed under the Creative Commons License.

package cmd

import (
	"github.com/spf13/cobra"
	"github.rajasoun/tips/controller"
)

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
		Example: `tips linux <command>
tips linux copy
"wget [Option] [URL]    :    DOWNLOAD APPLICATIONS/WEB PAGES DIRECTLY FROM THE WEB"`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = checklogger()
			toolName = "linux"
			input, err := isValidedArguments(args, toolName, cmd)
			if err == nil && input != "" {
				controller.GetTipForTopic(input, cmd.OutOrStdout())
			}
			return err
		},
	}
	return linuxcmd
}
