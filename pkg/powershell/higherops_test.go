package powershell

import (
	"fmt"
)

// GLogInfoLogger is a simple struct that provides ability to send logs to glog at Info level
type fmtPrintLogger struct {
}

func (logger fmtPrintLogger) Write(arg string) {
	fmt.Print("    In Logging : ", arg)
}

type callbackTest struct{}

func (c callbackTest) Callback(str string, input []Object, results CallbackResultsWriter) {
	fmt.Println("    In callback:", str)
	results.WriteString(str)
	for i, object := range input {
		if object.IsNull() {
			fmt.Println("    In callback: index", i, "Object Is Null") // ToString and Type are still valid
		}
		fmt.Println("    In callback: index", i, "type:", object.Type(), "with value:", object.ToString())
		results.Write(object, false)
	}
}

// PrintAndExecuteCommand executes a command in powershell and prints the results
func PrintAndExecuteCommand(runspace Runspace, command string, useLocalScope bool) {
	fmt.Print("Executing powershell command:", command, "\n")
	results := runspace.ExecStr(command, useLocalScope)
	defer results.Close()
	fmt.Print("Completed Executing powershell command:", command, "\n")
	if !results.Success() {
		fmt.Println("    Command threw exception of type", results.Exception.Type(), "and ToString", results.Exception.ToString())
	} else {
		fmt.Println("Command returned", len(results.Objects), "objects")
		for i, object := range results.Objects {
			fmt.Println("    Object", i, "is of type", object.Type(), "and ToString", object.ToString())
		}
	}
}

func ExampleCreateRunspaceWitLoggerWithCallback(){

	runspace := CreateRunspace(fmtPrintLogger{}, callbackTest{})
	// runspace := CreateRunspaceSimple()
	defer runspace.Delete()
input_str :=`
write-host "calling Write-Host"
write-debug "callng Write-Debug"
write-Information "calling write-Information"
write-verbose "calling Write-Verbose"
Send-HostCommand -message "Sending 0 objects to host" | out-null
"1","2" | Send-HostCommand -message "Sending 2 objects to host and returning them"`
PrintAndExecuteCommand(runspace, input_str, false)
	// Output:
// Executing powershell command:
// write-host "calling Write-Host"
// write-debug "callng Write-Debug"
// write-Information "calling write-Information"
// write-verbose "calling Write-Verbose"
// Send-HostCommand -message "Sending 0 objects to host" | out-null
// "1","2" | Send-HostCommand -message "Sending 2 objects to host and returning them"
//     In Logging : Debug: calling Write-Host
//     In Logging : Debug: callng Write-Debug
//     In Logging : calling write-Information    In Logging : Verbose: calling Write-Verbose
//     In callback: Sending 0 objects to host
//     In callback: Sending 2 objects to host and returning them
//     In callback: index 0 type: System.String with value: 1
//     In callback: index 1 type: System.String with value: 2
// Completed Executing powershell command:
// write-host "calling Write-Host"
// write-debug "callng Write-Debug"
// write-Information "calling write-Information"
// write-verbose "calling Write-Verbose"
// Send-HostCommand -message "Sending 0 objects to host" | out-null
// "1","2" | Send-HostCommand -message "Sending 2 objects to host and returning them"
// Command returned 3 objects
//     Object 0 is of type System.String and ToString Sending 2 objects to host and returning them
//     Object 1 is of type System.String and ToString 1
//     Object 2 is of type System.String and ToString 2
 }


func ExampleCreateRunspaceSimple(){
	runspace := CreateRunspaceSimple()
	defer runspace.Delete()
	results := runspace.ExecStr( `"emit this string"`, true)
	defer results.Close()
	// print the string result of the first object
	fmt.Println(results.Objects[0].ToString())
	// Output: emit this string
}