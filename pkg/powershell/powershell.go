package powershell

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

/*

#cgo CFLAGS: -I.
#cgo LDFLAGS: ./psh_host.dll


#include <stddef.h>
#include "powershell.h"

*/
import "C"

// Command represents a powershell command, must call Delete
type Command struct {
	handle C.PowershellHandle
}

// InvokeResults the results of an Invoke on a Command
type InvokeResults struct {
	Objects   []Object
	Exception Object
}

// CreateCommand using a runspace, still need to create a command in the powershell command
func (runspace Runspace) CreateCommand() Command {
	return Command{C.CreatePowershell(runspace.handle)}
}

// Delete and free a Command
// , call on all objects even those that are Invoked
func (command Command) Delete() {
	C.DeletePowershell(command.handle)
}

// AddCommand to an existing powershell command
func (command Command) AddCommand(commandlet string, useLocalScope bool) {
	cs, _ := windows.UTF16PtrFromString(commandlet)

	ptrwchar := unsafe.Pointer(cs)

	if useLocalScope {
		_ = C.AddCommandSpecifyScope(command.handle, (*C.wchar_t)(ptrwchar), 1)

	} else {
		_ = C.AddCommandSpecifyScope(command.handle, (*C.wchar_t)(ptrwchar), 0)

	}
}

// AddScript to an existing powershell command
func (command Command) AddScript(script string, useLocalScope bool) {
	cs, _ := windows.UTF16PtrFromString(script)

	ptrwchar := unsafe.Pointer(cs)

	if useLocalScope {
		_ = C.AddScriptSpecifyScope(command.handle, (*C.wchar_t)(ptrwchar), 1)

	} else {
		_ = C.AddScriptSpecifyScope(command.handle, (*C.wchar_t)(ptrwchar), 0)

	}
}

// AddArgumentString add a string argument to an existing powershell command
func (command Command) AddArgumentString(argument string) {
	cs, _ := windows.UTF16PtrFromString(argument)

	ptrwchar := unsafe.Pointer(cs)

	_ = C.AddArgument(command.handle, (*C.wchar_t)(ptrwchar))
}

// AddArgument add a Object argument to an existing powershell command
func (command Command) AddArgument(objects ...Object) {
	for _, object := range objects {
		_ = C.AddPSObjectArgument(command.handle, object.handle)
	}
}

// Invoke the powershell command
//
// If wanting to call another powershell command do not reuse after Invoke, create another Command object and use that one
//
// Must still call Delete on this object
func (command Command) Invoke() InvokeResults {

	var objects *C.PowerShellObject
	var count C.uint
	exception := C.InvokeCommand(command.handle, &objects, &count)
	return makeInvokeResults(objects, count, exception)
}

func makePowerShellObjectIndexed(objects *C.PowerShellObject, index uint32) Object {
	// I don't get why I have to use unsafe.Pointer on C memory
	ptr := unsafe.Pointer(objects)
	offset := (uintptr(index) * unsafe.Sizeof(*objects))
	var obj C.PowerShellObject = *(*C.PowerShellObject)(unsafe.Pointer(uintptr(ptr) + offset))
	return makePowerShellObject(obj)
}

func makePowerShellObject(object C.PowerShellObject) Object {
	// var obj uint64 = *(*uint64)(unsafe.Pointer(&object))
	// return Object{obj}
	return Object{object}
}

func makeInvokeResults(objects *C.PowerShellObject, count C.uint, exception C.PowerShellObject) (results InvokeResults) {
	goCount := uint32(count)
	results.Objects = make([]Object, count)
	for i := uint32(0); i < goCount; i++ {
		results.Objects[i] = makePowerShellObjectIndexed(objects, i)
	}
	results.Exception = makePowerShellObject(exception)
	return
}

// Close is a convenience wrapper to close all result objects, so you do not have to
func (results *InvokeResults) Close() {
	for _, object := range results.Objects {
		object.Close()
	}
	results.Exception.Close()
}

// Success returns true if the powershell command executed successfully (threw no terminating exceptions)
func (results *InvokeResults) Success() bool {
	return results.Exception.IsNull()
}
