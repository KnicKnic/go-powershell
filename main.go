package main

import (
	"os"
	"unsafe"

	"bitbucket.org/creachadair/shell"
	"golang.org/x/sys/windows"
)

/*

#cgo CFLAGS: -I.
#cgo LDFLAGS: ./psh_host.dll
#include <stddef.h>
#include "host.h"
#include <stdio.h>
#include <stdlib.h>
#include "sys/types.h"
// #include <metahost.h>
// #pragma comment(lib, "mscoree.lib")

// ICLRMetaHost       *pMetaHost       = NULL;
// ICLRMetaHostPolicy *pMetaHostPolicy = NULL;
// ICLRDebugging      *pCLRDebugging   = NULL;

void myprint(void* unknown) {
    wchar_t* s = (wchar_t*)unknown;
// HRESULT hr = CLRCreateInstance(&CLSID_CLRMetaHost, &IID_ICLRMetaHost,
//                     (LPVOID*)&pMetaHost);
	printf("\n%ws, old_main %d\n", s);
}

wchar_t * MakeWchar(void *unknown){
	wchar_t* s = (wchar_t*)unknown;;
	return s;
}

*/
import "C"

type Runspace struct {
	handle C.RunspaceHandle
}

type PowershellCommand struct {
	handle C.PowershellHandle
}

// CreateRunspace think of this kinda like a shell
func CreateRunspace() Runspace {
	return Runspace{C.CreateRunspace()}
}

// Delete and free a Runspace
func (runspace Runspace) Delete() {
	C.DeleteRunspace(runspace.handle)
}

// CreatePowershellCommand using a runspace, still need to create a command in the powershell command
func (runspace Runspace) CreatePowershellCommand() PowershellCommand {
	return PowershellCommand{C.CreatePowershell(runspace.handle)}
}

// Delete and free a PowershellCommand
func (command PowershellCommand) Delete() {
	C.DeletePowershell(command.handle)
}

// AddCommand to an existing powershell command
func (command PowershellCommand) AddCommand(commandlet string) {
	cs, _ := windows.UTF16PtrFromString(commandlet)

	ptrwchar := unsafe.Pointer(cs)

	_ = C.AddCommand(command.handle, C.MakeWchar(ptrwchar))
}

// AddArgument to an existing powershell command
func (command PowershellCommand) AddArgument(argument string) {
	cs, _ := windows.UTF16PtrFromString(argument)

	ptrwchar := unsafe.Pointer(cs)

	_ = C.AddArgument(command.handle, C.MakeWchar(ptrwchar))
}

// Invoke the powershell command, do not reuse afterwards
func (command PowershellCommand) Invoke() {

	_ = C.InvokeCommand(command.handle)
}

// ExecStr - executes a commandline in powershell
func (runspace Runspace) ExecStr(commandStr string) {
	command := runspace.CreatePowershellCommand()
	defer command.Delete()

	fields, ok := shell.Split(commandStr)
	if !ok {
		panic("command was invalid {" + commandStr + "}")
	}
	command.AddCommand(fields[0])
	for i := 1; i < len(fields); i++ {
		command.AddArgument(fields[i])
	}
	command.Invoke()
}

// Example on how to use powershell wrappers
func Example() {
	runspace := CreateRunspace()
	defer runspace.Delete()

	for i := 1; i < len(os.Args); i++ {
		runspace.ExecStr("mkdir " + os.Args[i])
	}
}
func main() {
	Example()
}
