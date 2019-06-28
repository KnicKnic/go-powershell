package main

// "bitbucket.org/creachadair/shell"
import (
	"flag"
	"strings"

	"github.com/golang/glog"
)

/*

#cgo CFLAGS: -I.
#cgo LDFLAGS: ./psh_host.dll


#include <stddef.h>
#include "powershell.h"

*/
import "C"

// ExecStr - executes a commandline in powershell
func (runspace Runspace) ExecStr(commandStr string) {
	command := runspace.CreatePowershellCommand()
	defer command.Delete()

	// fields, ok := shell.Split(commandStr)
	// if !ok {
	// 	panic("command was invalid {" + commandStr + "}")
	// }

	if strings.HasSuffix(commandStr, ".ps1") {
		command.AddCommand(commandStr, *useLocalScope)
	} else {
		command.AddScript(commandStr, *useLocalScope)
	}
	// for i := 1; i < len(fields); i++ {
	// 	command.AddArgument(fields[i])
	// }
	command.Invoke()
}

// Example on how to use powershell wrappers
func Example() {
	C.InitLibraryHelper()
	runspace := CreateRunspace()
	defer runspace.Delete()

	for i := 0; i < len(commandFlags); i++ {
		commandFlags[i] = strings.ReplaceAll(commandFlags[i], "\\", "\\\\")
		runspace.ExecStr(commandFlags[i])
	}
}

type arrayCommandFlags []string

func (i *arrayCommandFlags) String() string {
	return "my string representation"
}

func (i *arrayCommandFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var commandFlags arrayCommandFlags
var useLocalScope = flag.Bool("useLocalScope", false, "True if should execute scripts in the local scope")

func main() {
	flag.Var(&commandFlags, "command", "Command to run in powershell")
	flag.Parse()
	glog.Info(*useLocalScope, commandFlags)
	Example()
}
