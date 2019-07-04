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
func PrintAndExecuteCommand(runspace Runspace, command string, useLocalScope bool, args ...interface{}) {
	fmt.Print("Executing powershell command:", command, "\n")
	results := runspace.ExecStr(command, useLocalScope, args...)
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

func Example_createRunspaceWitLoggerWithCallback() {

	runspace := CreateRunspace(fmtPrintLogger{}, callbackTest{})
	// runspace := CreateRunspaceSimple()
	defer runspace.Delete()
	inputStr := `
write-host "calling Write-Host"
write-debug "callng Write-Debug"
write-Information "calling write-Information"
write-verbose "calling Write-Verbose"
Send-HostCommand -message "Sending 0 objects to host" | out-null
"1","2" | Send-HostCommand -message "Sending 2 objects to host and returning them"`
	PrintAndExecuteCommand(runspace, inputStr, false)
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

func Example_createRunspaceSimple() {
	runspace := CreateRunspaceSimple()
	defer runspace.Delete()
	results := runspace.ExecStr(`"emit this string"`, true)
	defer results.Close()
	// print the string result of the first object
	fmt.Println(results.Objects[0].ToString())
	// Output: emit this string
}

func Example_createRunspaceWithArgs() {
	runspace := CreateRunspace(fmtPrintLogger{}, callbackTest{})
	defer runspace.Delete()
	results := runspace.ExecStr(`$args[0]; $args[1] + "changed"`, true, "string1", "string2")
	defer results.Close()
	if !results.Success() {
		runspace.ExecStr(`write-host $args[0]`, true, results.Exception)
	}

	results2 := runspace.ExecStr(`write-debug $($args[0] + 1); write-debug $args[1]`, false, "myString", results.Objects[1])
	defer results2.Close()
	// Output:
	// In Logging : Debug: myString1
	//     In Logging : Debug: string2changed
}

func Example_createRunspaceUsingCommand() {
	runspace := CreateRunspace(fmtPrintLogger{}, callbackTest{})
	defer runspace.Delete()
	results := runspace.ExecStr(`..\\..\\tests\\t1.ps1`, false)
	defer results.Close()

	results2 := runspace.ExecStr(`..\\..\\tests\\t2.ps1`, false)
	defer results2.Close()
	// Output:
	// In Logging : Debug: ab  ba
	//     In Logging : Debug: ab  ba
	//     In Logging : Debug: ab  ba
	//     In callback: asdfasdf
	//     In callback: index 0 type: System.Int32 with value: 1
	//     In callback: index 1 type: System.Int32 with value: 2
	//     In callback: index 2 type: System.Int32 with value: 3
	//     In Logging : Debug: asdfasdf 1 2 3
	//     In callback: two
	//     In Logging : Debug: two
	//     In callback: three
	//     In Logging : Debug: three
	//     In callback: four
	//     In callback: index 0 Object Is Null
	//     In callback: index 0 type: nullptr with value: nullptr
	//     In callback: index 1 Object Is Null
	//     In callback: index 1 type: nullptr with value: nullptr
	//     In Logging : Debug: ab four   ba
	//     In Logging : Error: someerror
	//     In Logging : Debug: start t2
	//     In Logging : Debug: 5
	//     In Logging : Debug: 5
	//     In Logging : Debug: asdf
}

func Example_globalScope() {
	runspace := CreateRunspace(fmtPrintLogger{}, callbackTest{})
	defer runspace.Delete()
	results := runspace.ExecStr(`..\\..\\tests\\test_scope.ps1`, false)
	defer results.Close()

	results2 := runspace.ExecStr(`..\\..\\tests\\test_scope.ps1`, false)
	defer results2.Close()
	// Output:
	// In Logging : Debug: ab  ba
	//     In Logging : Debug: ab  ba
	//     In Logging : Debug: ab  ba
	//     In Logging : Debug: ab  ba
	//     In Logging : Debug: ab 1 ba
	//     In Logging : Debug: ab 2 ba
	//     In Logging : Debug: ab 3 ba
	//     In Logging : Debug: ab 4 ba
}

func Example_localScope() {
	runspace := CreateRunspace(fmtPrintLogger{}, callbackTest{})
	defer runspace.Delete()
	results := runspace.ExecStr(`..\\..\\tests\\test_scope.ps1`, true)
	defer results.Close()

	results2 := runspace.ExecStr(`..\\..\\tests\\test_scope.ps1`, true)
	defer results2.Close()
	// Output:
	// In Logging : Debug: ab  ba
	//     In Logging : Debug: ab  ba
	//     In Logging : Debug: ab  ba
	//     In Logging : Debug: ab  ba
	//     In Logging : Debug: ab  ba
	//     In Logging : Debug: ab  ba
	//     In Logging : Debug: ab 3 ba
	//     In Logging : Debug: ab  ba
}

type callbackAddRef struct{}

func (c callbackAddRef) Callback(str string, input []Object, results CallbackResultsWriter) {
	for _, object := range input {
		results.Write(object.AddRef(), true)
	}
}

func Example_callbackWriteTrue() {
	runspace := CreateRunspace(fmtPrintLogger{}, callbackAddRef{})
	defer runspace.Delete()
	results := runspace.ExecStr(`1 | send-hostcommand -message 'empty'`, true)
	defer results.Close()
	// Output:
}

type callbackAddRefSave struct {
	objects *[]Object
}

func (c callbackAddRefSave) Callback(str string, input []Object, results CallbackResultsWriter) {
	for _, object := range input {
		*c.objects = append(*c.objects, object.AddRef())
	}
}
func Example_callbackSaveObject() {
	var callback callbackAddRefSave
	callback.objects = &([]Object{})
	runspace := CreateRunspace(fmtPrintLogger{}, callback)
	defer runspace.Delete()
	results := runspace.ExecStr(`1 | send-hostcommand -message 'empty'`, true)
	defer results.Close()
	results2 := runspace.ExecStr(`write-host $args[0]`, true, (*callback.objects)[0])
	defer results2.Close()
	for _, object := range *callback.objects {
		object.Close()
	}
	// Output: In Logging : Debug: 1
}
