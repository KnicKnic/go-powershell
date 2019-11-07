package powershell

import (
	"encoding/json"
	"errors"
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

var record = &recordData{}

type fmtPrintLogger struct{}

func (logger fmtPrintLogger) Write(arg string) {
	record.Print("    In Logging : ", arg)
}

type callbackTest struct {
	lines *string
}

func (c callbackTest) Callback(_ Runspace, str string, input []Object, results CallbackResultsWriter) {
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
	results := runspace.ExecScript(command, useLocalScope, nil, args...)
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
	defer runspace.Close()
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
    In Logging : calling Write-Host
    In Logging : Debug: callng Write-Debug
    In Logging : calling write-Information
    In Logging : Verbose: calling Write-Verbose
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

type fullLogger struct {
}

func (_ fullLogger) Warning(arg string) {

	record.Print("    In Logging : Warning: ", arg)
}
func (_ fullLogger) Information(arg string) {
	record.Print("    In Logging : Information: ", arg)
}
func (_ fullLogger) Verbose(arg string) {
	record.Print("    In Logging : Verbose: ", arg)
}
func (_ fullLogger) Debug(arg string) {
	record.Print("    In Logging : Debug: ", arg)
}
func (_ fullLogger) Error(arg string) {
	record.Print("    In Logging : Error: ", arg)
}
func (_ fullLogger) Write(arg string) {
	record.Print("    In Logging : Write: ", arg)
}

func (_ fullLogger) Warningln(arg string) {
	arg = arg + "\n"
	record.Print("    In Logging : Line Warning: ", arg)
}
func (_ fullLogger) Informationln(arg string) {
	arg = arg + "\n"
	record.Print("    In Logging : Line Information: ", arg)
}
func (_ fullLogger) Verboseln(arg string) {
	arg = arg + "\n"
	record.Print("    In Logging : Line Verbose: ", arg)
}
func (_ fullLogger) Debugln(arg string) {
	arg = arg + "\n"
	record.Print("    In Logging : Line Debug: ", arg)
}
func (_ fullLogger) Errorln(arg string) {
	arg = arg + "\n"
	record.Print("    In Logging : Line Error: ", arg)
}
func (_ fullLogger) Writeln(arg string) {
	arg = arg + "\n"
	record.Print("    In Logging : Line Write: ", arg)
}

func TestCcreateRunspaceWitLoggerWithFullCallback(t *testing.T) {
	record.Reset()
	runspace := CreateRunspace(fullLogger{}, callbackTest{})
	// runspace := CreateRunspaceSimple()
	defer runspace.Close()
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
    In Logging : Line Write: calling Write-Host
    In Logging : Line Debug: callng Write-Debug
    In Logging : Line Write: calling write-Information
    In Logging : Line Verbose: calling Write-Verbose
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
	defer runspace.Close()
	results := runspace.ExecScript(`"emit this string"`, true, nil)
	defer results.Close()
	// print the string result of the first object
	record.Print(results.Objects[0].ToString())
	validate(t, "emit this string", record.lines)
}

func TestCcreateRunspaceWithArgs(t *testing.T) {
	record.Reset()
	runspace := CreateRunspace(fmtPrintLogger{}, callbackTest{})
	defer runspace.Close()
	results := runspace.ExecScript(`$args[0]; $args[1] + "changed"`, true, nil, "string1", "string2")
	defer results.Close()
	if !results.Success() {
		runspace.ExecScript(`write-host $args[0]`, true, nil, results.Exception)
	}

	results2 := runspace.ExecScript(`write-debug $($args[0] + 1); write-debug $args[1]`, false, nil, "myString", results.Objects[1])
	defer results2.Close()

	expected := `    In Logging : Debug: myString1
    In Logging : Debug: string2changed
`
	validate(t, expected, record.lines)
}

func TestCcreateRunspaceUsingCommand(t *testing.T) {
	record.Reset()
	runspace := CreateRunspace(fmtPrintLogger{}, callbackTest{})
	defer runspace.Close()
	results := runspace.ExecCommand(`..\\..\\tests\\t1.ps1`, false, nil)
	defer results.Close()

	results2 := runspace.ExecCommand(`..\\..\\tests\\t2.ps1`, false, nil)
	defer results2.Close()
	// Output:
	expected := `    In Logging : ab  ba
    In Logging : ab  ba
    In Logging : ab  ba
    In callback: asdfasdf
    In callback: index 0 type: System.Int32 with value: 1
    In callback: index 1 type: System.Int32 with value: 2
    In callback: index 2 type: System.Int32 with value: 3
    In Logging : asdfasdf 1 2 3
    In callback: two
    In Logging : two
    In callback: three
    In Logging : three
    In callback: four
    In callback: index 0 Object Is Null
    In callback: index 0 type: nullptr with value: nullptr
    In callback: index 1 Object Is Null
    In callback: index 1 type: nullptr with value: nullptr
    In Logging : ab four   ba
    In Logging : Error: someerror
    In Logging : start t2
    In Logging : 5
    In Logging : 5
    In Logging : asdf
`
	validate(t, expected, record.lines)
}

func TestCglobalScope(t *testing.T) {
	record.Reset()
	runspace := CreateRunspace(fmtPrintLogger{}, callbackTest{})
	defer runspace.Close()
	results := runspace.ExecCommand(`..\\..\\tests\\test_scope.ps1`, false, nil)
	defer results.Close()

	results2 := runspace.ExecCommand(`..\\..\\tests\\test_scope.ps1`, false, nil)
	defer results2.Close()
	// Output:
	expected := `    In Logging : ab  ba
    In Logging : ab  ba
    In Logging : ab  ba
    In Logging : ab  ba
    In Logging : ab 1 ba
    In Logging : ab 2 ba
    In Logging : ab 3 ba
    In Logging : ab 4 ba
`
	validate(t, expected, record.lines)
}

func TestClocalScope(t *testing.T) {
	record.Reset()
	runspace := CreateRunspace(fmtPrintLogger{}, callbackTest{})
	defer runspace.Close()
	results := runspace.ExecCommand(`..\\..\\tests\\test_scope.ps1`, true, nil)
	defer results.Close()

	results2 := runspace.ExecCommand(`..\\..\\tests\\test_scope.ps1`, true, nil)
	defer results2.Close()

	expected := `    In Logging : ab  ba
    In Logging : ab  ba
    In Logging : ab  ba
    In Logging : ab  ba
    In Logging : ab  ba
    In Logging : ab  ba
    In Logging : ab 3 ba
    In Logging : ab  ba
`
	validate(t, expected, record.lines)
}

type callbackAddRef struct{}

func (c callbackAddRef) Callback(_ Runspace, str string, input []Object, results CallbackResultsWriter) {
	results.WriteString(str)
	for _, object := range input {
		results.Write(object.AddRef(), true)
	}
}

func TestCcallbackWriteTrue(t *testing.T) {
	record.Reset()
	runspace := CreateRunspace(fmtPrintLogger{}, callbackAddRef{})
	defer runspace.Close()
	results := runspace.ExecScript(`1 | send-hostcommand -message 'empty'`, true, nil)
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

func (c *callbackAddRefSave) Callback(runspace Runspace, str string, input []Object, results CallbackResultsWriter) {
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
	defer runspace.Close()
	results := runspace.ExecScript(`1 | send-hostcommand -message 'empty'`, true, nil)
	defer results.Close()
	results2 := runspace.ExecScript(`$args | %{write-host $_}`, true, nil, psObjectsToInterface(callback.objects)...)
	defer results2.Close()
	for _, object := range callback.objects {
		object.Close()
	}
	// Output: In Logging : Debug: 1
	expected := `    In Logging : 1
`
	validate(t, expected, record.lines)
}

func TestCpowershellCommandWithNamedParameters(t *testing.T) {
	record.Reset()
	// create a runspace (where you run your powershell statements in)
	runspace := CreateRunspace(fmtPrintLogger{}, nil)
	// auto cleanup your runspace
	defer runspace.Close()

	paramResults := runspace.ExecScript(`"Software" + "Type"`, true, nil)
	defer paramResults.Close()

	results := runspace.ExecCommand("Get-ItemPropertyValue", true, map[string]interface{}{
		"Path": "HKLM:\\SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion",
		"Name": paramResults.Objects[0],
	})
	// auto cleanup the results
	defer results.Close()

	validate(t, results.Objects[0].ToString(), "System")
	if len(results.Objects) != 1 {
		t.Fail()
	}
	validate(t, "", record.lines)
}

type jsonMarshalFailure struct {
	err string
}

func (obj jsonMarshalFailure) MarshalJSON() ([]byte, error) {
	return nil, errors.New(obj.err)
}

func TestCpowershellCommandArgumentTypePanic(t *testing.T) {
	record.Reset()
	// create a runspace (where you run your powershell statements in)
	runspace := CreateRunspace(fmtPrintLogger{}, nil)
	// auto cleanup your runspace
	defer runspace.Close()

	var caughtUnknownArgumentType bool
	func() {
		defer func() {
			if r := recover(); r != nil {
				if r == "unknown argument type" {
					caughtUnknownArgumentType = true
				}
			}
		}()
		paramResults := runspace.ExecScript(`"Software" + "Type"`, true, nil, 1)
		defer paramResults.Close()
	}()

	if !caughtUnknownArgumentType {
		t.Fail()
	}
	caughtUnknownArgumentType = false

	func() {
		defer func() {
			if r := recover(); r != nil {
				if r == "unknown argument type" {
					caughtUnknownArgumentType = true
				}
			}
		}()
		results := runspace.ExecCommand("Get-ItemPropertyValue", true, map[string]interface{}{
			"Name": 1,
		})
		defer results.Close()
	}()
	if !caughtUnknownArgumentType {
		t.Fail()
	}

	caughtUnknownArgumentType = false
	func() {
		defer func() {
			if r := recover(); r != nil {
				err := r.(*json.MarshalerError)
				if err.Err.Error() == "2" {
					caughtUnknownArgumentType = true
				}
			}
		}()
		results := runspace.ExecCommandJSONMarshalUnknown("Get-ItemPropertyValue", true, map[string]interface{}{
			"Name": jsonMarshalFailure{"2"},
		})
		defer results.Close()
	}()
	if !caughtUnknownArgumentType {
		t.Fail()
	}
	caughtUnknownArgumentType = false
	func() {
		defer func() {
			if r := recover(); r != nil {
				err := r.(*json.MarshalerError)
				if err.Err.Error() == "3" {
					caughtUnknownArgumentType = true
				}
			}
		}()
		results := runspace.ExecScriptJSONMarshalUnknown("Get-ItemPropertyValue", true, nil, jsonMarshalFailure{"3"})
		defer results.Close()
	}()
	if !caughtUnknownArgumentType {
		t.Fail()
	}

	validate(t, "", record.lines)
}

func TestCpowershellJsonMarshal(t *testing.T) {
	record.Reset()
	// create a runspace (where you run your powershell statements in)
	runspace := CreateRunspace(fmtPrintLogger{}, nil)
	// auto cleanup your runspace
	defer runspace.Close()

	paramResults := runspace.ExecCommandJSONMarshalUnknown(`Write-Host`, true, map[string]interface{}{"Object": 17})
	defer paramResults.Close()
	expected := `    In Logging : 17
`
	validate(t, expected, record.lines)
}

func TestClogWchart_lookupFail(t *testing.T) {
	cStr := makeCStringUintptr("test String")
	caughtFailedToLoad := false
	defer func() {
		if r := recover(); r != nil {
			str := r.(string)
			validate(t, "failed to load context key:1", str)
			caughtFailedToLoad = true
		}
	}()
	loggerCallbackDebugln(1, cStr)
	if !caughtFailedToLoad {
		t.Fail()
	}
}
