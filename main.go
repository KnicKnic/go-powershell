package main

// "bitbucket.org/creachadair/shell"
import (
	"flag"
	"fmt"
	"strings"

	. "./powershell"
)

// GLogInfoLogger is a simple struct that provides ability to send logs to glog at Info level
type fmtPrintLogger struct {
}

func (logger fmtPrintLogger) Write(args ...interface{}) {
	fmt.Print("\tIn Logging : ", fmt.Sprint(args...))
}

type callbackTest struct{}

func (c callbackTest) Callback(str string, input []PowershellObject, results CallbackResultsWriter) {
	fmt.Println("\tIn callback:", str)
	results.WriteString(str)
	for i, object := range input {
		if object.IsNull() {
			fmt.Println("\tIn callback: index", i, "Object Is Null") // ToString and Type are still valid
		}
		fmt.Println("\tIn callback: index", i, "type:", object.Type(), "with value:", object.ToString())
		results.Write(object, false)
	}
}

func PrintAndExecuteCommand(runspace Runspace, command string, useLocalScope bool) {
	fmt.Println("Executing powershell command:", command)
	results := runspace.ExecStr(command, useLocalScope)
	defer results.Close()
	fmt.Println("Completed Executing powershell command:", command)
	if !results.Success() {
		fmt.Println("\tCommand threw exception of type", results.Exception.Type(), "and ToString", results.Exception.ToString())
	} else {
		fmt.Println("Command returned", len(results.Objects), "objects")
		for i, object := range results.Objects {
			fmt.Println("\tObject", i, "is of type", object.Type(), "and ToString", object.ToString())
		}
	}
}

// Example on how to use powershell wrappers
func Example() {
	runspace := CreateRunspace(fmtPrintLogger{}, callbackTest{})
	// runspace := CreateRunspaceSimple()
	defer runspace.Delete()

	for i := 0; i < len(commandFlags); i++ {
		command := strings.ReplaceAll(commandFlags[i], "\\", "\\\\")
		PrintAndExecuteCommand(runspace, command, *useLocalScope)

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
	Example()
}
