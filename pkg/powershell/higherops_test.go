package powershell

import (
	"fmt"
	"testing"
)

type recordData struct {
	lines string
}

func (logger *recordData) Write(arg string) {
	logger.lines = logger.lines + fmt.Sprintln("    In Logging : ", arg)
}

func (logger *recordData) Println(args ...interface{}) {
	logger.lines = logger.lines + fmt.Sprintln(args...)
	fmt.Println(args...)
}
func (logger *recordData) Print(args ...interface{}) {
	logger.lines = logger.lines + fmt.Sprint(args...)
	fmt.Println(args...)
}
func (logger *recordData) Reset() {
	logger.lines = ""
}

var recordPtrInstance recordData
var record *recordData = &recordPtrInstance

type fmtPrintLogger struct{}

func (logger fmtPrintLogger) Write(arg string) {
	record.Print("    In Logging : ", arg)
}

type callbackTest struct {
	lines *string
}

func (c callbackTest) Callback(str string, input []Object, results CallbackResultsWriter) {
	record.Println("    In callback:", str)
	results.WriteString(str)
	for i, object := range input {
		if object.IsNull() {
			record.Println("    In callback: index", i, "Object Is Null") // ToString and Type are still valid
		}
		record.Println("    In callback: index", i, "type:", object.Type(), "with value:", object.ToString())
		results.Write(object, false)
	}
}

