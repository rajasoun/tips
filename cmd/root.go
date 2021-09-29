// Licensed under the Creative Commons License.

package cmd

import (
	"io"

	"github.com/spf13/cobra"
)

var (
	rootCmd  = NewRootCmd()
	cmd      *cobra.Command
	toolName string
	// configPath
	fileName = "/.tips.yml"
	debug    bool
)

//  root tips cli functionality
func NewRootCmd() *cobra.Command {
	cmd = &cobra.Command{
		Use:     "tips",
		Long:    "Tips provides help for docker , git ,sudo , pip and linux cli commands ",
		Short:   "Tips for command line interface function",
		Aliases: []string{},
		Version: "0.1v",
		Args:    cobra.MaximumNArgs(1),
		Example: `SYNTAX: tips <tool_name> <command/topic>
E.g:
        "tips git saving"
        "tips docker ps"
        "tips linux move"
	"tips pip install"
	"tips sudo reboot"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			// checking logger status is set or not
			_ = checklogger()
			if len(args) == 0 {
				err = cmd.Help()
			} else if err := suggestedArgument(cmd.OutOrStdout(), args); err != nil {
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

func init() {
	cmd.PersistentFlags().BoolVar(&debug, "debug", false, "set log level to debug")
	_ = cmd.PersistentFlags().MarkHidden("debug")
	rootCmd.AddCommand(gitCmd)
	rootCmd.AddCommand(dockerCmd)
	rootCmd.AddCommand(linuxCmd)
	rootCmd.AddCommand(sudoCmd)
	rootCmd.AddCommand(pipCmd)

	// cmd.PersistentFlags().StringVarP(&configPath, "configPath", "", "", "config file (default is $HOME/.tips.yaml or $HOME/.tips/tips.json)")
	_ = tipsConfigurationSetting(fileName)
}
