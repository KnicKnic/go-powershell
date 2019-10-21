package powershell

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

/*

#cgo CFLAGS: -I.
#cgo LDFLAGS: -static ${SRCDIR}/../../bin/psh_host.dll


#include <stddef.h>
#include <string.h>
#include "powershell.h"

*/
import "C"

func makeString(str *C.wchar_t) string {
	ptr := unsafe.Pointer(str)
	count := C.wcslen(str) + 1
	arr := make([]uint16, count)
	ptrwchar := unsafe.Pointer(&arr[0])

	C.memcpy(ptrwchar, ptr, count*2)

	s := windows.UTF16ToString(arr)
	return s
}

func makeCString(str string) *C.wchar_t {
	cs, _ := windows.UTF16PtrFromString(str)
	ptrwchar := unsafe.Pointer(cs)
	return C.MallocCopy((*C.wchar_t)(ptrwchar))
}

//export logWchart
// logWchart the C function pointer that dispatches to the Golang function for SimpleLogging
func logWchart(context uint64, str *C.wchar_t) {
	if context != 0 {
		s := makeString(str)
		contextInterface := getRunspaceContext(context)
		contextInterface.log.Write(s)
	}
}

//export logWarningWchart
// logWarningWchart the C function pointer that dispatches to the Golang function for SimpleLogging
func logWarningWchart(context uint64, str *C.wchar_t) {
	if context != 0 {
		s := makeString(str)
		contextInterface := getRunspaceContext(context)
		contextInterface.log.Warning(s)
	}
}

//export logInformationWchart
// logInformationWchart the C function pointer that dispatches to the Golang function for SimpleLogging
func logInformationWchart(context uint64, str *C.wchar_t) {
	if context != 0 {
		s := makeString(str)
		contextInterface := getRunspaceContext(context)
		contextInterface.log.Information(s)
	}
}

//export logVerboseWchart
// logVerboseWchart the C function pointer that dispatches to the Golang function for SimpleLogging
func logVerboseWchart(context uint64, str *C.wchar_t) {
	if context != 0 {
		s := makeString(str)
		contextInterface := getRunspaceContext(context)
		contextInterface.log.Verbose(s)
	}
}

//export logDebugWchart
// logDebugWchart the C function pointer that dispatches to the Golang function for SimpleLogging
func logDebugWchart(context uint64, str *C.wchar_t) {
	if context != 0 {
		s := makeString(str)
		contextInterface := getRunspaceContext(context)
		contextInterface.log.Debug(s)
	}
}

//export logErrorWchart
// logErrorWchart the C function pointer that dispatches to the Golang function for SimpleLogging
func logErrorWchart(context uint64, str *C.wchar_t) {
	if context != 0 {
		s := makeString(str)
		contextInterface := getRunspaceContext(context)
		contextInterface.log.Error(s)
	}
}


//export loglnWchart
// loglnWchart the C function pointer that dispatches to the Golang function for SimpleLogging
func loglnWchart(context uint64, str *C.wchar_t) {
	if context != 0 {
		s := makeString(str)
		contextInterface := getRunspaceContext(context)
		contextInterface.log.Writeln(s)
	}
}

//export logWarninglnWchart
// logWarninglnWchart the C function pointer that dispatches to the Golang function for SimpleLogging
func logWarninglnWchart(context uint64, str *C.wchar_t) {
	if context != 0 {
		s := makeString(str)
		contextInterface := getRunspaceContext(context)
		contextInterface.log.Warningln(s)
	}
}

//export logInformationlnWchart
// logInformationlnWchart the C function pointer that dispatches to the Golang function for SimpleLogging
func logInformationlnWchart(context uint64, str *C.wchar_t) {
	if context != 0 {
		s := makeString(str)
		contextInterface := getRunspaceContext(context)
		contextInterface.log.Informationln(s)
	}
}

//export logVerboselnWchart
// logVerboselnWchart the C function pointer that dispatches to the Golang function for SimpleLogging
func logVerboselnWchart(context uint64, str *C.wchar_t) {
	if context != 0 {
		s := makeString(str)
		contextInterface := getRunspaceContext(context)
		contextInterface.log.Verboseln(s)
	}
}

//export logDebuglnWchart
// logDebuglnWchart the C function pointer that dispatches to the Golang function for SimpleLogging
func logDebuglnWchart(context uint64, str *C.wchar_t) {
	if context != 0 {
		s := makeString(str)
		contextInterface := getRunspaceContext(context)
		contextInterface.log.Debugln(s)
	}
}

//export logErrorlnWchart
// logErrorlnWchart the C function pointer that dispatches to the Golang function for SimpleLogging
func logErrorlnWchart(context uint64, str *C.wchar_t) {
	if context != 0 {
		s := makeString(str)
		contextInterface := getRunspaceContext(context)
		contextInterface.log.Errorln(s)
	}
}

//export commandWchart
// commandWchart the C function pointer that dispatches to the Golang function for Send-HostCommand
func commandWchart(context uint64, cMessage *C.wchar_t, input *C.NativePowerShell_PowerShellObject, inputCount uint64, ret *C.NativePowerShell_JsonReturnValues) {

	var resultsWriter callbackResultsWriter
	if context != 0 {
		contextInterface := getRunspaceContext(context)
		inputArr := make([]Object, inputCount)
		for i := uint32(0); uint64(i) < inputCount; i++ {
			inputArr[i] = makePowerShellObjectIndexed(input, i)
		}
		message := makeString(cMessage)
		contextInterface.callback.Callback(contextInterface.recreateRunspace(), message, inputArr, &resultsWriter)
	}
	resultsWriter.filloutResults(ret)
}
