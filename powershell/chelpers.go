package powershell

// "bitbucket.org/creachadair/shell"
import (
	"unsafe"

	"github.com/golang/glog"

	"golang.org/x/sys/windows"
)

/*

#cgo CFLAGS: -I.
#cgo LDFLAGS: ./psh_host.dll


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
func logWchart(context uint64, str *C.wchar_t) {
	if context != 0 {
		s := makeString(str)
		// glog.Info("golang log: ", s)

		contextInterface, ok := GetRunspaceContext(context)
		if ok {
			contextInterface.Log.Log.Verbose(s)
		} else {
			glog.Info("In Logging callback, failed to load context key: ", context)
		}
	}
}

//export commandWchart
func commandWchart(context uint64, cMessage *C.wchar_t, input *C.PowerShellObject, inputCount uint64, ret *C.JsonReturnValues) {

	if context != 0 {
		contextInterface, ok := GetRunspaceContext(context)
		if ok {
			inputArr := make([]PowershellObject, inputCount)
			for i := uint32(0); uint64(i) < inputCount; i++ {
				inputArr[i] = makePowerShellObjectIndexed(input, i)
			}
			message := makeString(cMessage)
			var resultsWriter callbackResultsWriter
			contextInterface.Callback.Callback(message, inputArr, &resultsWriter)
			// resultsWriter = callbackResultsWriter{}
			resultsWriter.filloutResults(ret)
			return
		} else {
			glog.Info("In Command callback, failed to load context key: ", context)
		}
	}
	var resultsWriter callbackResultsWriter
	resultsWriter.filloutResults(ret)
}