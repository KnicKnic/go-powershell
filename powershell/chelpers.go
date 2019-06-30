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
#include "powershell.h"

*/
import "C"

func makeString(str *C.wchar_t) string {
	var count C.int = 0
	var zero C.wchar_t = C.MakeNullTerminator()
	for ; C.GetChar(str, count) != zero; count++ {
	}
	count++
	arr := make([]uint16, count)
	arrPtr := &arr[0]
	ptrwchar := unsafe.Pointer(arrPtr)

	C.MemoryCopy(ptrwchar, str, count*2)

	s := windows.UTF16ToString(arr)
	return s
}

func makeCString(str string) *C.wchar_t {
	cs, _ := windows.UTF16PtrFromString(str)
	ptrwchar := unsafe.Pointer(cs)
	return C.MallocCopy(C.MakeWchar(ptrwchar))

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

type callbackTest struct{}

func (c callbackTest) Callback(str string, input []PowershellObject, results CallbackResultsWriter) {
	glog.Info("In callback: ", str)
	results.WriteString(str)
	for _, object := range input {
		results.Write(object)
	}
	// 	return "returned from callback: " + s
}
