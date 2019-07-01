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

type PowershellCommand struct {
	handle C.PowershellHandle
}

type InvokeResults struct {
	objects   []PowershellObject
	count     uint32
	exception PowershellObject
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
func (command PowershellCommand) AddCommand(commandlet string, useLocalScope bool) {
	cs, _ := windows.UTF16PtrFromString(commandlet)

	ptrwchar := unsafe.Pointer(cs)

	if useLocalScope {
		_ = C.AddCommandSpecifyScope(command.handle, C.MakeWchar(ptrwchar), 1)

	} else {
		_ = C.AddCommandSpecifyScope(command.handle, C.MakeWchar(ptrwchar), 0)

	}
}

// AddCommand to an existing powershell command
func (command PowershellCommand) AddScript(script string, useLocalScope bool) {
	cs, _ := windows.UTF16PtrFromString(script)

	ptrwchar := unsafe.Pointer(cs)

	if useLocalScope {
		_ = C.AddScriptSpecifyScope(command.handle, C.MakeWchar(ptrwchar), 1)

	} else {
		_ = C.AddScriptSpecifyScope(command.handle, C.MakeWchar(ptrwchar), 0)

	}
}

// AddArgument to an existing powershell command
func (command PowershellCommand) AddArgument(argument string) {
	cs, _ := windows.UTF16PtrFromString(argument)

	ptrwchar := unsafe.Pointer(cs)

	_ = C.AddArgument(command.handle, C.MakeWchar(ptrwchar))
}

// Invoke the powershell command, do not reuse afterwards
func (command PowershellCommand) Invoke() {

	var objects *C.PowerShellObject
	var count C.uint
	exception := C.InvokeCommand(command.handle, &objects, &count)
	makeInvokeResults(objects, count, exception)
}

func makePowerShellObjectIndexed(objects *C.PowerShellObject, index uint32) PowershellObject {
	// I don't get why I have to use unsafe.Pointer on C memory
	ptr := uintptr(unsafe.Pointer(objects))
	handle := ptr + (uintptr(index) * unsafe.Sizeof(*objects))
	var obj C.PowerShellObject = *(*C.PowerShellObject)(unsafe.Pointer(handle))
	return makePowerShellObject(obj)
}

func makePowerShellObject(object C.PowerShellObject) PowershellObject {
	// var obj uint64 = *(*uint64)(unsafe.Pointer(&object))
	// return PowershellObject{obj}
	return PowershellObject{object}
}

func makeInvokeResults(objects *C.PowerShellObject, count C.uint, exception C.PowerShellObject) (results InvokeResults) {
	results.count = uint32(count)
	results.objects = make([]PowershellObject, count)
	for i := uint32(0); i < results.count; i++ {
		results.objects[i] = makePowerShellObjectIndexed(objects, i)
	}
	results.exception = makePowerShellObject(exception)
	return
}
