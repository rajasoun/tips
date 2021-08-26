// Licensed under the Creative Commons License.

package cmd

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.rajasoun/tips/controller"

	"gopkg.in/yaml.v2"
)

var (
	gitCmd                   = GitCommand()
	rootCmd                  = NewRootCmd()
	dockerCmd                = DockerCommand()
	cmd                      *cobra.Command
	debug, cfgFile, toolName string

	createFile = os.Create
	copyData   = io.Copy
	getRequest = http.Get
	path       = os.Getenv("HOME")
	fileName   = "/.tips.yml"
)

const (
	validLen    int    = 1
	validArg    string = "git"
	emptyString string = ""
	firstLetter string = "g"
	dataLink    string = "https://raw.githubusercontent.com/rajasoun/tips/main/data/tips.json"
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

// git command functionality
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

// Execute executes the root command.
func Execute(writer io.Writer) error {
	rootCmd.SetOutput(writer)
	return rootCmd.Execute()
}

//  Checking argument and pass input to controller
func isValidTopic(args []string, toolName string, cmd *cobra.Command) error {
	switch {
	case len(args) == 0 && debug == emptyString:
		err := cmd.Help()
		return err
	case len(args) == 0 && debug != emptyString:
		return errors.New("please add an argument to debug")
	case args[0] != emptyString || debug != emptyString:
		if debug != emptyString {
			err := setUpLogs(cmd.OutOrStdout(), debug)
			if err != nil {
				return err
			}
		}
		logrus.WithField("loglevel", debug).Debug("successfully set logger level to debug ")
		input, err := getTopic(args)
		if err != nil {
			logrus.WithField("err", err).Debug("invalid user input")
			return err
		}
		input = toolName + " " + input
		logrus.WithField("userInput", input).Debug("successfully getting valid input ")
		controller.GetTipForTopic(input, cmd.OutOrStdout())
	}

	return nil
}

// getting topic with checking validation
func getTopic(args []string) (string, error) {
	var userInput = args[0]
	if isValidInput(userInput) {
		logrus.WithField("topic", userInput).Debug("successfully validation checked")
		return userInput, nil
	}
	return "", errors.New("argument should be greater than 1")
}

//  checking  userinput validation
func isValidInput(userInput string) bool {
	if len(userInput) > validLen && len(userInput) != 0 || userInput == validArg {
		return true
	}
	return false
}

// checking valid arguments
func isValidArguments(writer io.Writer, args []string) error {
	if args[0] != validArg {
		if string(args[0][0]) == firstLetter {
			fmt.Fprint(writer, "Did you mean this? \n git\n\n ")
		}
	}
	fmt.Fprint(writer, "unknown command ", args[0], " for tips \n")
	logrus.WithField("command", args[0]).Debug("unknown command for tips ")
	return errors.New("invalid command for tips")
}

// setting log level
func setUpLogs(out io.Writer, level string) error {
	logrus.SetOutput(out)
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	logrus.SetLevel(logLevel)
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	return nil
}

func init() {
	_ = InitializeTipsTool(fileName)
	cmd.PersistentFlags().StringVarP(&debug, "debug", "", "", "verbose logging")
	_ = cmd.PersistentFlags().MarkHidden("debug")
	rootCmd.AddCommand(gitCmd)
	rootCmd.AddCommand(dockerCmd)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tips.yaml)")
}

type confYml struct {
	Dir string `yaml:"tipsDataPath"`
}

func isExist(fileDir string) bool {
	if fileInfo, err := os.Stat(fileDir); err != nil {
		if os.IsNotExist(err) && fileInfo == nil {
			return false
		}
	}
	return true
}
func checkTipsData(fileName string) bool {
	pathv := path + fileName
	return isExist(pathv)
}

func createDir(dirPath string) error {
	info, err := createFile(dirPath) // create tips.yml
	if err != nil {
		fmt.Println("File info", info)
		return errors.New("issue on creating file")
	}
	configData := map[string]string{
		"tipsDataPath": "/.tips/tips.json",
	}
	if cfgFile != "" {
		configData["tipsDataPath"] = cfgFile
	}
	fmt.Println(configData["tipsDataPath"])

	data, _ := yaml.Marshal(&configData)
	_ = ioutil.WriteFile(dirPath, data, 0)
	return nil
}

// download url json file and save in filepath
func downloadFileFromURL(url string, filepath string) error {
	out, err := createFile(filepath)
	if err != nil {
		logrus.WithField("err", err).Debug("issue in creating file")
		return err
	}
	defer out.Close()
	resp, err := getRequest(url)
	if err != nil {
		logrus.WithField("err", err).Debug("getting err on get http request")
		return err
	}
	defer resp.Body.Close()
	_, err = copyData(out, resp.Body)
	if err != nil {
		logrus.WithField("err", err).Debug("getting error on copy data in dir")
		return err
	}
	return nil
}

// reading yml file data
func readfromYMLConfig(filePath string) (confYml, error) {
	config := confYml{}
	yamlFile, err := ioutil.ReadFile(path + filePath)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(yamlFile, &config)
	return config, err
}

func InitializeTipsTool(fileName string) error {
	// to do more condition
	// if .tips.yml not exist -- create it and download json
	// if exist ---check .tips dir is exist or not
	// if not create .tips dir and download json
	// if all are exist then update the .tips.yml and json file
	if !checkTipsData(fileName) {
		err := createDir(path + fileName)
		if err != nil {
			return err
		}
	}
	// to do add more condition when user add --config flag
	_ = os.Mkdir(path+"/.tips", 0700)
	pathValue, _ := readfromYMLConfig(fileName)
	err := downloadFileFromURL(dataLink, path+pathValue.Dir)
	return err
}
