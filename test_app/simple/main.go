package main

import (
	"fmt"
	"github.com/KnicKnic/go-powershell/pkg/powershell"
)

func main() {
	// create a runspace (where you run your powershell statements in)
	runspace := powershell.CreateRunspaceSimple()
	// auto cleanup your runspace
	defer runspace.Delete()

	// execute a statement in powershell consisting of "emit this string"
	// this will output that object into the results
	results := runspace.ExecStr(`"emit this string"`, true)
	// auto cleanup all results returned
	defer results.Close()

	// print the string result of the first object
	fmt.Println(results.Objects[0].ToString())
	// Output: emit this string
}
