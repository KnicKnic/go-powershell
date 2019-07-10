[![Build Status](https://dev.azure.com/oneeyedelf1/powershell.native/_apis/build/status/KnicKnic.go-powershell?branchName=master)](https://dev.azure.com/oneeyedelf1/powershell.native/_build/latest?definitionId=3&branchName=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/KnicKnic/go-powershell)](https://goreportcard.com/report/github.com/KnicKnic/go-powershell)
[![gopherbadger](https://img.shields.io/badge/Go%20Coverage-99%25-brightgreen.svg?longCache=true&style=flat)](./scripts/code_coverage.ps1)
[![GoDoc](https://godoc.org/github.com/KnicKnic/go-powershell/pkg/powershell?status.svg)](https://godoc.org/github.com/KnicKnic/go-powershell/pkg/powershell)
[![GitHub commits since latest release (branch)](https://img.shields.io/github/commits-since/KnicKnic/go-powershell/latest.svg)](https://github.com/KnicKnic/go-powershell/releases/latest)

# Goal
The goal of this project is to enable you to quickly write golang code and interact with windows via powershell and not use exec. Because powershell is a powerful scripting language you will sometimes want to call back into golang. This is also enabled by this project. Also due to sometimes wanting to host .net and powershell giving you an easy way to wrap .net modules and functions and objects, this project also enables that.

Features:
1. Call from golang to powershell
1. Call from powershell to golang (via special Send-HostCommand commandlet)
1. Easy logging - Trap host output in powershell and call custom logging routines in golang


# Status
It works
1. call scripts / cmdlets
1. reuse variables between calls / invocation
1. Call from golang to powershell
1. Call from powershell back to golang (via special Send-HostCommand commandlet)
1. trap host output in powershell and call custom logging routines in golang
1. has automted tests
1. Docs - if you missed the badge above go to https://godoc.org/github.com/KnicKnic/go-powershell/pkg/powershell

This project is not api stable, however I believe it will be simple if you do use the current api to migrate to any future changes. 

**TODO:** 

- [x] add some code for easy back and forth json serialization of complex objects
- [X] more examples / tests
- [ ] example / helper classes around exception
- [ ] a doc overview
- [ ] support for default loggers, like glog or log (in seperate package)

# Usage
```go
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
	//not deferring close as we do not need the results
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
```

## Dependencies
This project has a dependency on [native-powershell](https://github.com/KnicKnic/native-powershell). This is a c++/cli project that enables interacting with powershell through a C DLL interface.

### Using native-powershell
1. Simply fetch the dependencies, `go get -d .` and then make sure to build, `go build`
1. Copy the precompiled psh_host.dll into your location so it can be found when running the app
    1. cmd - `copy %GOPATH%\src\github.com\KnicKnic\go-powershell\native-powershell\native-powershell-bin\psh_host.dll .`
    1. powershell - `copy "$($env:GOPATH)\src\github.com\KnicKnic\go-powershell\native-powershell\native-powershell-bin\psh_host.dll" .`
1. I ended up checking in the psh_host.dll and host.h (to make things easy)
    1. I could not find a better way to go about this and still have things be easy.

### Getting cgo (so you can compile)
Windows - install dependencies - you need gcc. I Use chocolatey to install (easiest way to install gcc)

1. Install chocolatey
	1. https://chocolatey.org/docs/installation#installing-chocolatey
1. `choco install mingw -y`


# Docs
https://grokbase.com/t/gg/golang-nuts/154m672a6t/go-nuts-linking-cgo-with-visual-studio-x64-release-libraries-on-windows
