package powershell

import (
	"fmt"
	"strconv"
)

// the callback we want to add
type callbackAdd10Nested struct {
}

func (callback callbackAdd10Nested) Callback(runspace Runspace, str string, input []Object, results CallbackResultsWriter) {
	switch str {
	// check if we are processing the "add 10" message
	case "add 10":
		// iterate through all items passed in
		for _, object := range input {
			numStr := object.ToString()
			num, _ := strconv.Atoi(numStr)
			num += 10

			// // write the object back to powershell as a string
			// results.WriteString(fmt.Sprint(num))

			// or write them back as a powershell integer
			//
			// convert object into a powershell integer
			//
			// execute in anonyous function to get scoped cleanup of results
			func() {
				execResults := runspace.ExecScript(`[int]$args[0]`, true, nil, fmt.Sprint(num))
				defer execResults.Close()

				// we need to close our execResults.Object[0] for us after it has been processed
				// however we do not know when that is, so tell the results to auto do it
				// WE MUST NOT CLOSE IT OURSELVES IF SPECIFYING TRUE!
				results.Write(execResults.Objects[0], true)
				execResults.RemoveObjectFromClose(0)
			}()
		}
	}
}
func Example_powershellCallbackNested() {
	callback := callbackAdd10Nested{}

	// create a runspace (where you run your powershell statements in)
	runspace := CreateRunspace(nil, callback)
	// auto cleanup your runspace
	defer runspace.Close()

	statements := `1..3 | Send-HostCommand -message "add 10"`
	results := runspace.ExecScript(statements, true, nil)
	// auto cleanup all results returned
	defer results.Close()

	for _, num := range results.Objects {
		fmt.Println(num.ToString())
	}

	// OUTPUT:
	// 11
	// 12
	// 13
}
