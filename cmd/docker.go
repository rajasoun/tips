// Licensed under the Creative Commons License.

package cmd

import "github.com/spf13/cobra"

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
"List all containers : docker ps -a "`,
		Args: cobra.MaximumNArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			toolName = "docker"
			err := isValidTopic(args, toolName, cmd)
			return err
		},
	}
	return dockercmd
}