// PrintAndExecuteCommand executes a command in powershell and prints the results
func PrintAndExecuteCommand(runspace Runspace, command string, useLocalScope bool, args ...interface{}) {
	record.Print("Executing powershell command:", command, "\n")
	results := runspace.ExecStr(command, useLocalScope, args...)
	defer results.Close()
	record.Print("Completed Executing powershell command:", command, "\n")
	if !results.Success() {
		record.Println("    Command threw exception of type", results.Exception.Type(), "and ToString", results.Exception.ToString())
	} else {
		record.Println("Command returned", len(results.Objects), "objects")
		for i, object := range results.Objects {
			record.Println("    Object", i, "is of type", object.Type(), "and ToString", object.ToString())
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
	record.Reset()
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
	expected := `Executing powershell command:
write-host "calling Write-Host"
write-debug "callng Write-Debug"
write-Information "calling write-Information"
write-verbose "calling Write-Verbose"
Send-HostCommand -message "Sending 0 objects to host" | out-null
"1","2" | Send-HostCommand -message "Sending 2 objects to host and returning them"
    In Logging : Debug: calling Write-Host
    In Logging : Debug: callng Write-Debug
    In Logging : calling write-Information    In Logging : Verbose: calling Write-Verbose
    In callback: Sending 0 objects to host
    In callback: Sending 2 objects to host and returning them
    In callback: index 0 type: System.String with value: 1
    In callback: index 1 type: System.String with value: 2
Completed Executing powershell command:
write-host "calling Write-Host"
write-debug "callng Write-Debug"
write-Information "calling write-Information"
write-verbose "calling Write-Verbose"
Send-HostCommand -message "Sending 0 objects to host" | out-null
"1","2" | Send-HostCommand -message "Sending 2 objects to host and returning them"
Command returned 3 objects
    Object 0 is of type System.String and ToString Sending 2 objects to host and returning them
    Object 1 is of type System.String and ToString 1
    Object 2 is of type System.String and ToString 2
`
	validate(t, expected, record.lines)
}

func TestCcreateRunspaceSimple(t *testing.T) {
	record.Reset()
	runspace := CreateRunspaceSimple()
	defer runspace.Delete()
	results := runspace.ExecStr(`"emit this string"`, true)
	defer results.Close()
	// print the string result of the first object
	record.Print(results.Objects[0].ToString())
	validate(t, "emit this string", record.lines)
}

func TestCcreateRunspaceWithArgs(t *testing.T) {
	record.Reset()
	runspace := CreateRunspace(fmtPrintLogger{}, callbackTest{})
	defer runspace.Delete()
	results := runspace.ExecStr(`$args[0]; $args[1] + "changed"`, true, "string1", "string2")
	defer results.Close()
	if !results.Success() {
		runspace.ExecStr(`write-host $args[0]`, true, results.Exception)
	}

	results2 := runspace.ExecStr(`write-debug $($args[0] + 1); write-debug $args[1]`, false, "myString", results.Objects[1])
	defer results2.Close()

	expected := `    In Logging : Debug: myString1
    In Logging : Debug: string2changed
`
	validate(t, expected, record.lines)
}

func TestCcreateRunspaceUsingCommand(t *testing.T) {
	record.Reset()
	runspace := CreateRunspace(fmtPrintLogger{}, callbackTest{})
	defer runspace.Delete()
	results := runspace.ExecStr(`..\\..\\tests\\t1.ps1`, false)
	defer results.Close()

	results2 := runspace.ExecStr(`..\\..\\tests\\t2.ps1`, false)
	defer results2.Close()
	// Output:
	expected := `    In Logging : Debug: ab  ba
    In Logging : Debug: ab  ba
    In Logging : Debug: ab  ba
    In callback: asdfasdf
    In callback: index 0 type: System.Int32 with value: 1
    In callback: index 1 type: System.Int32 with value: 2
    In callback: index 2 type: System.Int32 with value: 3
    In Logging : Debug: asdfasdf 1 2 3
    In callback: two
    In Logging : Debug: two
    In callback: three
    In Logging : Debug: three
    In callback: four
    In callback: index 0 Object Is Null
    In callback: index 0 type: nullptr with value: nullptr
    In callback: index 1 Object Is Null
    In callback: index 1 type: nullptr with value: nullptr
    In Logging : Debug: ab four   ba
    In Logging : Error: someerror
    In Logging : Debug: start t2
    In Logging : Debug: 5
    In Logging : Debug: 5
    In Logging : Debug: asdf
`
	validate(t, expected, record.lines)
}

func TestCglobalScope(t *testing.T) {
	record.Reset()
	runspace := CreateRunspace(fmtPrintLogger{}, callbackTest{})
	defer runspace.Delete()
	results := runspace.ExecStr(`..\\..\\tests\\test_scope.ps1`, false)
	defer results.Close()

	results2 := runspace.ExecStr(`..\\..\\tests\\test_scope.ps1`, false)
	defer results2.Close()
	// Output:
	expected := `    In Logging : Debug: ab  ba
    In Logging : Debug: ab  ba
    In Logging : Debug: ab  ba
    In Logging : Debug: ab  ba
    In Logging : Debug: ab 1 ba
    In Logging : Debug: ab 2 ba
    In Logging : Debug: ab 3 ba
    In Logging : Debug: ab 4 ba
`
	validate(t, expected, record.lines)
}

func TestClocalScope(t *testing.T) {
	record.Reset()
	runspace := CreateRunspace(fmtPrintLogger{}, callbackTest{})
	defer runspace.Delete()
	results := runspace.ExecStr(`..\\..\\tests\\test_scope.ps1`, true)
	defer results.Close()

	results2 := runspace.ExecStr(`..\\..\\tests\\test_scope.ps1`, true)
	defer results2.Close()

	expected := `    In Logging : Debug: ab  ba
    In Logging : Debug: ab  ba
    In Logging : Debug: ab  ba
    In Logging : Debug: ab  ba
    In Logging : Debug: ab  ba
    In Logging : Debug: ab  ba
    In Logging : Debug: ab 3 ba
    In Logging : Debug: ab  ba
`
	validate(t, expected, record.lines)
}

type callbackAddRef struct{}

func (c callbackAddRef) Callback(str string, input []Object, results CallbackResultsWriter) {
	results.WriteString(str)
	for _, object := range input {
		results.Write(object.AddRef(), true)
	}
}

func TestCcallbackWriteTrue(t *testing.T) {
	record.Reset()
	runspace := CreateRunspace(fmtPrintLogger{}, callbackAddRef{})
	defer runspace.Delete()
	results := runspace.ExecStr(`1 | send-hostcommand -message 'empty'`, true)
	defer results.Close()
	// Output:
	validate(t, results.Objects[0].ToString(), "empty")
	validate(t, results.Objects[1].ToString(), "1")
	if len(results.Objects) != 2 {
		t.Fail()
	}
	validate(t, "", record.lines)
}

type callbackAddRefSave struct {
	objects []Object
}

func (c *callbackAddRefSave) Callback(str string, input []Object, results CallbackResultsWriter) {
	for _, object := range input {
		c.objects = append(c.objects, object.AddRef())
	}
}

func psObjectsToInterface(obj []Object) []interface{} {
	s := make([]interface{}, len(obj))
	for i, v := range obj {
		s[i] = v
	}
	return s
}

func TestCcallbackSaveObject(t *testing.T) {
	callback := callbackAddRefSave{}
	record.Reset()
	runspace := CreateRunspace(fmtPrintLogger{}, &callback)
	defer runspace.Delete()
	results := runspace.ExecStr(`1 | send-hostcommand -message 'empty'`, true)
	defer results.Close()
	results2 := runspace.ExecStr(`$args | %{write-host $_}`, true, psObjectsToInterface(callback.objects)...)
	defer results2.Close()
	for _, object := range callback.objects {
		object.Close()
	}
	// Output: In Logging : Debug: 1
	expected := `    In Logging : Debug: 1
`
	validate(t, expected, record.lines)
}
