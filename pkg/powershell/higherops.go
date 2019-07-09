package powershell

func processArgs(command psCommand, namedArgs map[string]interface{}, args ...interface{}) *InvokeResults {
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			command.AddArgumentString(v)
		case Object:
			command.AddArgument(v)
		default:
			panic("unknown argument type")
		}
	}

	if namedArgs != nil {
		for name, arg := range namedArgs {
			switch v := arg.(type) {
			case string:
				command.AddParameterString(name, v)
			case Object:
				command.AddParameter(name, v)
			default:
				panic("unknown argument type")
			}
		}
	}
	return command.Invoke()
}

// ExecScript - executes a series of statements (not to be confused with .ps1 files which are commands) in powershell
func (runspace Runspace) ExecScript(commandStr string, useLocalScope bool, namedArgs map[string]interface{}, args ...interface{}) *InvokeResults {
	command := runspace.createCommand()
	defer command.Close()
	command.AddScript(commandStr, useLocalScope)
	return processArgs(command, namedArgs, args...)
}

// ExecCommand - executes a command (cmdlets, command files (.ps1), functions, ...) in powershell
func (runspace Runspace) ExecCommand(commandStr string, useLocalScope bool, namedArgs map[string]interface{}, args ...interface{}) *InvokeResults {
	command := runspace.createCommand()
	defer command.Close()
	command.AddCommand(commandStr, useLocalScope)
	return processArgs(command, namedArgs, args...)
}
