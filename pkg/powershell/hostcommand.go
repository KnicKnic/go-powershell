package powershell

/*

#cgo CFLAGS: -I.
#cgo LDFLAGS: ${SRCDIR}/../../native-powershell/native-powershell-bin/psh_host.dll


#include <stddef.h>
#include "powershell.h"

*/
import "C"

// CallbackResultsWriter allows you to write values to powershell when inside Send-HostCommand
type CallbackResultsWriter interface {
	WriteString(string)
	Write(object Object, needsClose bool)
}

// CallbackHolder callback function pointer for Send-HostCommand callbacks
type CallbackHolder interface {
	Callback(str string, input []Object, results CallbackResultsWriter)
}

// callbackResultsWriter is the internal implementation of CallbackResultsWriter
type callbackResultsWriter struct {
	objects []C.GenericPowershellObject
}

// WriteString accumulates a string object to return from Send-HostCommand
func (writer *callbackResultsWriter) WriteString(str string) {
	cStr := makeCString(str)
	var obj C.GenericPowershellObject
	C.SetGenericPowershellString(&obj, cStr, 1)
	writer.objects = append(writer.objects, obj)
}

// Write accumulates a string object to return from Send-HostCommand
func (writer *callbackResultsWriter) Write(handle Object, needsClose bool) {
	var obj C.GenericPowershellObject
	var autoClose C.char
	if needsClose {
		autoClose = 1
	}
	C.SetGenericPowershellHandle(&obj, handle.toCHandle(), autoClose)
	writer.objects = append(writer.objects, obj)
}

// filloutResults takes accumulated objects from Write calls and prepares them to cross the C boundary
func (writer *callbackResultsWriter) filloutResults(results *C.JsonReturnValues) {
	results.objects = nil
	results.count = 0
	if writer.objects != nil {
		results.count = C.ulong(len(writer.objects))
		results.objects = C.MallocCopyGenericPowershellObject(&writer.objects[0], C.ulonglong(len(writer.objects)))
	}
}
