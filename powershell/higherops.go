package powershell

import (
	"strings"
)

// ExecStr - executes a commandline in powershell
func (runspace Runspace) ExecStr(commandStr string, useLocalScope bool) InvokeResults {
	command := runspace.CreateCommand()
	defer command.Delete()

	if strings.HasSuffix(commandStr, ".ps1") {
		command.AddCommand(commandStr, useLocalScope)
	} else {
		command.AddScript(commandStr, useLocalScope)
	}
	
	return command.Invoke()
}
