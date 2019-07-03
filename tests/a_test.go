package tests

import (
	"fmt"
	powershell "github.com/KnicKnic/go-powershell"
)

// GLogInfoLogger is a simple struct that provides ability to send logs to glog at Info level
type fmtPrintLogger struct {
}

func (logger fmtPrintLogger) Write(arg string) {
	fmt.Print("    In Logging : ", arg)
}

type callbackTest struct{}

func (c callbackTest) Callback(str string, input []powershell.Object, results powershell.CallbackResultsWriter) {
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
func PrintAndExecuteCommand(runspace powershell.Runspace, command string, useLocalScope bool) {
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

// // Example on how to use powershell wrappers
// func TestFoo(*testing.T) {
// 	runspace := powershell.CreateRunspace(fmtPrintLogger{}, callbackTest{})
// 	// runspace := CreateRunspaceSimple()
// 	defer runspace.Delete()

// 	PrintAndExecuteCommand(runspace, ".\\t1.ps1", false)	
// 	PrintAndExecuteCommand(runspace, ".\\t2.ps1", false)	
// }

// func ExampleCreateRunspaceSimple(){

// 	runspace := powershell.CreateRunspace(fmtPrintLogger{}, callbackTest{})
// 	// runspace := CreateRunspaceSimple()
// 	defer runspace.Delete()

// 	PrintAndExecuteCommand(runspace, ".\\t1.ps1", false)	
// 	PrintAndExecuteCommand(runspace, ".\\t2.ps1", false)	
// 	// Output:
// // Executing powershell command: .\t1.ps1
// //     In Logging : Debug: ab  ba
// //     In Logging : Debug: ab  ba
// //     In Logging : Debug: ab  ba
// //     In callback: asdfasdf
// //     In callback: index 0 type: System.Int32 with value: 1
// //     In callback: index 1 type: System.Int32 with value: 2
// //     In callback: index 2 type: System.Int32 with value: 3
// //     In Logging : Debug: asdfasdf 1 2 3
// //     In callback: two
// //     In Logging : Debug: two
// //     In callback: three
// //     In Logging : Debug: three
// //     In callback: four
// //     In callback: index 0 Object Is Null
// //     In callback: index 0 type: nullptr with value: nullptr
// //     In callback: index 1 Object Is Null
// //     In callback: index 1 type: nullptr with value: nullptr
// //     In Logging : Debug: four  
// //     In Logging : Error: someerror
// // Completed Executing powershell command: .\t1.ps1
// // Command returned 0 objects
// // Executing powershell command: .\t2.ps1
// //     In Logging : Debug: start t2
// //     In Logging : Debug: 5
// //     In Logging : Debug: 5
// //     In Logging : Debug: asdf
// // Completed Executing powershell command: .\t2.ps1
// // Command returned 0 objects
// }

func ExampleCreateRunspaceWitLoggerWithCallback(){

	runspace := powershell.CreateRunspace(fmtPrintLogger{}, callbackTest{})
	// runspace := CreateRunspaceSimple()
	defer runspace.Delete()
input_str :=`
write-host 'simple output statement'
"s1"
"s2"`
	PrintAndExecuteCommand(runspace, input_str, false)
	// Output:
// Executing powershell command:
// write-host 'simple output statement'
// "s1"
// "s2"
//     In Logging : Debug: simple output statement
// Completed Executing powershell command:
// write-host 'simple output statement'
// "s1"
// "s2"
// Command returned 2 objects
//     Object 0 is of type System.String and ToString s1
//     Object 1 is of type System.String and ToString s2
}

// func ExampleCreateRunspaceSimple(){
// 	// ensure directory is gone
// 	_ = os.Remove("ExampleCreateRunspaceSimple")
	
// 	// show directory is gone
// 	printDirectoryStatus := func (){
// 		_ = os.stat("ExampleCreateRunspaceSimple")
// 		if _, err := os.Stat(path); os.IsNotExist(err) {
// 			fmt.Println("ExampleCreateRunspaceSimple does not exist")
// 		}else {
// 			fmt.Println("ExampleCreateRunspaceSimple exists")
// 		}
// 	}

// 	runspace := powershell.CreateRunspaceSimple()
// 	defer runspace.Delete()
// 	runspace.ExecStr( `[System.Console]::WriteLine("hi")`, true)
// 	// Output: hi
// }


func ExampleCreateRunspaceSimple(){
	runspace := powershell.CreateRunspaceSimple()
	defer runspace.Delete()
	results := runspace.ExecStr( `"emit this string"`, true)
	defer results.Close()
	// print the string result of the first object
	fmt.Println(results.Objects[0].ToString())
	// Output: emit this string
}