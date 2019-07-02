package main

// "bitbucket.org/creachadair/shell"
import (
	"flag"
	"strings"

	. "./powershell"

	"github.com/golang/glog"
)

// GLogInfoLogger is a simple struct that provides ability to send logs to glog at Info level
type gLogInfoLogger struct {
}

func (logger gLogInfoLogger) Write(args ...interface{}) {
	glog.Info(args...)
}

type callbackTest struct{}

func (c callbackTest) Callback(str string, input []PowershellObject, results CallbackResultsWriter) {
	glog.Info("In callback: ", str)
	results.WriteString(str)
	for i, object := range input {
		if object.IsNull() {
			glog.Info("In callback: index ", i, " Object Is Null: ") // ToString and Type are still valid
		}
		glog.Info("In callback: index ", i, " type: ", object.Type(), " with value: ", object.ToString())
		results.Write(object, false)
	}
}

// Example on how to use powershell wrappers
func Example() {
	runspace := CreateRunspace(gLogInfoLogger{}, callbackTest{})
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
