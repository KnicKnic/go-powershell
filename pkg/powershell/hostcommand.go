package powershell

import "unsafe"

// CallbackResultsWriter allows you to write values to powershell when inside Send-HostCommand
type CallbackResultsWriter interface {
	WriteString(string)
	Write(object Object, needsClose bool)
}

// CallbackHolder callback function pointer for Send-HostCommand callbacks
type CallbackHolder interface {
	Callback(runspace Runspace, message string, input []Object, results CallbackResultsWriter)
}

// CallbackFuncPtr a simple implementation of CallbackHolder that lets you pass in a function pointer for the callback
type CallbackFuncPtr struct {
	FuncPtr func(runspace Runspace, message string, input []Object, results CallbackResultsWriter)
}

// Callback is the function that will call the function pointer in CallbackFuncPtr
func (callback CallbackFuncPtr) Callback(runspace Runspace, message string, input []Object, results CallbackResultsWriter) {
	callback.FuncPtr(runspace, message, input, results)
}

// callbackResultsWriter is the internal implementation of CallbackResultsWriter
type callbackResultsWriter struct {
	objects []nativePowerShell_GenericPowerShellObject
}

func setGenericPowershellString(object *nativePowerShell_GenericPowerShellObject, value uintptr, autoRelease byte) {
	object.typeEnum = nativePowerShell_PowerShellObjectTypeString
	object.object = makeUint64FromPtr(value)
	object.releaseObject = autoRelease
}

func setGenericPowerShellHandle(object *nativePowerShell_GenericPowerShellObject, value nativePowerShell_PowerShellObject, autoRelease byte) {
	object.typeEnum = nativePowerShell_PowerShellObjectHandle
	object.object = value
	object.releaseObject = autoRelease
}

// WriteString accumulates a string object to return from Send-HostCommand
func (writer *callbackResultsWriter) WriteString(str string) {
	cStr := makeCStringUintptr(str)
	var obj nativePowerShell_GenericPowerShellObject
	setGenericPowershellString(&obj, cStr, 1)
	writer.objects = append(writer.objects, obj)
}

// Write accumulates a string object to return from Send-HostCommand
func (writer *callbackResultsWriter) Write(handle Object, needsClose bool) {
	var obj nativePowerShell_GenericPowerShellObject
	var autoClose = boolToByte(needsClose)

	setGenericPowerShellHandle(&obj, handle.toCHandle(), autoClose)
	writer.objects = append(writer.objects, obj)
}

func mallocCopyGenericPowerShellObject(input *nativePowerShell_GenericPowerShellObject, inputCount uint64) *nativePowerShell_GenericPowerShellObject {

	size := uintptr(inputCount) * unsafe.Sizeof(*input)

	return (*nativePowerShell_GenericPowerShellObject)(unsafe.Pointer(mallocCopy(uintptr(unsafe.Pointer(input)), size)))
}

// filloutResults takes accumulated objects from Write calls and prepares them to cross the C boundary
func (writer *callbackResultsWriter) filloutResults(res uintptr) {
	results := (*nativePowerShell_JsonReturnValues)(unsafe.Pointer(res))
	results.objects = nil
	results.count = 0
	if writer.objects != nil && len(writer.objects) > 0 {
		results.count = uint32(len(writer.objects))
		results.objects = mallocCopyGenericPowerShellObject(&writer.objects[0], uint64(len(writer.objects)))
	}
}
