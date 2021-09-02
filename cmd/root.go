// Licensed under the Creative Commons License.

package cmd

import (
	"io"

	"github.com/spf13/cobra"
)

var (
	rootCmd         = NewRootCmd()
	cmd             *cobra.Command
	debug, toolName string
	// cfgFile
	fileName = "/.tips.yml"
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

func init() {
	cmd.PersistentFlags().StringVarP(&debug, "debug", "", "", "verbose logging")
	_ = cmd.PersistentFlags().MarkHidden("debug")
	rootCmd.AddCommand(gitCmd)
	rootCmd.AddCommand(dockerCmd)
	// cmd.PersistentFlags().StringVarP(&cfgFile, "cfgFile", "", "", "config file (default is $HOME/.tips.yaml or $HOME/.tips/tips.json)")
	_ = tipsConfigurationSetting(fileName)
}
