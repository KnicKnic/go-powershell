package powershell

// "bitbucket.org/creachadair/shell"

/*

#cgo CFLAGS: -I.
#cgo LDFLAGS: ./psh_host.dll


#include <stddef.h>
#include "powershell.h"

*/
import "C"

type CallbackResultsWriter interface {
	WriteString(string)
	Write(object PowershellObject, needsClose bool)
}
type CallbackHolder interface {
	Callback(str string, input []PowershellObject, results CallbackResultsWriter)
}

type callbackResultsWriter struct {
	objects []C.GenericPowershellObject
}

func (writer *callbackResultsWriter) WriteString(str string) {
	cStr := makeCString(str)
	var obj C.GenericPowershellObject
	C.SetGenericPowershellString(&obj, cStr, C.char(1))
	writer.objects = append(writer.objects, obj)
}

func (writer *callbackResultsWriter) Write(handle PowershellObject, needsClose bool) {
	var obj C.GenericPowershellObject
	var autoClose C.char = C.char(1)
	if needsClose {
		autoClose = 1
	}
	C.SetGenericPowershellHandle(&obj, C.ulonglong(handle.handle), autoClose)
	writer.objects = append(writer.objects, obj)
}

func (writer *callbackResultsWriter) filloutResults(results *C.JsonReturnValues) {
	results.objects = nil
	results.count = 0
	if writer.objects != nil {
		results.count = C.ulong(len(writer.objects))
		results.objects = C.MallocCopyGenericPowershellObject(&writer.objects[0], C.ulonglong(len(writer.objects)))
	}
}
