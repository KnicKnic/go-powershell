package powershell

import (
	"strings"
)

// ExecStr - executes a commandline in powershell
func (runspace Runspace) ExecStr(commandStr string, useLocalScope bool) {
	command := runspace.CreatePowershellCommand()
	defer command.Delete()

	// fields, ok := shell.Split(commandStr)
	// if !ok {
	// 	panic("command was invalid {" + commandStr + "}")
	// }

	if strings.HasSuffix(commandStr, ".ps1") {
		command.AddCommand(commandStr, useLocalScope)
	} else {
		command.AddScript(commandStr, useLocalScope)
	}
	// for i := 1; i < len(fields); i++ {
	// 	command.AddArgument(fields[i])
	// }
	command.Invoke()
}
