package powershell

import (
	"fmt"
	"testing"
)

type fmtPrintLogger struct{}

// GLogInfoLogger is a simple struct that provides ability to send logs to glog at Info level
func (logger fmtPrintLogger) Write(arg string) {
	fmt.Print("    In Logging : ", arg)
}

type recordData struct {
	lines *string
}

func (logger recordData) Write(arg string) {
	*logger.lines = *logger.lines + fmt.Sprintln("    In Logging : ", arg)
}

type callbackTest struct {
	lines *string
}

func (c callbackTest) Callback(str string, input []Object, results CallbackResultsWriter) {
	*c.lines = *c.lines + fmt.Sprintln("    In callback:", str)
	results.WriteString(str)
	for i, object := range input {
		if object.IsNull() {
			*c.lines = *c.lines + fmt.Sprintln("    In callback: index", i, "Object Is Null") // ToString and Type are still valid
		}
		*c.lines = *c.lines + fmt.Sprintln("    In callback: index", i, "type:", object.Type(), "with value:", object.ToString())
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

// PrintAndExecuteCommand executes a command in powershell and prints the results
func PrintAndExecuteCommandlines(lines *string, runspace Runspace, command string, useLocalScope bool, args ...interface{}) {
	*lines = *lines + fmt.Sprint("Executing powershell command:", command, "\n")
	results := runspace.ExecStr(command, useLocalScope, args...)
	defer results.Close()
	*lines = *lines + fmt.Sprint("Completed Executing powershell command:", command, "\n")
	if !results.Success() {
		*lines = *lines + fmt.Sprint("    Command threw exception of type", results.Exception.Type(), "and ToString", results.Exception.ToString())
	} else {
		*lines = *lines + fmt.Sprint("Command returned", len(results.Objects), "objects")
		for i, object := range results.Objects {
			*lines = *lines + fmt.Sprint("    Object", i, "is of type", object.Type(), "and ToString", object.ToString())
		}
	}
}

func validate(t *testing.T, expected string, received string) {
	if expected != received {
		fmt.Printf("expected [%s]\ngot     [%s]1", expected, received)
		t.Fatalf("expected [%s]\ngot     [%s]1", expected, received)
	}
}

func TestCcreateRunspaceWitLoggerWithCallback(t *testing.T) {

	var v string
	lines := &v
	runspace := CreateRunspace(recordData{lines}, callbackTest{lines})
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
	expected := `    In Logging :  Debug: calling Write-Host

    In Logging :  Debug: callng Write-Debug

    In Logging :  calling write-Information
    In Logging :  Verbose: calling Write-Verbose

    In callback: Sending 0 objects to host
    In callback: Sending 2 objects to host and returning them
    In callback: index 0 type: System.String with value: 1
    In callback: index 1 type: System.String with value: 2
`
	validate(t, expected, *lines)
}

// func Example_createRunspaceSimple() {
// 	runspace := CreateRunspaceSimple()
// 	defer runspace.Delete()
// 	results := runspace.ExecStr(`"emit this string"`, true)
// 	defer results.Close()
// 	// print the string result of the first object
// 	fmt.Println(results.Objects[0].ToString())
// 	// Output: emit this string
// }

// func Example_createRunspaceWithArgs() {
// 	runspace := CreateRunspace(fmtPrintLogger{}, callbackTest{})
// 	defer runspace.Delete()
// 	results := runspace.ExecStr(`$args[0]; $args[1] + "changed"`, true, "string1", "string2")
// 	defer results.Close()
// 	if !results.Success() {
// 		runspace.ExecStr(`write-host $args[0]`, true, results.Exception)
// 	}

// 	results2 := runspace.ExecStr(`write-debug $($args[0] + 1); write-debug $args[1]`, false, "myString", results.Objects[1])
// 	defer results2.Close()
// 	// Output:
// 	// In Logging : Debug: myString1
// 	//     In Logging : Debug: string2changed
// }

// func Example_createRunspaceUsingCommand() {
// 	runspace := CreateRunspace(fmtPrintLogger{}, callbackTest{})
// 	defer runspace.Delete()
// 	results := runspace.ExecStr(`..\\..\\tests\\t1.ps1`, false)
// 	defer results.Close()

// 	results2 := runspace.ExecStr(`..\\..\\tests\\t2.ps1`, false)
// 	defer results2.Close()
// 	// Output:
// 	// In Logging : Debug: ab  ba
// 	//     In Logging : Debug: ab  ba
// 	//     In Logging : Debug: ab  ba
// 	//     In callback: asdfasdf
// 	//     In callback: index 0 type: System.Int32 with value: 1
// 	//     In callback: index 1 type: System.Int32 with value: 2
// 	//     In callback: index 2 type: System.Int32 with value: 3
// 	//     In Logging : Debug: asdfasdf 1 2 3
// 	//     In callback: two
// 	//     In Logging : Debug: two
// 	//     In callback: three
// 	//     In Logging : Debug: three
// 	//     In callback: four
// 	//     In callback: index 0 Object Is Null
// 	//     In callback: index 0 type: nullptr with value: nullptr
// 	//     In callback: index 1 Object Is Null
// 	//     In callback: index 1 type: nullptr with value: nullptr
// 	//     In Logging : Debug: ab four   ba
// 	//     In Logging : Error: someerror
// 	//     In Logging : Debug: start t2
// 	//     In Logging : Debug: 5
// 	//     In Logging : Debug: 5
// 	//     In Logging : Debug: asdf
// }

// func Example_globalScope() {
// 	runspace := CreateRunspace(fmtPrintLogger{}, callbackTest{})
// 	defer runspace.Delete()
// 	results := runspace.ExecStr(`..\\..\\tests\\test_scope.ps1`, false)
// 	defer results.Close()

// 	results2 := runspace.ExecStr(`..\\..\\tests\\test_scope.ps1`, false)
// 	defer results2.Close()
// 	// Output:
// 	// In Logging : Debug: ab  ba
// 	//     In Logging : Debug: ab  ba
// 	//     In Logging : Debug: ab  ba
// 	//     In Logging : Debug: ab  ba
// 	//     In Logging : Debug: ab 1 ba
// 	//     In Logging : Debug: ab 2 ba
// 	//     In Logging : Debug: ab 3 ba
// 	//     In Logging : Debug: ab 4 ba
// }

// func Example_localScope() {
// 	runspace := CreateRunspace(fmtPrintLogger{}, callbackTest{})
// 	defer runspace.Delete()
// 	results := runspace.ExecStr(`..\\..\\tests\\test_scope.ps1`, true)
// 	defer results.Close()

// 	results2 := runspace.ExecStr(`..\\..\\tests\\test_scope.ps1`, true)
// 	defer results2.Close()
// 	// Output:
// 	// In Logging : Debug: ab  ba
// 	//     In Logging : Debug: ab  ba
// 	//     In Logging : Debug: ab  ba
// 	//     In Logging : Debug: ab  ba
// 	//     In Logging : Debug: ab  ba
// 	//     In Logging : Debug: ab  ba
// 	//     In Logging : Debug: ab 3 ba
// 	//     In Logging : Debug: ab  ba
// }

// type callbackAddRef struct{}

// func (c callbackAddRef) Callback(str string, input []Object, results CallbackResultsWriter) {
// 	for _, object := range input {
// 		results.Write(object.AddRef(), true)
// 	}
// }

// func Example_callbackWriteTrue() {
// 	runspace := CreateRunspace(fmtPrintLogger{}, callbackAddRef{})
// 	defer runspace.Delete()
// 	results := runspace.ExecStr(`1 | send-hostcommand -message 'empty'`, true)
// 	defer results.Close()
// 	// Output:
// }

// type callbackAddRefSave struct {
// 	objects *[]Object
// }

// func (c callbackAddRefSave) Callback(str string, input []Object, results CallbackResultsWriter) {
// 	for _, object := range input {
// 		*c.objects = append(*c.objects, object.AddRef())
// 	}
// }
// func Example_callbackSaveObject() {
// 	var callback callbackAddRefSave
// 	callback.objects = &([]Object{})
// 	runspace := CreateRunspace(fmtPrintLogger{}, callback)
// 	defer runspace.Delete()
// 	results := runspace.ExecStr(`1 | send-hostcommand -message 'empty'`, true)
// 	defer results.Close()
// 	results2 := runspace.ExecStr(`write-host $args[0]`, true, (*callback.objects)[0])
// 	defer results2.Close()
// 	for _, object := range *callback.objects {
// 		object.Close()
// 	}
// 	// Output: In Logging : Debug: 1
// }
