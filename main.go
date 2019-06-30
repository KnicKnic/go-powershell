package main

// "bitbucket.org/creachadair/shell"
import (
	"flag"
	"strings"

	. "./powershell"

	"github.com/golang/glog"
)

type callbackTest struct{}

func (c callbackTest) Callback(str string, input []PowershellObject, results CallbackResultsWriter) {
	glog.Info("In callback: ", str)
	results.WriteString(str)
	for _, object := range input {
		results.Write(object)
	}
}

// Example on how to use powershell wrappers
func Example() {
	runspace := CreateRunspace(GLogInfoLogger{}, callbackTest{})
	defer runspace.Delete()

	for i := 0; i < len(commandFlags); i++ {
		commandFlags[i] = strings.ReplaceAll(commandFlags[i], "\\", "\\\\")
		runspace.ExecStr(commandFlags[i], *useLocalScope)
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
