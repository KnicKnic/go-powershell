package powershell

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// psCommand represents a powershell command, must call Close
type psCommand struct {
	handle  nativePowerShell_PowerShellHandle
	context *runspaceContext
}

// InvokeResults the results of an Invoke on a psCommand
type InvokeResults struct {
	Objects        []Object
	Exception      Object
	objectsNoClose map[int]bool //only using as a set
}

// createCommand using a runspace, still need to create a command in the powershell command
func (runspace Runspace) createCommand() psCommand {
	currentlyInvoking := runspace.invoking
	if len(currentlyInvoking) != 0 {
		currentCommand := currentlyInvoking[len(currentlyInvoking)-1]
		return currentCommand.createNested()
	}
	handle := nativePowerShell_CreatePowerShell(runspace.handle)

	return psCommand{handle, runspace.runspaceContext}
}

// createNested a nested powershell command
func (command psCommand) createNested() psCommand {
	return psCommand{nativePowerShell_CreatePowerShellNested(command.handle), command.context}
}

// Close and free a psCommand
// , call on all objects even those that are Invoked
func (command psCommand) Close() {
	nativePowerShell_DeletePowershell(command.handle)
}

func boolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}

// AddCommand to an existing powershell command
func (command psCommand) AddCommand(commandlet string, useLocalScope bool) {
	cs, _ := windows.UTF16PtrFromString(commandlet)

	localScope := boolToByte(useLocalScope)

	_ = nativePowerShell_AddCommandSpecifyScope(command.handle, cs, localScope)
}

// AddScript to an existing powershell command
func (command psCommand) AddScript(script string, useLocalScope bool) {
	cs, _ := windows.UTF16PtrFromString(script)

	localScope := boolToByte(useLocalScope)

	_ = nativePowerShell_AddScriptSpecifyScope(command.handle, cs, localScope)
}

// AddArgumentString add a string argument to an existing powershell command
func (command psCommand) AddArgumentString(argument string) {
	cs, _ := windows.UTF16PtrFromString(argument)

	_ = nativePowerShell_AddArgument(command.handle, cs)
}

// AddArgument add a Object argument to an existing powershell command
func (command psCommand) AddArgument(object Object) {
	_ = nativePowerShell_AddPSObjectArgument(command.handle, object.handle)
}

// AddParameterString add a string with a parameter name to an existing powershell command
func (command psCommand) AddParameterString(paramName string, paramValue string) {
	cName, _ := windows.UTF16PtrFromString(paramName)

	cValue, _ := windows.UTF16PtrFromString(paramValue)
	_ = nativePowerShell_AddParameterString(command.handle, cName, cValue)
}

// AddParameter add a Object with a parameter name to an existing powershell command
func (command psCommand) AddParameter(paramName string, object Object) {

	cName, _ := windows.UTF16PtrFromString(paramName)

	_ = nativePowerShell_AddParameterObject(command.handle, cName, object.handle)
}

func (command psCommand) completeInvoke() {
	command.context.invoking = command.context.invoking[:len(command.context.invoking)-1]
}

// Invoke the powershell command
//
// If wanting to call another powershell command do not reuse after Invoke, create another psCommand object and use that one
//
// Must still call Close on this object
func (command psCommand) Invoke() *InvokeResults {

	var objects uintptr
	var count uint
	command.context.invoking = append(command.context.invoking, command)
	defer command.completeInvoke()
	exception := nativePowerShell_InvokeCommand(command.handle, &objects, &count)
	return makeInvokeResults(objects, count, exception)
}

func makePowerShellObjectIndexed(objects uintptr, index uint32) Object {
	// I don't get why I have to use unsafe.Pointer on C memory
	var handle nativePowerShell_PowerShellObject

	offset := (uintptr(index) * unsafe.Sizeof(handle))
	handle = *(*nativePowerShell_PowerShellObject)(unsafe.Pointer(objects + offset))
	return makePowerShellObject(handle)
}

func makePowerShellObject(object nativePowerShell_PowerShellObject) Object {
	// var obj uint64 = *(*uint64)(unsafe.Pointer(&object))
	// return Object{obj}
	return Object{object}
}

func makeInvokeResults(objects uintptr, count uint, exception nativePowerShell_PowerShellObject) *InvokeResults {
	results := InvokeResults{Objects: make([]Object, count),
		Exception:      makePowerShellObject(exception),
		objectsNoClose: make(map[int]bool),
	}
	goCount := uint32(count)
	for i := uint32(0); i < goCount; i++ {
		results.Objects[i] = makePowerShellObjectIndexed(objects, i)
	}
	return &results
}

// RemoveObjectFromClose remove object from objects that get closed from Close routine.
// Does not alter InvokeResults.Objects
//
// Useful when calling powershell from inside a powershell callback and need to
// to call CallbackResultsWriter.Write(object, true) to have powershell cleanup the reference
func (results *InvokeResults) RemoveObjectFromClose(index int) {
	results.objectsNoClose[index] = true
}

// Close is a convenience wrapper to close all result objects, and the exception so you do not have to
func (results *InvokeResults) Close() {
	for i, object := range results.Objects {
		if !results.objectsNoClose[i] {
			object.Close()
		}
	}
	results.Exception.Close()
}

// Success returns true if the powershell command executed successfully (threw no terminating exceptions)
func (results *InvokeResults) Success() bool {
	return results.Exception.IsNull()
}
