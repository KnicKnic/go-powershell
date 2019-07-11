package powershell

import (
	"encoding/json"
)

func processArgs(command psCommand, jsonMarshalUnknowns bool, namedArgs map[string]interface{}, args ...interface{}) *InvokeResults {
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			command.AddArgumentString(v)
		case Object:
			command.AddArgument(v)
		default:
			if jsonMarshalUnknowns {
				bytes, err := json.Marshal(arg)
				if err != nil {
					panic(err)
				}
				command.AddArgumentString(string(bytes))
			} else {
				panic("unknown argument type")
			}
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
				if jsonMarshalUnknowns {
					bytes, err := json.Marshal(arg)
					if err != nil {
						panic(err)
					}
					command.AddParameterString(name, string(bytes))
				} else {
					panic("unknown argument type")
				}
			}
		}
	}
	return command.Invoke()
}

// ExecScript - executes a series of statements (not to be confused with .ps1 files which are commands) in powershell
//
// useLocalScope - true means to create a child scope, false means to use the global scope (this is probably what you want)
//
// although namedArgs and args both state that they take interface, the only types they currently take are string and Object (a result from a previous powershell invocation)
func (runspace Runspace) ExecScript(commandStr string, useLocalScope bool, namedArgs map[string]interface{}, args ...interface{}) *InvokeResults {
	command := runspace.createCommand()
	defer command.Close()
	command.AddScript(commandStr, useLocalScope)
	return processArgs(command, false, namedArgs, args...)
}

// ExecCommand - executes a command (cmdlets, command files (.ps1), functions, ...) in powershell
//
// useLocalScope - true means to create a child scope, false means to use the global scope (this is probably what you want)
//
// although namedArgs and args both state that they take interface, the only types they currently take are string and Object (a result from a previous powershell invocation)
func (runspace Runspace) ExecCommand(commandStr string, useLocalScope bool, namedArgs map[string]interface{}, args ...interface{}) *InvokeResults {
	command := runspace.createCommand()
	defer command.Close()
	command.AddCommand(commandStr, useLocalScope)
	return processArgs(command, false, namedArgs, args...)
}

// ExecScriptJSONMarshalUnknown - executes a series of statements (not to be confused with .ps1 files which are commands) in powershell
//
// useLocalScope - true means to create a child scope, false means to use the global scope (this is probably what you want)
//
// Although namedArgs and args both state that they take interface, string and Object (a result from a previous powershell invocation) will get sent to powershell directly. Any other type will first get marshaled to json using json.Marshal and sent as a string to powershell
func (runspace Runspace) ExecScriptJSONMarshalUnknown(commandStr string, useLocalScope bool, namedArgs map[string]interface{}, args ...interface{}) *InvokeResults {
	command := runspace.createCommand()
	defer command.Close()
	command.AddScript(commandStr, useLocalScope)
	return processArgs(command, true, namedArgs, args...)
}

// ExecCommandJSONMarshalUnknown - executes a command (cmdlets, command files (.ps1), functions, ...) in powershell
//
// useLocalScope - true means to create a child scope, false means to use the global scope (this is probably what you want)
//
// Although namedArgs and args both state that they take interface, string and Object (a result from a previous powershell invocation) will get sent to powershell directly. Any other type will first get marshaled to json using json.Marshal and sent as a string to powershell
func (runspace Runspace) ExecCommandJSONMarshalUnknown(commandStr string, useLocalScope bool, namedArgs map[string]interface{}, args ...interface{}) *InvokeResults {
	command := runspace.createCommand()
	defer command.Close()
	command.AddCommand(commandStr, useLocalScope)
	return processArgs(command, true, namedArgs, args...)
}
