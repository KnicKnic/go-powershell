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

func CreateRunspace() C.RunspaceHandle {
	return C.CreateRunspace()
}
func DeleteRunspace(handle C.RunspaceHandle) {
	C.DeleteRunspace(handle)
}
func CreatePowershell(handle C.RunspaceHandle) C.PowershellHandle {
	return C.CreatePowershell(handle)
}
func DeletePowershell(handle C.PowershellHandle) {
	C.DeletePowershell(handle)
}

func AddCommand(handle C.PowershellHandle, command string) {
	cs, _ := windows.UTF16PtrFromString(command)

	ptrwchar := unsafe.Pointer(cs)

	_ = C.AddCommand(handle, C.MakeWchar(ptrwchar))
}
func AddArgument(handle C.PowershellHandle, argument string) {
	cs, _ := windows.UTF16PtrFromString(argument)

	ptrwchar := unsafe.Pointer(cs)

	_ = C.AddArgument(handle, C.MakeWchar(ptrwchar))
}
func InvokeCommand(handle C.PowershellHandle) {

	_ = C.InvokeCommand(handle)
}

type RunspaceHandle struct {
	handle C.RunspaceHandle
}

func ExecStr(runspace C.RunspaceHandle, command string) {
	powershell := CreatePowershell(runspace)
	defer DeletePowershell(powershell)

	fields, ok := shell.Split(command)
	if !ok {
		panic("command was invalid {" + command + "}")
	}
	AddCommand(powershell, fields[0])
	for i := 1; i < len(fields); i++ {
		AddArgument(powershell, fields[i])
	}
	InvokeCommand(powershell)
}
func (runspace RunspaceHandle) ExecStr2(command string) {
	ExecStr(runspace.handle, command)
}

func MakeDir(handle C.PowershellHandle, path string) {
	AddCommand(handle, "mkdir")
	AddArgument(handle, path)
	InvokeCommand(handle)

}

func RunMakeDir(runspace C.RunspaceHandle, path string) {

	powershell := CreatePowershell(runspace)
	defer DeletePowershell(powershell)

	println("mkdir ", path)
	MakeDir(powershell, path)
}

func Example() {
	runspace := CreateRunspace()
	defer DeleteRunspace(runspace)

	for i := 1; i < len(os.Args); i++ {
		// RunMakeDir(runspace, os.Args[i])
		handle := RunspaceHandle{runspace}
		// ExecStr(runspace, "mkdir "+os.Args[i])
		handle.ExecStr2("mkdir " + os.Args[i])
	}
}
func main() {
	Example()
}
