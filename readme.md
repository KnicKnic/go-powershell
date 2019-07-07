[![Build Status](https://dev.azure.com/oneeyedelf1/powershell.native/_apis/build/status/KnicKnic.go-powershell?branchName=master)](https://dev.azure.com/oneeyedelf1/powershell.native/_build/latest?definitionId=3&branchName=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/KnicKnic/go-powershell)](https://goreportcard.com/report/github.com/KnicKnic/go-powershell)
[![gopherbadger](https://img.shields.io/badge/Go%20Coverage-97%25-brightgreen.svg?longCache=true&style=flat)](./scripts/code_coverage.ps1)
[![GoDoc](https://godoc.org/github.com/KnicKnic/go-powershell/pkg/powershell?status.svg)](https://godoc.org/github.com/KnicKnic/go-powershell/pkg/powershell)

# Status
It works
1. call scripts / cmdlets
1. reuse variables between calls / invocation
1. communicate from golang to powershell
1. communicate from powershell to golang
1. trap host output in powershell and call custom logging routines in golang
1. has tests
1. Docs - if you missed the badge above go to https://godoc.org/github.com/KnicKnic/go-powershell/pkg/powershell

This project is not api stable, however I believe it will be simple if you do use the current api to migrate to any future changes. 

**TODO:** 
1. add some code for easy back and forth json serialization of complex objects
2. more examples / tests
3. example / helper classes around exception
4. a doc overview

# Goal
The goal of this project is to enable you to quickly write golang code and interact with windows via powershell. Because powershell is a powerful scripting language you will sometimes want to call back into golang. This is also permitted. Also due to sometimes wanting to host .net and powershell giving you an easy way to wrap .net modules and functions and objects, this project also enables that.

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
	defer runspace.Delete()
	
	// execute a statement in powershell consisting of "emit this string"
	// this will output that object into the results
	results := runspace.ExecScript( `"emit this string"`, true, nil)
	// auto cleanup all results returned
	defer results.Close()
	
	// print the string result of the first object
	fmt.Println(results.Objects[0].ToString())
	// Output: emit this string
}
```

## Dependencies
This project has a dependency on [native-powershell](https://github.com/KnicKnic/native-powershell). This is a c++/cli project that enables interacting with powershell through a C DLL interface.

### Using native-powershell
1. Download host.h & psh_host.dll from https://github.com/KnicKnic/native-powershell/releases
1. copy host.h into this /pkg/powershell
1. Copy the compiled psh_host.dll into
    1. /pkg/powershell
    1. the same folder where you distribute the golang binary

### Getting cgo (so you can compile)
Windows - install dependencies - Use choco (easiest way to install gcc)

1. `choco install mingw -y`


# Docs
https://grokbase.com/t/gg/golang-nuts/154m672a6t/go-nuts-linking-cgo-with-visual-studio-x64-release-libraries-on-windows
