package powershell

/*

#cgo CFLAGS: -I.
#cgo LDFLAGS: ./psh_host.dll


#include <stddef.h>
#include "powershell.h"

*/
import "C"
import "unsafe"

// a context handle representing an object return from a powershell invocation
//
// Needs to be called Close on if returned from Invoke,
// unless given back to powershell and told powershell to call close.
// You do not need to call Close on those objects presented during Callbacks.
//
// This behavior is useful in send-hostcommand when you cannot execute after returnign to powershell to call close
type PowershellObject struct {
	handle C.PowerShellObject
}

// toCHandle gets the backing handle of PowershellObject
func (obj PowershellObject) toCHandle() C.PowerShellObject {
	// return *((*C.PowerShellObject)(unsafe.Pointer(&obj.handle)))
	return obj.handle
}

// Close allows the memory for the powershell object to be reclaimed
//
// Should be called on all objects returned from PowershellCommand.Invoke unless you have called CallbackResultsWriter.Write() with autoclose
//
// Needs to be called for every object returned from AddRef
func (obj PowershellObject) Close() {
	C.ClosePowerShellObject(obj.toCHandle())
}

// AddRef returns a new PowershellObject that has to also be called Close on
//
// This is usefull in Callback processing, as those PowershellObjects are auto closed, and to keep
// a reference after the function returns use AddRef
func (obj PowershellObject) AddRef() PowershellObject {
	handle := C.AddPSObjectHandle(obj.toCHandle())
	return makePowerShellObject(handle)
}

// IsNull returns true if the backing powershell object is null
func (obj PowershellObject) IsNull() bool {
	return C.IsPSObjectNullptr(obj.toCHandle()) == 1
}

// Type returns the (System.Object).GetType().ToString() function
//
// for nullptr returns nullptr
func (obj PowershellObject) Type() string {
	if obj.IsNull() {
		return "nullptr"
	}

	var str *C.wchar_t = C.GetPSObjectType(obj.toCHandle())
	defer C.FreeWrapper(unsafe.Pointer(str))
	return makeString(str)
}

// ToString returns the (System.Object).ToString() function
//
// for nullptr returns nullptr
func (obj PowershellObject) ToString() string {
	if obj.IsNull() {
		return "nullptr"
	}

	var str *C.wchar_t = C.GetPSObjectToString(obj.toCHandle())
	defer C.FreeWrapper(unsafe.Pointer(str))
	return makeString(str)
}
