# Tips

## The Problem

As developers we need to remember several commands related to git, docker, tdd, general development environment setup making it overwhelming

## Solution

Command Line Tool to provide tips on the command to be used based on the topic

## Tips Tool
1. Tips is a command line tool
1. Tips tool interact with user and provide command related to user input
1. Tips helps provide the user on the commands that may pertain and brings up the  exact command  related  to the raw search.

## Usage

### Tips tool Usage
```
$ tips

  tips provides help for docker , linux and git cli commands

Usage:
  tips [flags]
  tips [command]

Examples:
-> tips <tool_name> <command/topic>

tips git push
tips docker ps
tips linux move

Available Commands:
  completion  generate the autocompletion script for the specified shell
  docker      Docker provides the ability to package and run an application.
  git         Git is a DevOps tool used for source code management.
  help        Help about any command
  linux       Linux is an open source operating system (OS)

Flags:
      --config string   config file (default is $HOME/.tips.yaml)
  -h, --help            help for tips
  -v, --version         version for tips

Use "tips [command] --help" for more information about a command.
```
## Features

1. Ease of Search.
1. Integration with multiple tools.

## Libraries

1. Cobra library  is used to build Tips command line app [cli].
2. Logrus library is used to set the log status (i.e debug).
3. We followed TDD design while building the Tips cli app, Also used Testify library for testing test cases.

##  Installation
```
curl -fsSL https://raw.githubusercontent.com/rajasoun/tips/main/install.sh | bash
export PATH="/opt/shellspec:/tools:$PATH"
```
## Future Scope :

1. User can add custom commands and also add commands data in tips tool
1. Develop feedback template where user can give their feedback.
