package powershell

import (
	"fmt"
)

func Example_powershellStatement() {
	// create a runspace (where you run your powershell statements in)
	runspace := CreateRunspaceSimple()
	// auto cleanup your runspace
	defer runspace.Close()

	statements := `$os = $env:OS;
				   "emitting your os is $os"`
	// execute a statement in powershell consisting of "emitting your os is $os"
	// $os will be Windows_NT
	results := runspace.ExecScript(statements, true, nil)
	// auto cleanup all results returned
	defer results.Close()

	fmt.Println(results.Objects[0].ToString())

	// OUTPUT: emitting your os is Windows_NT
}

func Example_savingVariablesAcrossStatements() {
	// create a runspace (where you run your powershell statements in)
	runspace := CreateRunspaceSimple()
	// auto cleanup your runspace
	defer runspace.Close()

	// gets whatever environment variable we request
	//     wrapping $args[0] inside $() so powershell understands [0] associated with $args
	getEnvironmentVariable := `$environVariable = get-childitem "env:\$($args[0])";`

	// Execute the statement
	// false - says to not execute the statement in a temporary child scope
	//     meaning that the variables will be available to future invocations
	// nil - means we didn't name any arguments
	// "OS" - after first 3 parameters comes the unnamed arguments which we reference via $args[index]
	results1 := runspace.ExecScript(getEnvironmentVariable, false, nil, "OS")
	//not defering close as we do not need the results
	results1.Close()

	returnEnvironmentInfo := `"emitting your $($environVariable.Name) is $($environVariable.Value)"`
	// true - we are choosing the create in a temporary child scope, the parent scope variables are still accessible to us
	//     we could however choose to specify false and be in the same scope
	results2 := runspace.ExecScript(returnEnvironmentInfo, false, nil)
	// auto cleanup all results returned
	defer results2.Close()

	// print the string result of the first object from the last statement (which happens to already be a string)
	fmt.Println(results2.Objects[0].ToString())

	// Output: emitting your OS is Windows_NT
}

func Example_powershellCommand() {
	// create a runspace (where you run your powershell statements in)
	runspace := CreateRunspaceSimple()
	// auto cleanup your runspace
	defer runspace.Close()

	// this will get the registry key for HKEY_LOCAL_MACHINE
	results := runspace.ExecCommand("get-item", true, nil, `hklm:\`)
	// auto cleanup the results
	defer results.Close()

	// print the .ToString() of a registry key, which is the key name
	fmt.Println(results.Objects[0].ToString())

	// OUTPUT: HKEY_LOCAL_MACHINE
}

func Example_powershellCommandWithNamedParameters() {
	// create a runspace (where you run your powershell statements in)
	runspace := CreateRunspace(fmtPrintLogger{}, nil)
	// auto cleanup your runspace
	defer runspace.Close()

	results := runspace.ExecCommand("Get-ItemPropertyValue", true, map[string]interface{}{
		"Path": "HKLM:\\SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion",
		"Name": "SoftwareType",
	})
	// auto cleanup the results
	defer results.Close()

	// print the .ToString() of a registry key, which is the key name
	fmt.Println(results.Objects[0].ToString())

	// OUTPUT: System
}

type printer struct{}

func (_ printer) Write(arg string) {
	fmt.Print(arg)
}

func Example_powershellJsonMarshal() {

	// create a runspace (where you run your powershell statements in)
	runspace := CreateRunspace(printer{}, nil)
	// auto cleanup your runspace
	defer runspace.Close()

	// write to host the parameters that are passed in
	command := `write-host "$args"; foreach($x in $args) {write-host $x};`
	results := runspace.ExecScriptJsonMarshalUnknown(command, true, nil, 1, 2, false, "test string", []int{1, 2, 3}, map[string]string{"fruit": "apple", "vegetable": "celery"})
	// auto cleanup the results
	defer results.Close()

	// OUTPUT:
	// Debug: 1 2 false test string [1,2,3] {"fruit":"apple","vegetable":"celery"}
	// Debug: 1
	// Debug: 2
	// Debug: false
	// Debug: test string
	// Debug: [1,2,3]
	// Debug: {"fruit":"apple","vegetable":"celery"}
}

type person struct {
	Category int
	Name     string
	Human    bool
}

func Example_powershellJsonUnmarshal() {

	// create a runspace (where you run your powershell statements in)
	runspace := CreateRunspace(printer{}, nil)
	// auto cleanup your runspace
	defer runspace.Close()

	// emit a json object with the following fields
	command := `@{"Name"= "Knic";"Category"=4;"Human"=$true} |ConvertTo-Json -Depth 3`
	results := runspace.ExecScript(command, true, nil)
	// auto cleanup the results
	defer results.Close()

	// Unmarshal into custom object person
	var me person
	results.Objects[0].JsonUnmarshal(&me)

	fmt.Print("Name: ", me.Name, ", Category: ", me.Category, ", Human: ", me.Human)
	// OUTPUT: Name: Knic, Category: 4, Human: true
}

// func Example_powershellCommandWithNamedParametersComplex() {
// 	// create a runspace (where you run your powershell statements in)
// 	runspace := CreateRunspace(fmtPrintLogger{}, nil)
// 	// auto cleanup your runspace
// 	defer runspace.Close()

// 	command := runspace.createCommand()
// 	// auto cleanup your command
// 	defer command.Close()

// 	// Get-ItemPropertyValue "HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion" -Name SoftwareType
// 	command.AddCommand("Get-ItemPropertyValue", true)
// 	command.AddParameterString("Path", "HKLM:\\SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion")
// 	command.AddParameterString("Name", "SoftwareType")
// 	// this will get the registry value for System
// 	results := command.Invoke()
// 	// auto cleanup the results
// 	defer results.Close()

// 	// print the .ToString() of a registry key, which is the key name
// 	fmt.Println(results.Objects[0].ToString())

// 	// OUTPUT: System
// }

/**
if !results.Exception.IsNull() {
	results2 := runspace.ExecScript("args[0].ToString()", true, nil, results.Exception)
	defer results2.Close()
	fmt.Println(results2.Objects[0].ToString())
}*/
