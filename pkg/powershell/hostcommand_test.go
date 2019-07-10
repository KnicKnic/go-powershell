package powershell

import (
	"fmt"
	"strconv"
)

// the callback we want to add
type callbackAdd10 struct {
}

func (callbackAdd10) Callback(runspace Runspace, str string, input []Object, results CallbackResultsWriter) {
	switch str {
	// check if we are processing the "add 10" message
	case "add 10":
		// iterate through all items passed in
		for _, object := range input {
			numStr := object.ToString()
			num, _ := strconv.Atoi(numStr)

			// write the object back to powershell as a string
			results.WriteString(fmt.Sprint(num + 10))
		}
	}
}
func ExampleCallbackHolder() {
	// create a runspace (where you run your powershell statements in)
	runspace := CreateRunspace(nil, callbackAdd10{})
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
