package powershell

/*

#cgo CFLAGS: -I.
#cgo LDFLAGS: ./psh_host.dll


#include <stddef.h>
#include "powershell.h"

*/
import "C"
import "unsafe"

// Object representing an object return from a powershell invocation
//
// Needs to be called Close on if returned from Invoke,
// unless given back to powershell and told powershell to call close.
// You do not need to call Close on those objects presented during Callbacks.
//
// This behavior is useful in send-hostcommand when you cannot execute after returnign to powershell to call close
type Object struct {
	handle C.PowerShellObject
}

// toCHandle gets the backing handle of Object
func (obj Object) toCHandle() C.PowerShellObject {
	// return *((*C.PowerShellObject)(unsafe.Pointer(&obj.handle)))
	return obj.handle
}

// // toCHandle gets the backing handle of Object
// func makeCHandles(objects []Object) []C.PowerShellObject {
// 	cHandles := make([]C.PowerShellObject, len(objects))
// 	for i,object := range(objects){
// 		cHandles[i] = object.handle
// 	}
// 	return cHandles
// }

// Close allows the memory for the powershell object to be reclaimed
//
// Should be called on all objects returned from Command.Invoke unless you have called CallbackResultsWriter.Write() with autoclose
//
// Needs to be called for every object returned from AddRef
func (obj Object) Close() {
	C.ClosePowerShellObject(obj.toCHandle())
}

// AddRef returns a new Object that has to also be called Close on
//
// This is useful in Callback processing, as those PowershellObjects are auto closed, and to keep
// a reference after the function returns use AddRef
func (obj Object) AddRef() Object {
	handle := C.AddPSObjectHandle(obj.toCHandle())
	return makePowerShellObject(handle)
}

// IsNull returns true if the backing powershell object is null
func (obj Object) IsNull() bool {
	return C.IsPSObjectNullptr(obj.toCHandle()) == 1
}

// Type returns the (System.Object).GetType().ToString() function
//
// for nullptr returns nullptr
func (obj Object) Type() string {
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
func (obj Object) ToString() string {
	if obj.IsNull() {
		return "nullptr"
	}

	var str *C.wchar_t = C.GetPSObjectToString(obj.toCHandle())
	defer C.FreeWrapper(unsafe.Pointer(str))
	return makeString(str)
}
