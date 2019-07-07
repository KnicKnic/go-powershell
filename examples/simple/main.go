package main

import (
	"fmt"

	"github.com/KnicKnic/go-powershell/pkg/powershell"
)

func main() {
	// create a runspace (where you run your powershell statements in)
	runspace := powershell.CreateRunspaceSimple()
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
