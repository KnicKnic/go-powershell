package powershell

/*

#cgo CFLAGS: -I.
#cgo LDFLAGS: ${SRCDIR}/../../native-powershell/native-powershell-bin/psh_host.dll


#include <stddef.h>
#include "powershell.h"

*/
import "C"
import "unsafe"
import "encoding/json"

// Object representing an object return from a powershell invocation
//
// Should be called on all objects returned from a powershell invocation (not callback parameters)
//
// See note on Object.Close for exceptions & more rules about Close
type Object struct {
	handle C.NativePowerShell_PowerShellObject
}

// toCHandle gets the backing handle of Object
func (obj Object) toCHandle() C.NativePowerShell_PowerShellObject {
	// return *((*C.NativePowerShell_PowerShellObject)(unsafe.Pointer(&obj.handle)))
	return obj.handle
}

// // toCHandle gets the backing handle of Object
// func makeCHandles(objects []Object) []C.NativePowerShell_PowerShellObject {
// 	cHandles := make([]C.NativePowerShell_PowerShellObject, len(objects))
// 	for i,object := range(objects){
// 		cHandles[i] = object.handle
// 	}
// 	return cHandles
// }

// Close allows the memory for the powershell object to be reclaimed
//
// Should be called on all objects returned from a powershell invocation (not callback parameters)
//
//     Exception: Do not call Close on the object when inside a callback and calling CallbackResultsWriter.Write() with autoclose
//
// Needs to be called for every object returned from AddRef
func (obj Object) Close() {
	C.NativePowerShell_ClosePowerShellObject(obj.toCHandle())
}

// AddRef returns a new Object that has to also be called Close on
//
// This is useful in Callback processing, as those NativePowerShell_PowerShellObjects are auto closed, and to keep
// a reference after the function returns use AddRef
func (obj Object) AddRef() Object {
	handle := C.NativePowerShell_AddPSObjectHandle(obj.toCHandle())
	return makePowerShellObject(handle)
}

// IsNull returns true if the backing powershell object is null
func (obj Object) IsNull() bool {
	return C.NativePowerShell_IsPSObjectNullptr(obj.toCHandle()) == 1
}

// Type returns the (System.Object).GetType().ToString() function
//
// for nullptr returns nullptr
func (obj Object) Type() string {
	if obj.IsNull() {
		return "nullptr"
	}

	var str *C.wchar_t = C.NativePowerShell_GetPSObjectType(obj.toCHandle())
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

	var str *C.wchar_t = C.NativePowerShell_GetPSObjectToString(obj.toCHandle())
	defer C.FreeWrapper(unsafe.Pointer(str))
	return makeString(str)
}

// JSONUnmarshal calls the ToString function and unmarshals it into the supplied object
func (obj Object) JSONUnmarshal(userObject interface{}) error {
	bytes := []byte(obj.ToString())
	return json.Unmarshal(bytes, userObject)
}
