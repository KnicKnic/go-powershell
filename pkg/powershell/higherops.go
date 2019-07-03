package powershell

import (
	"strings"
)

// ExecStr - executes a commandline in powershell
func (runspace Runspace) ExecStr(commandStr string, useLocalScope bool, args... interface{}) InvokeResults {
	command := runspace.CreateCommand()
	defer command.Delete()

	if strings.HasSuffix(commandStr, ".ps1") {
		command.AddCommand(commandStr, useLocalScope)
	} else {
		command.AddScript(commandStr, useLocalScope)
	}
	for _,arg := range args{
		switch v := arg.(type){
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
