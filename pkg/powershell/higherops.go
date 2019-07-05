package powershell

import (
	"strings"
)

func processArgs(command Command, args ...interface{}) InvokeResults {
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			command.AddArgumentString(v)
		case Object:
			command.AddArgument(v)
		default:
			panic("unknown argument")
		}
	}
	return command.Invoke()
}

// ExecScript - executes a series of statements (not to be confused with .ps1 files which are commands) in powershell
func (runspace Runspace) ExecScript(commandStr string, useLocalScope bool, args ...interface{}) InvokeResults {
	command := runspace.CreateCommand()
	defer command.Delete()
	command.AddScript(commandStr, useLocalScope)
	return processArgs(command, args...)
}

// ExecCommand - executes a command (cmdlets, command files (.ps1), functions, ...) in powershell
func (runspace Runspace) ExecCommand(commandStr string, useLocalScope bool, args ...interface{}) InvokeResults {
	command := runspace.CreateCommand()
	defer command.Delete()
	command.AddCommand(commandStr, useLocalScope)
	return processArgs(command, args...)
}

// ExecStr - executes a commandline in powershell
func (runspace Runspace) ExecStr(commandStr string, useLocalScope bool, args ...interface{}) InvokeResults {
	command := runspace.CreateCommand()
	defer command.Delete()

	if strings.HasSuffix(commandStr, ".ps1") {
		command.AddCommand(commandStr, useLocalScope)
	} else {
		command.AddScript(commandStr, useLocalScope)
	}
	return processArgs(command, args...)
}
