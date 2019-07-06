package powershell

import (
	"fmt"
)

func Example_powershellStatement() {
	// create a runspace (where you run your powershell statements in)
	runspace := CreateRunspaceSimple()
	// auto cleanup your runspace
	defer runspace.Delete()

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
	defer runspace.Delete()

	statement1 := `$os = $env:OS;`
	// The false says to not execute the statement in a seperate scope
	// meaning that the variables will be availale to future invocations
	results1 := runspace.ExecScript(statement1, false, nil)
	//not defering close as we do not need the results
	results1.Close()

	statement2 := `"emitting your os is $os"`
	// we are choosing the create in a different scope, the parent scope variables are accessible to us
	// we could however choose to specify false and be in the same scope
	results2 := runspace.ExecScript(statement2, true, nil)
	// auto cleanup all results returned
	defer results2.Close()

	fmt.Println(results2.Objects[0].ToString())

	// OUTPUT: emitting your os is Windows_NT
}

func Example_powershellCommand() {
	// create a runspace (where you run your powershell statements in)
	runspace := CreateRunspaceSimple()
	// auto cleanup your runspace
	defer runspace.Delete()

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
	defer runspace.Delete()

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

// func Example_powershellCommandWithNamedParametersComplex() {
// 	// create a runspace (where you run your powershell statements in)
// 	runspace := CreateRunspace(fmtPrintLogger{}, nil)
// 	// auto cleanup your runspace
// 	defer runspace.Delete()

// 	command := runspace.CreateCommand()
// 	// auto cleanup your command
// 	defer command.Delete()

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
