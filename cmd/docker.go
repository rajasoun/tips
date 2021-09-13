// Licensed under the Creative Commons License.

package cmd

import (
	"github.com/spf13/cobra"
	"github.rajasoun/tips/controller"
)

var (
	dockerCmd = DockerCommand()
)

//  docker command functionality
func DockerCommand() *cobra.Command {
	var dockercmd = &cobra.Command{
		Use:   "docker",
		Short: "Docker provides the ability to package and run an application.",
		Long: ` "Docker is a software platform that simplifies the process of building, running,
managing and distributing applications."`,
		Aliases: []string{},

		Version: "0.1v",
		Example: `tips docker <command>
tips docker ps
"docker ps -a    :    LIST ALL CONTAINERS"`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = checklogger()
			toolName = "docker"
			input, err := isValidedArguments(args, toolName, cmd)
			if err == nil && input != "" {
				controller.GetTipForTopic(input, cmd.OutOrStdout())
			}
			return err
		},
	}
	return dockercmd
}
